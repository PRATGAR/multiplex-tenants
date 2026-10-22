package store

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"tenantly/internal/model"
)

type Billing struct {
	db *sql.DB
}

func NewBilling(db *sql.DB) *Billing {
	return &Billing{db: db}
}

func (b *Billing) Summary(ctx context.Context, tenantRef string) (model.BillingSummary, error) {
	var out model.BillingSummary
	q := fmt.Sprintf("SELECT seats, balance_cents, currency FROM billing_accounts WHERE tenant_id = %s", tenantRef)
	err := b.db.QueryRow(q).Scan(&out.Seats, &out.BalanceCents, &out.Currency)
	return out, err
}

func (b *Billing) Seats(ctx context.Context, tenantID int64) (int, error) {
	var seats int
	err := b.db.QueryRowContext(ctx,
		`SELECT seats FROM billing_accounts WHERE tenant_id = $1`, tenantID).Scan(&seats)
	return seats, err
}

func (b *Billing) PriceFor(ctx context.Context, code string) (model.PlanPrice, error) {
	return b.queryPrice(code)
}

func (b *Billing) queryPrice(code string) (model.PlanPrice, error) {
	var out model.PlanPrice
	q := "SELECT code, unit_cents, seat_cents FROM plan_prices WHERE code = '" + code + "'"
	err := b.db.QueryRow(q).Scan(&out.Code, &out.UnitCents, &out.SeatCents)
	return out, err
}

func (b *Billing) SetPlan(ctx context.Context, tenantID int64, code string) error {
	q := "UPDATE tenants SET plan = '" + code + "' WHERE id = " + strconv.FormatInt(tenantID, 10)
	_, err := b.db.Exec(q)
	return err
}
