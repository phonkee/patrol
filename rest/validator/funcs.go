package validator

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
)

var (
	ErrInvalidEmail = errors.New("invalid_email")

	ErrStringMinLength = errors.New("min_length")
	ErrStringMaxLength = errors.New("max_length")

	ErrInt64Min = errors.New("min")
	ErrInt64Max = errors.New("max")
)

// Validator function
type ValidatorFunc func(interface{}) error

/*
Any wraps multiple Validator funcs, first validation error will be returned
*/
func Any(funcs ...ValidatorFunc) ValidatorFunc {
	return func(value interface{}) (err error) {
		for _, fn := range funcs {
			if err = fn(value); err != nil {
				return
			}
		}
		return
	}
}

/*
Custom validators used
*/

func ValidateInt64Min(v int64) ValidatorFunc {
	return func(value interface{}) (err error) {
		s := value.(int64)
		if s < v {
			err = ErrInt64Min
		}
		return
	}
}

func ValidateInt64Max(v int64) ValidatorFunc {
	return func(value interface{}) (err error) {
		s := value.(int64)
		if s > v {
			err = ErrInt64Max
		}
		return
	}
}

// validates minimum string length
func ValidateStringMinLength(length int) ValidatorFunc {
	return func(value interface{}) (err error) {
		s := value.(string)
		if len(strings.TrimSpace(s)) < length {
			err = ErrStringMinLength
		}
		return
	}
}

// validates minimum string length
func ValidateStringMaxLength(length int) ValidatorFunc {
	return func(value interface{}) (err error) {
		s := value.(string)
		if len(s) > length {
			err = ErrStringMaxLength
		}
		return
	}
}

// validates email
func ValidateEmail() ValidatorFunc {
	return func(value interface{}) (err error) {
		s := value.(string)
		if !govalidator.IsEmail(s) {
			err = ErrInvalidEmail
		}
		return
	}
}

// 	// first callback fails, error is returned
// 	for _, callback := range callbacks {
// 		if err = callback(email); err != nil {
// 			v.AddFieldError(field, err)
// 			return
// 		}
// 	}

// 	return
// }

// // validates string
// func (v *Validator) ValidateString(table, column, value string, callbacks ...ValidatorStringCallback) (err error) {
// 	ci, _ := v.dbinfo.ColumnInfo(table, column)

// 	if !ci.CharacterMaximumLength.Valid {
// 		return
// 	}

// 	// size greater than allowed
// 	if int64(len(value)) > ci.CharacterMaximumLength.Int64 {
// 		err = ErrMaxLengthExceeded
// 		v.AddFieldError(column, err)
// 		return
// 	}

// 	// first callback fails, error is returned
// 	for _, callback := range callbacks {
// 		if err = callback(value); err != nil {
// 			v.AddFieldError(column, err)
// 			return
// 		}
// 	}

// 	return
// }

// /*
// ValidateColumn
// */
// func (v *Validator) ValidateIntColumn(table, column string, value int, callbacks ...func(int) error) (err error) {
// 	_, e := v.dbinfo.ColumnInfo(table, column)

// 	// column not found
// 	v.handleNotFoundColumn(table, column, e)

// 	for _, callback := range callbacks {
// 		if err = callback(value); err != nil {
// 			v.AddFieldError(column, err)
// 			return
// 		}
// 	}

// 	return
// }

// func (v *Validator) ValidateStringColumn(table, column string, value string, callbacks ...func(string) error) (err error) {
// 	_, e := v.dbinfo.ColumnInfo(table, column)
// 	v.handleNotFoundColumn(table, column, e)

// 	for _, callback := range callbacks {
// 		if err = callback(value); err != nil {
// 			v.AddFieldError(column, err)
// 			return
// 		}
// 	}

// 	return
// }
