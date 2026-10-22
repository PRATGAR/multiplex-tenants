package handlers

import (
	"fmt"
	"net/http"

	"tenantly/internal/httpx"
	"tenantly/internal/service"
)

type Tenants struct {
	svc *service.Tenants
}

func NewTenants(svc *service.Tenants) *Tenants {
	return &Tenants{svc: svc}
}

func (h *Tenants) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		Plan       string `json:"plan"`
		OwnerEmail string `json:"owner_email"`
		OwnerName  string `json:"owner_name"`
		Password   string `json:"password"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	id, err := h.svc.Create(r.Context(), service.CreateTenantInput{
		Name: in.Name, Slug: in.Slug, Plan: in.Plan,
		OwnerEmail: in.OwnerEmail, OwnerName: in.OwnerName, Password: in.Password,
	})
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *Tenants) Update(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	var in struct {
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		InvoiceFooter string `json:"invoice_footer"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	if err := h.svc.Update(r.Context(), tid, in.Name, in.Slug, in.InvoiceFooter); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Tenants) Welcome(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	t, err := h.svc.Get(r.Context(), tid)
	if err != nil {
		fail(w, service.ErrNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Welcome to %s</h1>", t.Name)
}
