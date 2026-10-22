package service

import (
	"context"

	"tenantly/internal/model"
	"tenantly/internal/store"
)

type Projects struct {
	repo *store.Projects
}

func NewProjects(repo *store.Projects) *Projects {
	return &Projects{repo: repo}
}

func (s *Projects) ForTenant(ctx context.Context, tenantRef string) ([]model.Project, error) {
	return s.repo.ListForTenant(ctx, tenantRef)
}

func (s *Projects) Search(ctx context.Context, tenantID int64, term, orderBy string) ([]model.Project, error) {
	return s.repo.Search(ctx, tenantID, term, orderBy)
}

type Users struct {
	users *store.Users
	audit *store.Audit
}

func NewUsers(users *store.Users, audit *store.Audit) *Users {
	return &Users{users: users, audit: audit}
}

func (s *Users) Activity(ctx context.Context, tenantID, userID int64) ([]model.AuditEntry, error) {
	u, err := s.users.Get(ctx, tenantID, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.audit.ByActor(ctx, tenantID, u.DisplayName)
}
