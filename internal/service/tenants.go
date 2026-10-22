package service

import (
	"context"
	"strings"

	"tenantly/internal/config"
	"tenantly/internal/exporter"
	"tenantly/internal/model"
	"tenantly/internal/store"
	"tenantly/internal/validate"
)

type Tenants struct {
	tenants *store.Tenants
	users   *store.Users
	files   *exporter.Store
	salt    string
}

func NewTenants(t *store.Tenants, u *store.Users, files *exporter.Store, salt string) *Tenants {
	return &Tenants{tenants: t, users: u, files: files, salt: salt}
}

type CreateTenantInput struct {
	Name       string
	Slug       string
	Plan       string
	OwnerEmail string
	OwnerName  string
	Password   string
}

func (s *Tenants) Create(ctx context.Context, in CreateTenantInput) (int64, error) {
	if !validate.Slug(in.Slug) || !validate.Email(in.OwnerEmail) || !validate.DisplayName(in.Name) {
		return 0, ErrInvalid
	}
	if len(in.Password) < 10 {
		return 0, ErrInvalid
	}
	plan, ok := config.ResolvePlan(in.Plan)
	if !ok {
		plan = config.Plans["starter"]
	}
	id, err := s.tenants.Create(ctx, strings.TrimSpace(in.Name), in.Slug, plan.Code)
	if err != nil {
		return 0, err
	}
	if _, err := s.users.Create(ctx, id, strings.TrimSpace(in.OwnerEmail), in.OwnerName, "owner",
		HashPassword(s.salt, in.Password)); err != nil {
		return 0, err
	}
	return id, s.files.Prepare(id)
}

func (s *Tenants) Get(ctx context.Context, id int64) (model.Tenant, error) {
	return s.tenants.Get(ctx, id)
}

func (s *Tenants) Update(ctx context.Context, id int64, name, slug, footer string) error {
	if len(slug) > 64 || len(name) > 120 || len(footer) > 2000 {
		return ErrInvalid
	}
	return s.tenants.Update(ctx, id, strings.TrimSpace(name), strings.TrimSpace(slug), footer)
}
