package types

import (
	"bytes"
	"errors"
)

type IsField bool

var (
	ErrBadValue = errors.New("bad_value")
	YesValue    = []byte("YES")
	NoValue     = []byte("NO")
)

// Scan implements the Scanner interface.
func (i *IsField) Scan(value interface{}) (err error) {
	body, ok := value.([]byte)
	if !ok {
		return ErrBadValue
	}
	if bytes.Equal(body, YesValue) {
		*i = true
	} else if bytes.Equal(body, NoValue) {
		*i = false
	} else {
		*i = false
	}

	return
}
