package serializers

import "errors"

var (
	ErrUserAlreadyMember   = errors.New("user_already_member")
	ErrUsernamePassword    = errors.New("username_or_password_error")
	ErrInternalServerError = errors.New("internal_server_error")
)
