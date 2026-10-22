package store

import (
	"context"
	"database/sql"
)

type Exports struct {
	db *sql.DB
}

func NewExports(db *sql.DB) *Exports {
	return &Exports{db: db}
}

func (e *Exports) Insert(ctx context.Context, tenantID int64, rows int) (int64, error) {
	var id int64
	err := e.db.QueryRowContext(ctx,
		`INSERT INTO exports (tenant_id, row_count) VALUES ($1, $2) RETURNING id`, tenantID, rows).Scan(&id)
	return id, err
}

func (e *Exports) Exists(ctx context.Context, tenantID, id int64) (bool, error) {
	var n int
	err := e.db.QueryRowContext(ctx,
		`SELECT count(*) FROM exports WHERE tenant_id = $1 AND id = $2`, tenantID, id).Scan(&n)
	return n > 0, err
}
