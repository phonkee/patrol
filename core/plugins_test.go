package core

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Plugin system tests

const (
	TEST_PLUGIN_ID = "example-plugin"
)

// Custom made dummy plugin
type ExamplePlugin struct {
	Plugin
}

func (e *ExamplePlugin) Init() error {
	return nil
}

func (e *ExamplePlugin) Id() string {
	return TEST_PLUGIN_ID
}

func TestPlugin(t *testing.T) {
	Convey("test example plugin", t, func() {
		p := &ExamplePlugin{}
		So(p.ID(), ShouldEqual, TEST_PLUGIN_ID)
		So(p.Commands(), ShouldResemble, []Commander{})
		So(p.URLViews(), ShouldResemble, []*URLView{})
		So(p.Init(), ShouldBeNil)
	})
}

func TestPluginRegistry(t *testing.T) {
	pr := NewPluginRegistry()
	p := &ExamplePlugin{}

	Convey("test plugin registry", t, func() {
		err := pr.RegisterPlugin(p)
		So(err, ShouldBeNil)

		errRe := pr.RegisterPlugin(p)
		So(errRe, ShouldNotBeNil)

		r, e := pr.Plugin(p.ID())
		So(e, ShouldBeNil)
		So(r, ShouldResemble, p)

		rn, en := pr.Plugin("non-existing-plugin-id")
		So(en, ShouldEqual, ErrPluginNotFound)
		So(rn, ShouldBeNil)

	})

	Convey("test plugin registry do method", t, func() {
		pr.RegisterPlugin(p)
		err := pr.Do(func(p Pluginer) error {
			return nil
		})
		So(err, ShouldBeNil)
	})

	Convey("test plugin registry do method with error", t, func() {
		pr.RegisterPlugin(p)
		err := pr.Do(func(p Pluginer) error {
			return fmt.Errorf("some error")
		})
		So(err, ShouldNotBeNil)
	})
}
