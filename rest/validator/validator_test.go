package validator

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidator(t *testing.T) {

	Convey("Test New Validator", t, func() {
		validator := New()
		So(validator, ShouldNotBeNil)
	})

	Convey("Test Validator add func", t, func() {
		validator := New()

		vf := func(interface{}) error { return nil }

		validator.Add("field", vf)
		So(len(validator), ShouldEqual, 1)

	})

	Convey("Test Validate", t, func() {

		type User struct {
			Username string `json:"username" validator:"some,other,etc"`
		}

		alwayserr := func(interface{}) error { return errors.New("hehe") }

		validator := New().Add("some", alwayserr)
		result := validator.Validate(&User{})
		fmt.Printf("this is result %+v\n", result)

	})

}
