package service

import (
	"context"
	"strings"

	"tenantly/internal/model"
	"tenantly/internal/store"
	"tenantly/internal/tokens"
	"tenantly/internal/validate"
	"tenantly/internal/verifier"
)

type Domains struct {
	repo     *store.Domains
	verifier *verifier.Verifier
}

func NewDomains(repo *store.Domains, v *verifier.Verifier) *Domains {
	return &Domains{repo: repo, verifier: v}
}

func (s *Domains) Add(ctx context.Context, tenantID int64, name string) (model.Domain, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if !validate.Hostname(name) {
		return model.Domain{}, ErrInvalid
	}
	secret, err := tokens.Secret()
	if err != nil {
		return model.Domain{}, err
	}
	id, err := s.repo.Insert(ctx, tenantID, name, secret)
	if err != nil {
		return model.Domain{}, err
	}
	return model.Domain{ID: id, TenantID: tenantID, Name: name, VerifyToken: secret, Status: "pending"}, nil
}

func (s *Domains) Verify(ctx context.Context, tenantID, domainID int64) (model.Domain, error) {
	d, err := s.repo.Get(ctx, tenantID, domainID)
	if err != nil {
		return model.Domain{}, ErrNotFound
	}
	ok, err := s.verifier.Check(ctx, d.Name, d.VerifyToken)
	if err != nil || !ok {
		return model.Domain{}, ErrNotVerified
	}
	if err := s.repo.SetStatus(ctx, d.ID, "active"); err != nil {
		return model.Domain{}, err
	}
	d.Status = "active"
	return d, nil
}
