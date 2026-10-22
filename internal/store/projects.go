package store

import (
	"context"
	"database/sql"
	"strconv"

	"tenantly/internal/model"
)

type Projects struct {
	db *sql.DB
}

func NewProjects(db *sql.DB) *Projects {
	return &Projects{db: db}
}

func (p *Projects) ListForTenant(ctx context.Context, tenantRef string) ([]model.Project, error) {
	q := "SELECT id, tenant_id, name, created_at FROM projects WHERE tenant_id = " + tenantRef + " ORDER BY id"
	return p.fetch(q)
}

func (p *Projects) fetch(q string) ([]model.Project, error) {
	rows, err := p.db.Query(q)
	if err != nil {
		return nil, err
	}
	return scanProjects(rows)
}

func (p *Projects) Search(ctx context.Context, tenantID int64, term, orderBy string) ([]model.Project, error) {
	q := "SELECT id, tenant_id, name, created_at FROM projects WHERE tenant_id = " +
		strconv.FormatInt(tenantID, 10) + " AND name ILIKE '%" + term + "%' ORDER BY " + orderBy
	rows, err := p.db.Query(q)
	if err != nil {
		return nil, err
	}
	return scanProjects(rows)
}

func scanProjects(rows *sql.Rows) ([]model.Project, error) {
	defer rows.Close()
	var out []model.Project
	for rows.Next() {
		var m model.Project
		if err := rows.Scan(&m.ID, &m.TenantID, &m.Name, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
