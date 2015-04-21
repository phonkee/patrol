package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFields(t *testing.T) {

	Convey("Test new field with simple setters", t, func() {
		field := NewField()

		label := "this is label"
		field.SetLabel(label)
		So(field.Label, ShouldEqual, label)

		field.SetRequired(true)
		So(field.Required, ShouldEqual, true)

		t := "mytype"
		field.SetType(t)
		So(field.Type, ShouldEqual, t)

		helptext := "Lorem ipsum something..."
		field.SetHelpText(helptext)
		So(field.HelpText, ShouldEqual, helptext)

	})

	Convey("Test new fields in field", t, func() {
		field := NewField()
		name := "myfield"
		So(field.HasField(name), ShouldBeFalse)

		// create subfield
		f1 := field.Field(name)
		So(field.HasField(name), ShouldBeTrue)

		// returns existing
		f2 := field.Field(name)
		So(f1, ShouldPointTo, f2)
	})

}
