package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtils(t *testing.T) {

	Convey("Test cleanMethodName", t, func() {
		data := []struct {
			in  string
			out string
		}{
			{"get", "GET"},
			{" post ", "POST"},
		}

		for _, item := range data {
			So(cleanMethodName(item.in), ShouldEqual, item.out)
		}
	})

	Convey("Test parseTag", t, func() {
		parsed, options := ParseTag("-")
		So(parsed, ShouldEqual, "-")
		So(options.Contains("omitempty"), ShouldBeFalse)

		parsed, options = ParseTag(",omitempty")
		So(parsed, ShouldEqual, "")
		So(options.Contains("omitempty"), ShouldBeTrue)

		_, options = ParseTag(",omitempty,other")
		So(options.Contains("other"), ShouldBeTrue)
		So(options.Contains("nonexisting"), ShouldBeFalse)

	})

}
