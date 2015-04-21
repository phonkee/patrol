package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChoices(t *testing.T) {

	Convey("Test new choice", t, func() {
		value, display := 1, "display"
		choice := NewChoice(value, display)
		So(choice.Value, ShouldEqual, value)
		So(choice.Display, ShouldEqual, display)
	})

	Convey("Test choices", t, func() {
		choices := NewChoices()

		So(choices.Len(), ShouldBeZeroValue)

		choices.Add("value", "display")
		So(choices.Len(), ShouldEqual, 1)

		choices.Flush()
		So(choices.Len(), ShouldBeZeroValue)

	})

}
