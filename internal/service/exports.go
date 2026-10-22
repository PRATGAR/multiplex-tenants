package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"os/exec"
	"strconv"

	"tenantly/internal/exporter"
	"tenantly/internal/model"
	"tenantly/internal/store"
)

type Exports struct {
	tenants *store.Tenants
	users   *store.Users
	repo    *store.Exports
	files   *exporter.Store
}

func NewExports(t *store.Tenants, u *store.Users, repo *store.Exports, files *exporter.Store) *Exports {
	return &Exports{tenants: t, users: u, repo: repo, files: files}
}

func (s *Exports) Create(ctx context.Context, tenantID int64) (model.Export, error) {
	t, err := s.tenants.Get(ctx, tenantID)
	if err != nil {
		return model.Export{}, ErrNotFound
	}
	members, err := s.users.All(ctx, tenantID)
	if err != nil {
		return model.Export{}, err
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"id", "email", "display_name", "role"})
	for _, m := range members {
		_ = w.Write([]string{strconv.FormatInt(m.ID, 10), m.Email, m.DisplayName, m.Role})
	}
	w.Flush()
	id, err := s.repo.Insert(ctx, tenantID, len(members))
	if err != nil {
		return model.Export{}, err
	}
	if err := s.files.WriteData(tenantID, id, buf.Bytes()); err != nil {
		return model.Export{}, err
	}
	if err := s.files.PublishManifest(t.Slug, id, len(members)); err != nil {
		return model.Export{}, err
	}
	return model.Export{ID: id, TenantID: tenantID, Rows: len(members)}, nil
}

func (s *Exports) Read(ctx context.Context, tenantID, id int64) ([]byte, error) {
	ok, err := s.repo.Exists(ctx, tenantID, id)
	if err != nil || !ok {
		return nil, ErrNotFound
	}
	return s.files.ReadData(tenantID, id)
}

func (s *Exports) Gzip(ctx context.Context, tenantID, id int64, level string) error {
	ok, err := s.repo.Exists(ctx, tenantID, id)
	if err != nil || !ok {
		return ErrNotFound
	}
	source := s.files.DataPath(tenantID, id)
	cmd := exec.Command("sh", "-c", "gzip -"+level+" -k "+source)
	return cmd.Run()
}
