package validator

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResult(t *testing.T) {

	Convey("Test Add Field Error", t, func() {
		result := NewResult()
		result.AddFieldError("test", errors.New("error"))
		So(len(result.GetFieldErrors("test")), ShouldEqual, 1)
		result.AddFieldError("test", errors.New("other"))
		So(len(result.GetFieldErrors("test")), ShouldEqual, 2)

		result.AddFieldError("test2", errors.New("error"))
		So(len(result.GetFieldErrors("test2")), ShouldEqual, 1)
	})

	Convey("Test Add Duplicate Field Error", t, func() {
		result := NewResult()
		result.AddFieldError("test", errors.New("error"))
		result.AddFieldError("test", errors.New("error"))
		So(len(result.GetFieldErrors("test")), ShouldEqual, 1)
	})

	Convey("Test Non Existing Field errors", t, func() {
		result := NewResult()
		So(len(result.GetFieldErrors("test")), ShouldEqual, 0)
	})

	Convey("Test Has Field Errors", t, func() {
		result := NewResult()
		result.AddFieldError("test", errors.New("error"))
		So(result.HasFieldErrors("test"), ShouldBeTrue)
		So(result.HasFieldErrors("nonexisting"), ShouldBeFalse)
	})

	Convey("Test Exclude", t, func() {
		result := NewResult()
		result.AddFieldError("test", errors.New("error"))
		So(result.HasFieldErrors("test"), ShouldBeTrue)
		So(result.Exclude("test").HasFieldErrors("test"), ShouldBeFalse)
		So(result.Exclude("nonexisting").HasFieldErrors("test"), ShouldBeTrue)
	})

	Convey("Test Exclude", t, func() {
		result := NewResult()
		result.AddFieldError("test", errors.New("error"))
		So(result.HasFieldErrors("test"), ShouldBeTrue)
		So(result.Allow("test").HasFieldErrors("test"), ShouldBeTrue)
		So(result.Allow("nonexisting").HasFieldErrors("test"), ShouldBeFalse)
	})

	Convey("Test Append", t, func() {
		result := NewResult()
		result.AddFieldError("test", errors.New("error"))
		So(result.HasFieldErrors("test"), ShouldBeTrue)
		So(result.HasFieldErrors("nonexisting"), ShouldBeFalse)

		test := NewResult()
		test.Append(result)
		So(test.HasFieldErrors("test"), ShouldBeTrue)
		So(test.HasFieldErrors("nonexisting"), ShouldBeFalse)

	})

}
