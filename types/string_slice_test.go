package types

import (
	"testing"

	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStringSlice(t *testing.T) {
	context, err := utils.NewTestPatrolContext()
	_ = context
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	Convey("test add/index", t, func() {
		s := StringSlice{}
		s.Add("test")
		So(len(s), ShouldEqual, 1)
		s.AddUnique("test")
		So(len(s), ShouldEqual, 1)

		s.AddUnique("test2")
		So(s.Index("test2"), ShouldEqual, 1)
		So(s.Index("test"), ShouldEqual, 0)
	})

	Convey("test remove", t, func() {
		s := StringSlice{}
		s.Add("test")
		s.Add("test2")
		s.Add("test3")

		s.Remove("test2")
		So(len(s), ShouldEqual, 2)
	})

	Convey("test compact", t, func() {
		s := StringSlice{}
		s.Add("test")
		s.Add("test2")
		s.Add("test3")
		s.Add("test")
		s.Add("test")

		s.Compact()
		So(len(s), ShouldEqual, 2)
	})

}
