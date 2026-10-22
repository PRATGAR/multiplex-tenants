package handlers

import (
	"net/http"
	"strings"

	"tenantly/internal/httpx"
	"tenantly/internal/service"
)

type Auth struct {
	svc *service.Auth
}

func NewAuth(svc *service.Auth) *Auth {
	return &Auth{svc: svc}
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Tenant   string `json:"tenant"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	token, err := h.svc.Login(r.Context(), in.Tenant, in.Email, in.Password)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Auth) Callback(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusFound)
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	_ = h.svc.Logout(r.Context(), token)
	http.Redirect(w, r, r.URL.Query().Get("return"), http.StatusFound)
}
