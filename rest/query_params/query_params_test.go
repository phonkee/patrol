package query_params

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryParams(t *testing.T) {
	Convey("test int", t, func() {

		data := []struct {
			input    string
			def      int
			expected int
		}{
			{"123", -1, 123},
			{"invalid", 5, 5},
		}

		varname := "var"

		for _, item := range data {
			uv := make(url.Values)
			uv.Set(varname, item.input)

			qp := New(uv)
			So(qp.GetInt(varname, item.def), ShouldEqual, item.expected)
		}

	})

	Convey("test bool", t, func() {

		data := []struct {
			input    string
			def      bool
			expected bool
		}{
			{"t", false, true},
			{"true", false, true},
			{"on", false, true},
			{"T", false, true},
			{"True", false, true},
			{"On", false, true},
			{"1", false, true},

			{"f", true, false},
			{"false", true, false},
			{"off", true, false},
			{"0", true, false},
			{"False", true, false},
			{"F", true, false},

			// test def
			{"", true, true},
		}

		varname := "var"

		for _, item := range data {
			uv := make(url.Values)
			uv.Set(varname, item.input)

			qp := New(uv)
			So(qp.GetBool(varname, item.def), ShouldEqual, item.expected)
		}

	})

	Convey("test float", t, func() {

		data := []struct {
			input    string
			def      float64
			expected float64
		}{
			{"0.1", 0, 0.1},
			{"-0.1", 0, -0.1},
			{"invalid", 3.1, 3.1},
		}

		varname := "var"

		for _, item := range data {
			uv := make(url.Values)
			uv.Set(varname, item.input)

			qp := New(uv)
			So(qp.GetFloat(varname, item.def), ShouldEqual, item.expected)
		}

	})

}
