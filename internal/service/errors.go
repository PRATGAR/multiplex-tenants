package service

import "errors"

var (
	ErrInvalid        = errors.New("invalid input")
	ErrNotFound       = errors.New("not found")
	ErrBadCredentials = errors.New("bad credentials")
	ErrNotVerified    = errors.New("domain ownership not verified")
	ErrSeatLimit      = errors.New("seat count exceeds plan limit")
)
