package store

import (
	"context"
	"database/sql"

	"tenantly/internal/model"
)

type Domains struct {
	db *sql.DB
}

func NewDomains(db *sql.DB) *Domains {
	return &Domains{db: db}
}

func (d *Domains) Insert(ctx context.Context, tenantID int64, name, verifyToken string) (int64, error) {
	var id int64
	err := d.db.QueryRowContext(ctx,
		`INSERT INTO domains (tenant_id, name, verify_token) VALUES ($1, $2, $3) RETURNING id`,
		tenantID, name, verifyToken).Scan(&id)
	return id, err
}

func (d *Domains) Get(ctx context.Context, tenantID, id int64) (model.Domain, error) {
	var out model.Domain
	err := d.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, verify_token, status FROM domains WHERE tenant_id = $1 AND id = $2`,
		tenantID, id).Scan(&out.ID, &out.TenantID, &out.Name, &out.VerifyToken, &out.Status)
	return out, err
}

func (d *Domains) SetStatus(ctx context.Context, id int64, status string) error {
	_, err := d.db.ExecContext(ctx, `UPDATE domains SET status = $2 WHERE id = $1`, id, status)
	return err
}
