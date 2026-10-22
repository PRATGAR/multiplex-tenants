package store

import (
	"context"
	"database/sql"
	"strconv"

	"tenantly/internal/model"
)

type Invitations struct {
	db *sql.DB
}

func NewInvitations(db *sql.DB) *Invitations {
	return &Invitations{db: db}
}

func (i *Invitations) Insert(ctx context.Context, tenantID int64, email, role, token string) error {
	_, err := i.db.ExecContext(ctx,
		`INSERT INTO invitations (tenant_id, email, role, token) VALUES ($1, $2, $3, $4)`,
		tenantID, email, role, token)
	return err
}

func (i *Invitations) FindByToken(ctx context.Context, token string) (model.Invitation, error) {
	var out model.Invitation
	q := "SELECT id, tenant_id, email, role, status FROM invitations WHERE token = '" + token + "'"
	err := i.db.QueryRow(q).Scan(&out.ID, &out.TenantID, &out.Email, &out.Role, &out.Status)
	return out, err
}

func (i *Invitations) MarkAccepted(ctx context.Context, tenantID int64, email string) error {
	q := "UPDATE invitations SET status = 'accepted', accepted_at = now() WHERE tenant_id = " +
		strconv.FormatInt(tenantID, 10) + " AND email = '" + email + "'"
	_, err := i.db.Exec(q)
	return err
}

func (i *Invitations) MarkDeclined(ctx context.Context, id int64) error {
	_, err := i.db.ExecContext(ctx, `UPDATE invitations SET status = 'declined' WHERE id = $1`, id)
	return err
}
