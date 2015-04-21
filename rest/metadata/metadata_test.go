package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMetadata(t *testing.T) {

	Convey("Test set name", t, func() {
		name := "this is name"
		md := New(name)

		So(md.Name, ShouldEqual, name)
	})

	Convey("Test set description", t, func() {
		description := "my awesome metadata"
		md := New("name").SetDescription(description)
		So(md.Description, ShouldEqual, description)
	})

	Convey("Test actions", t, func() {
		md := New("name")

		action := md.Action("GET")
		other := md.Action("GET")

		So(action, ShouldEqual, other)
		So(md.Methods(), ShouldResemble, []string{"GET"})
		So(md.HasAction("get"), ShouldBeTrue)

		// test remove action
		md.RemoveAction("get")
		So(md.HasAction("get"), ShouldBeFalse)
	})

	Convey("Test actions fields", t, func() {
		md := New("name")
		md.Action("GET").Field("myfield").SetRequired(true).SetLabel("label")

		// silence
		_ = md.String()
		_ = md.Pretty()
	})

	Convey("Test aliases", t, func() {
		var md *Metadata

		md = New("md")
		md.ActionCreate()
		So(md.HasAction("POST"), ShouldBeTrue)

		md = New("md")
		md.ActionUpdate()
		So(md.HasAction("PUT"), ShouldBeTrue)

		md = New("md")
		md.ActionRetrieve()
		So(md.HasAction("GET"), ShouldBeTrue)

		md = New("md")
		md.ActionDelete()
		So(md.HasAction("DELETE"), ShouldBeTrue)
	})

}
