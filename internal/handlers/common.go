package handlers

import (
	"errors"
	"net/http"

	"tenantly/internal/httpx"
	"tenantly/internal/middleware"
	"tenantly/internal/model"
	"tenantly/internal/service"
)

func fail(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrBadCredentials):
		httpx.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrNotVerified), errors.Is(err, service.ErrSeatLimit):
		httpx.Error(w, http.StatusUnprocessableEntity, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal error")
	}
}

func session(r *http.Request) model.Session {
	sess, _ := middleware.SessionFrom(r.Context())
	return sess
}
