package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type ExampleEventParser struct{ EventParser }

func (e *ExampleEventParser) Parse(body []byte) (result []*RawEvent, err error) { return }

func TestEventParserRegistry(t *testing.T) {

	Convey("TestRegister", t, func() {
		var err error

		registry := NewEventParserRegistry()
		err = registry.Register("4", func() EventParserer { return &ExampleEventParser{} })
		So(err, ShouldBeNil)

		err = registry.Register("5", func() EventParserer { return &ExampleEventParser{} })
		So(err, ShouldBeNil)

		err = registry.Register("4", func() EventParserer { return &ExampleEventParser{} })
		So(err, ShouldEqual, ErrEventParserAlreadyRegistered)

	})

}
