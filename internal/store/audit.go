package store

import (
	"context"
	"database/sql"
	"strconv"

	"tenantly/internal/model"
)

type Audit struct {
	db *sql.DB
}

func NewAudit(db *sql.DB) *Audit {
	return &Audit{db: db}
}

func (a *Audit) Record(ctx context.Context, tenantID int64, actor, action string) error {
	_, err := a.db.ExecContext(ctx,
		`INSERT INTO audit_log (tenant_id, actor, action) VALUES ($1, $2, $3)`, tenantID, actor, action)
	return err
}

func (a *Audit) ByActor(ctx context.Context, tenantID int64, actor string) ([]model.AuditEntry, error) {
	q := "SELECT action, at FROM audit_log WHERE tenant_id = " + strconv.FormatInt(tenantID, 10) +
		" AND actor = '" + actor + "' ORDER BY at DESC LIMIT 100"
	rows, err := a.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.AuditEntry
	for rows.Next() {
		var e model.AuditEntry
		if err := rows.Scan(&e.Action, &e.At); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
