package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Store struct {
	dir string
}

func New(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) Dir() string {
	return s.dir
}

func (s *Store) Prepare(tenantID int64) error {
	for _, d := range []string{
		filepath.Join(s.dir, strconv.FormatInt(tenantID, 10)),
		filepath.Join(s.dir, "manifests"),
	} {
		if err := os.MkdirAll(d, 0o750); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DataPath(tenantID, id int64) string {
	return filepath.Join(s.dir, strconv.FormatInt(tenantID, 10), fmt.Sprintf("%d.csv", id))
}

func (s *Store) WriteData(tenantID, id int64, data []byte) error {
	return os.WriteFile(s.DataPath(tenantID, id), data, 0o640)
}

func (s *Store) ReadData(tenantID, id int64) ([]byte, error) {
	return os.ReadFile(s.DataPath(tenantID, id))
}

type manifest struct {
	ExportID  int64     `json:"export_id"`
	Rows      int       `json:"rows"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) PublishManifest(slug string, id int64, rows int) error {
	target := filepath.Join(s.dir, "manifests", slug+".json")
	return s.writeManifest(target, manifest{ExportID: id, Rows: rows, CreatedAt: time.Now().UTC()})
}

func (s *Store) writeManifest(target string, m manifest) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(m)
}
