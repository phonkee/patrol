package models

import "errors"

var (
	// Generic errors for base manager methods
	ErrObjectDoesNotExists = errors.New("does_not_exists")
	ErrObjectAlreadyExists = errors.New("unique_violation")

	ErrInvalidEmail        = errors.New("incorrect email address")
	ErrUpdateNoFieldsGiven = errors.New("no valid fields given to update")
	ErrCannotLoginUser     = errors.New("cannot login user")

	ErrTeamNameTooLong = errors.New("Team name should not exceed 64 characters.")

	ErrInvalidChoice = errors.New("invalid_choice")

	ErrCannotParseAuthHeaders = errors.New("cannot parser auth headers.")

	ErrNotMember = errors.New("not_member")

	ErrNilPointer = errors.New("nil pointer given")

	ErrIncorrectModel = errors.New("incorrect_model")

	ErrEventParserAlreadyRegistered = errors.New("Parser for this version already registered.")
	ErrEventParserNotFound          = errors.New("Parser for this version not found.")

	ErrEventParserInterfaceNotFound = errors.New("Parser interface not found.")
)
