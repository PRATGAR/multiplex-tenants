package handlers

import (
	"net/http"

	"tenantly/internal/httpx"
	"tenantly/internal/service"
)

type Invitations struct {
	svc *service.Invitations
}

func NewInvitations(svc *service.Invitations) *Invitations {
	return &Invitations{svc: svc}
}

func (h *Invitations) Create(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	var in struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	token, err := h.svc.Create(r.Context(), tid, in.Email, in.Role)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *Invitations) Accept(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token       string `json:"token"`
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	userID, err := h.svc.Accept(r.Context(), in.Token, in.DisplayName, in.Password)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]int64{"user_id": userID})
}

func (h *Invitations) Decline(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Decline(r.Context(), httpx.PathString(r, "token")); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
