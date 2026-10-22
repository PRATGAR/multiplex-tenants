package handlers

import (
	"io"
	"net/http"
	"os"

	"tenantly/internal/httpx"
	"tenantly/internal/service"
)

type Domains struct {
	svc *service.Domains
}

func NewDomains(svc *service.Domains) *Domains {
	return &Domains{svc: svc}
}

func (h *Domains) Add(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	d, err := h.svc.Add(r.Context(), tid, in.Name)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, d)
}

func (h *Domains) Verify(w http.ResponseWriter, r *http.Request) {
	tid, okT := httpx.PathInt(r, "tid")
	did, okD := httpx.PathInt(r, "did")
	if !okT || !okD {
		httpx.Error(w, http.StatusBadRequest, "bad id")
		return
	}
	d, err := h.svc.Verify(r.Context(), tid, did)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, d)
}

type Exports struct {
	svc *service.Exports
}

func NewExports(svc *service.Exports) *Exports {
	return &Exports{svc: svc}
}

func (h *Exports) Create(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	exp, err := h.svc.Create(r.Context(), tid)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, exp)
}

func (h *Exports) File(w http.ResponseWriter, r *http.Request) {
	tid, okT := httpx.PathInt(r, "tid")
	eid, okE := httpx.PathInt(r, "eid")
	if !okT || !okE {
		httpx.Error(w, http.StatusBadRequest, "bad id")
		return
	}
	data, err := h.svc.Read(r.Context(), tid, eid)
	if err != nil {
		fail(w, err)
		return
	}
	w.Write(data)
}

func (h *Exports) Download(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(r.URL.Query().Get("file"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "no such file")
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = io.Copy(w, f)
}

func (h *Exports) Compress(w http.ResponseWriter, r *http.Request) {
	eid, ok := httpx.PathInt(r, "eid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad export id")
		return
	}
	tid := session(r).TenantID
	level := r.URL.Query().Get("level")
	if level == "" {
		level = "6"
	}
	if err := h.svc.Gzip(r.Context(), tid, eid, level); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
