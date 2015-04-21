package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestActions(t *testing.T) {

	Convey("Test create field", t, func() {
		action := NewAction()

		name := "field"
		action.Field(name)

		So(len(action), ShouldEqual, 1)
		action.Field(name)

		So(len(action), ShouldEqual, 1)

		So(action.HasField(name), ShouldBeTrue)
	})

	Convey("Test action.From", t, func() {
		type Product struct {
			Name  string   `json:"name"`
			Price *float64 `json:"price"`
		}

		action := NewAction().From(Product{})
		So(action.HasField("name"), ShouldBeTrue)
		So(action.HasField("price"), ShouldBeTrue)
		So(action.Field("name").Required, ShouldBeTrue)
		So(action.Field("price").Required, ShouldBeFalse)

		actionnew := NewAction().From(&Product{})
		So(actionnew.HasField("name"), ShouldBeTrue)
		So(actionnew.HasField("price"), ShouldBeTrue)
		So(actionnew.Field("name").Required, ShouldBeTrue)
		So(actionnew.Field("price").Required, ShouldBeFalse)

		So(func() { NewAction().From(1) }, ShouldPanic)

	})
}
