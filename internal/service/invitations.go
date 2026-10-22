package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"tenantly/internal/model"
	"tenantly/internal/store"
	"tenantly/internal/tokens"
	"tenantly/internal/validate"
)

const inviteTTL = 72 * time.Hour

func inviteKey(token string) string {
	return "invite:" + token
}

type Invitations struct {
	repo  *store.Invitations
	users *store.Users
	audit *store.Audit
	rdb   *redis.Client
	salt  string
}

func NewInvitations(repo *store.Invitations, users *store.Users, audit *store.Audit, rdb *redis.Client, salt string) *Invitations {
	return &Invitations{repo: repo, users: users, audit: audit, rdb: rdb, salt: salt}
}

func (s *Invitations) Create(ctx context.Context, tenantID int64, email, role string) (string, error) {
	email = strings.TrimSpace(email)
	if !validate.Email(email) {
		return "", ErrInvalid
	}
	if role != "admin" && role != "member" && role != "viewer" {
		return "", ErrInvalid
	}
	token, err := tokens.Invite()
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(model.InvitePayload{TenantID: tenantID, Email: email, Role: role})
	if err != nil {
		return "", err
	}
	if err := s.rdb.Set(ctx, inviteKey(token), payload, inviteTTL).Err(); err != nil {
		return "", err
	}
	if err := s.repo.Insert(ctx, tenantID, email, role, token); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Invitations) Accept(ctx context.Context, token, displayName, password string) (int64, error) {
	if !validate.DisplayName(displayName) || len(password) < 10 {
		return 0, ErrInvalid
	}
	raw, err := s.rdb.Get(ctx, inviteKey(token)).Result()
	if err != nil {
		return 0, ErrNotFound
	}
	var payload model.InvitePayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return 0, ErrNotFound
	}
	inv, err := s.repo.FindByToken(ctx, token)
	if err != nil || inv.Status != "pending" {
		return 0, ErrNotFound
	}
	userID, err := s.users.Create(ctx, payload.TenantID, payload.Email, strings.TrimSpace(displayName),
		payload.Role, HashPassword(s.salt, password))
	if err != nil {
		return 0, err
	}
	if err := s.repo.MarkAccepted(ctx, payload.TenantID, payload.Email); err != nil {
		return 0, err
	}
	_ = s.audit.Record(ctx, payload.TenantID, strings.TrimSpace(displayName), "invitation.accepted")
	return userID, s.rdb.Del(ctx, inviteKey(token)).Err()
}

func (s *Invitations) Decline(ctx context.Context, token string) error {
	inv, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		return ErrNotFound
	}
	if err := s.repo.MarkDeclined(ctx, inv.ID); err != nil {
		return err
	}
	return s.rdb.Del(ctx, inviteKey(token)).Err()
}
