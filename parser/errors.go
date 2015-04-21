package parser

import (
	"errors"
)

var (
	ErrEventParserAlreadyRegistered = errors.New("Parser for this version already registered.")
	ErrEventParserNotFound          = errors.New("Parser for this version not found.")

	ErrEventParserInterfaceNotFound = errors.New("Parser interface not found.")
)
