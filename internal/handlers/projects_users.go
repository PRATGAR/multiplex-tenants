package handlers

import (
	"net/http"
	"strconv"

	"tenantly/internal/httpx"
	"tenantly/internal/service"
)

type Projects struct {
	svc *service.Projects
}

func NewProjects(svc *service.Projects) *Projects {
	return &Projects{svc: svc}
}

func (h *Projects) ForTenant(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ForTenant(r.Context(), httpx.PathString(r, "tid"))
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}

func (h *Projects) Mine(w http.ResponseWriter, r *http.Request) {
	sess := session(r)
	list, err := h.svc.ForTenant(r.Context(), strconv.FormatInt(sess.TenantID, 10))
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}

func (h *Projects) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, err := h.svc.Search(r.Context(), session(r).TenantID, q.Get("q"), q.Get("sort"))
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}

type Users struct {
	svc *service.Users
}

func NewUsers(svc *service.Users) *Users {
	return &Users{svc: svc}
}

func (h *Users) Activity(w http.ResponseWriter, r *http.Request) {
	tid, okT := httpx.PathInt(r, "tid")
	uid, okU := httpx.PathInt(r, "uid")
	if !okT || !okU {
		httpx.Error(w, http.StatusBadRequest, "bad id")
		return
	}
	entries, err := h.svc.Activity(r.Context(), tid, uid)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, entries)
}
