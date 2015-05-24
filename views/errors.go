package views

import "errors"

var (
	ErrForbidden           = errors.New("forbidden")
	ErrInternalServerError = errors.New("internal_server_error")
	ErrInvalidParam        = errors.New("invalid_param")
	ErrNotFound            = errors.New("not_found")
	ErrUnauthorized        = errors.New("unauthorized")
)
