package store

import (
	"context"
	"database/sql"

	"tenantly/internal/model"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Create(ctx context.Context, tenantID int64, email, displayName, role, hash string) (int64, error) {
	var id int64
	err := u.db.QueryRowContext(ctx,
		`INSERT INTO users (tenant_id, email, display_name, role, password_hash)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		tenantID, email, displayName, role, hash).Scan(&id)
	return id, err
}

func (u *Users) Get(ctx context.Context, tenantID, userID int64) (model.User, error) {
	var out model.User
	err := u.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, email, display_name, role FROM users WHERE tenant_id = $1 AND id = $2`,
		tenantID, userID).Scan(&out.ID, &out.TenantID, &out.Email, &out.DisplayName, &out.Role)
	return out, err
}

func (u *Users) ByLogin(ctx context.Context, slug, email string) (model.User, error) {
	var out model.User
	err := u.db.QueryRowContext(ctx,
		`SELECT u.id, u.tenant_id, u.email, u.display_name, u.role, u.password_hash
		   FROM users u JOIN tenants t ON t.id = u.tenant_id
		  WHERE t.slug = $1 AND u.email = $2`, slug, email).
		Scan(&out.ID, &out.TenantID, &out.Email, &out.DisplayName, &out.Role, &out.PasswordHash)
	return out, err
}

func (u *Users) All(ctx context.Context, tenantID int64) ([]model.User, error) {
	rows, err := u.db.QueryContext(ctx,
		`SELECT id, tenant_id, email, display_name, role FROM users WHERE tenant_id = $1 ORDER BY id`, tenantID)
	if err != nil {
		return nil, err
	}
	return scanUsers(rows)
}

func scanUsers(rows *sql.Rows) ([]model.User, error) {
	defer rows.Close()
	var out []model.User
	for rows.Next() {
		var m model.User
		if err := rows.Scan(&m.ID, &m.TenantID, &m.Email, &m.DisplayName, &m.Role); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
