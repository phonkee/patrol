package views

import "errors"

var (
	ErrInternalServerError = errors.New("internal_server_error")
	ErrInvalidParam        = errors.New("invalid_param")
	ErrNotFound            = errors.New("not_found")
	ErrUnauthorized        = errors.New("unauthorized")
)
