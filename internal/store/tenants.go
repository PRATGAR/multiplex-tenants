package store

import (
	"context"
	"database/sql"

	"tenantly/internal/model"
)

type Tenants struct {
	db *sql.DB
}

func NewTenants(db *sql.DB) *Tenants {
	return &Tenants{db: db}
}

func (t *Tenants) Create(ctx context.Context, name, slug, plan string) (int64, error) {
	var id int64
	err := t.db.QueryRowContext(ctx,
		`INSERT INTO tenants (name, slug, plan) VALUES ($1, $2, $3) RETURNING id`,
		name, slug, plan).Scan(&id)
	if err != nil {
		return 0, err
	}
	_, err = t.db.ExecContext(ctx,
		`INSERT INTO billing_accounts (tenant_id) VALUES ($1)`, id)
	return id, err
}

func (t *Tenants) Get(ctx context.Context, id int64) (model.Tenant, error) {
	var out model.Tenant
	err := t.db.QueryRowContext(ctx,
		`SELECT id, name, slug, plan, invoice_footer, created_at FROM tenants WHERE id = $1`, id).
		Scan(&out.ID, &out.Name, &out.Slug, &out.Plan, &out.InvoiceFooter, &out.CreatedAt)
	return out, err
}

func (t *Tenants) Update(ctx context.Context, id int64, name, slug, footer string) error {
	_, err := t.db.ExecContext(ctx,
		`UPDATE tenants
		    SET name = COALESCE(NULLIF($2, ''), name),
		        slug = COALESCE(NULLIF($3, ''), slug),
		        invoice_footer = COALESCE(NULLIF($4, ''), invoice_footer)
		  WHERE id = $1`, id, name, slug, footer)
	return err
}
