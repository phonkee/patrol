package utils

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDBArrayPlaceholder(t *testing.T) {
	Convey("test DBArrayPlaceholder", t, func() {
		ph := ""
		ph = DBArrayPlaceholder(3)
		So(ph, ShouldEqual, "?,?,?")

		ph = DBArrayPlaceholder(0)
		So(ph, ShouldEqual, "")
	})
}

func TestStringIndex(t *testing.T) {

	words := "zero one two three four five six seven eight nine ten"

	Convey("test string index from start", t, func() {
		data := []struct {
			list     []string
			s        string
			position int
		}{
			{
				strings.Split(words, " "),
				"three",
				3,
			},
		}

		for _, item := range data {
			So(StringIndex(item.list, item.s), ShouldEqual, item.position)
		}
	})

	Convey("test string index from custom start", t, func() {
		data := []struct {
			list     []string
			start    int
			s        string
			position int
		}{
			{strings.Split(words, " "), 4, "nine", 9},
			{strings.Split(words, " "), -2, "eight", -1},
			{strings.Split(words, " "), -2, "nine", 9},
		}

		for _, item := range data {
			si := StringIndex(item.list, item.s, item.start)
			So(si, ShouldEqual, item.position)
		}
	})

	Convey("test split identifier", t, func() {
		data := []struct {
			identifier  string
			def         string
			migrationId string
			pluginId    string
			shouldError bool
		}{
			{"auth:initial-migration", "", "initial-migration", "auth", false},
			{"auth:initial:migration", "", "initial:migration", "auth", false},
			{":initial:migration", "asdf", "initial:migration", "asdf", false},
			{"initial-migration", "asdf", "initial-migration", "asdf", false},
			{"asdf:", "", "", "", true},
			{":", "", "", "", true},
			{"", "", "", "", true},
		}

		for _, item := range data {
			a, b, err := SplitIdentifier(item.identifier, item.def)
			if item.shouldError {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(b, ShouldEqual, item.migrationId)
				So(a, ShouldEqual, item.pluginId)
			}
		}

	})
}

func TestStringPadLeft(t *testing.T) {
	Convey("test string pad left", t, func() {
		data := []struct {
			input  string
			char   string
			num    int
			output string
		}{
			{"test", " ", 10, "      test"},
			{"testing", " ", 2, "testing"},
		}

		for _, item := range data {
			result := StringPadLeft(item.input, item.char, item.num)
			So(result, ShouldEqual, item.output)
		}

	})
}

func TestStringPadRight(t *testing.T) {
	Convey("test string pad right", t, func() {
		data := []struct {
			input  string
			char   string
			num    int
			output string
		}{
			{"test", "a", 10, "testaaaaaa"},
			{"testing", "a", 2, "testing"},
		}

		for _, item := range data {
			result := StringPadRight(item.input, item.char, item.num)
			So(result, ShouldEqual, item.output)
		}
	})
}

func TestRandomString(t *testing.T) {
	Convey("test random string", t, func() {

		s := RandomString(10)
		So(len(s), ShouldEqual, 10)

		s2 := RandomString(15, "abc")
		So(len(s2), ShouldEqual, 15)

		for _, c := range s2 {
			So(strings.ContainsRune("abc", c), ShouldBeTrue)
		}

	})
}
