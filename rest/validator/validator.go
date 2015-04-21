package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/phonkee/patrol/rest/metadata"
)

const (
	VALIDATE_TAG = "validator"
)

type Cleaner interface {
	Clean()
}

// Validator
type Validator map[string]ValidatorFunc

func New() Validator {
	return Validator{}
}

// add validator function
func (v Validator) Add(name string, fn ValidatorFunc) Validator {
	v[name] = fn
	return v
}

// validates target and returns validator Result
func (v Validator) Validate(s interface{}) (result *Result) {
	if v, ok := s.(Cleaner); ok {
		v.Clean()
	}
	return v.validatePrefixed(s, "")
}

func (v Validator) validatePrefixed(s interface{}, prefix string) (result *Result) {
	result = NewResult()
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	t := val.Type()
	if t == nil || t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)
		if !fv.CanInterface() {
			continue
		}
		val := fv.Interface()
		tag := f.Tag.Get(VALIDATE_TAG)
		if tag == "" {
			continue
		}
		vts := strings.Split(tag, ",")

		for _, vt := range vts {
			name := f.Name
			tag, _ := metadata.ParseTag(f.Tag.Get("json"))
			if tag != "" && tag != "-" {
				name = tag
			}

			if len(prefix) > 0 {
				name = prefix + "." + name
			}

			if vt == "struct" {
				structResult := v.validatePrefixed(val, name)
				result.Append(structResult)
				continue
			}

			vf := v[vt]
			if vf == nil {
				result.AddUnboundError(fmt.Errorf("undefined validator: %q", vt))
				continue
			}
			if err := vf(val); err != nil {
				result.AddFieldError(name, err)
			}
		}

	}

	return
}
