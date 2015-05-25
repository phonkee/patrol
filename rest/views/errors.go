package views

import "errors"

var (
	ErrMethodNotFound = errors.New("method not found")

	ErrForbidden           = errors.New("forbidden")
	ErrInternalServerError = errors.New("internal_server_error")
	ErrInvalidParam        = errors.New("invalid_param")
	ErrNotFound            = errors.New("not_found")
	ErrUnauthorized        = errors.New("unauthorized")

	ErrBreakRequest   = errors.New("break request")
	ErrParamNotFound  = errors.New("param not found")
	ErrMuxVarNotFound = errors.New("mux var not found")
)
