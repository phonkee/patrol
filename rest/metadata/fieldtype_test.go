package metadata

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldType(t *testing.T) {

	Convey("Test kinds", t, func() {

		i, i8, i16, i32, i64 := int(1), int8(1), int16(1), int32(32), int64(64)
		u, u8, u16, u32, u64 := uint(1), uint8(1), uint16(1), uint32(1), uint64(1)
		f32, f64 := float32(32), float64(64)
		s := "string"
		b := true

		data := []struct {
			values    []interface{}
			fieldtype string
			required  bool
		}{
			{[]interface{}{i, i8, i16, i32, i64}, IntegerField, true},
			{[]interface{}{u, u8, u16, u32, u64}, UnsignedIntegerField, true},
			{[]interface{}{s}, StringField, true},
			{[]interface{}{b}, BoolField, true},
			{[]interface{}{f32, f64}, FloatField, true},
			// pointers (required is false)
			{[]interface{}{&i, &i8, &i16, &i32, &i64}, IntegerField, false},
			{[]interface{}{&u, &u8, &u16, &u32, &u64}, UnsignedIntegerField, false},
			{[]interface{}{&s}, StringField, false},
			{[]interface{}{&b}, BoolField, false},
			{[]interface{}{&f32, &f64}, FloatField, false},
		}

		for _, item := range data {
			for _, value := range item.values {
				f := GetField(reflect.TypeOf(value))
				So(f.Type, ShouldEqual, item.fieldtype)
				So(f.Required, ShouldEqual, item.required)
			}
		}
	})

	Convey("Test struct", t, func() {
		type Profile struct {
			Status int `json:"status"`
			Other  int `json:"other"`
			// Datetime time.Time `json:"datetime"`
		}

		type Product struct {
			ID        int        `json:"id"`
			Name      *string    `json:"name"`
			Something string     `json:"-"`
			Profile   *Profile   `json:"profile"`
			Start     time.Time  `json:"start"`
			End       *time.Time `json:"end"`
		}
		field := GetField(reflect.TypeOf(&Product{}))
		So(field.HasField("profile", "status"), ShouldBeTrue)
		So(field.HasField("profile", "non"), ShouldBeFalse)

		// silence
		_ = field.Pretty()
	})

	Convey("Test array", t, func() {
		type Test struct {
		}
		type Profile struct {
			Tests []int `json:"tests"`
		}

		x := GetField(reflect.TypeOf(Profile{}))
		So(x.Required, ShouldBeTrue)
		So(x.HasField("tests", "value"), ShouldBeTrue)
	})

	Convey("Test map", t, func() {
		mm := map[int]string{}
		x := GetField(reflect.TypeOf(mm))

		So(x.Field("key").Type, ShouldEqual, IntegerField)
		So(x.Field("value").Type, ShouldEqual, StringField)
	})

	Convey("Test custom types", t, func() {
		tt := "customtype"

		type CustomStruct struct{}
		RegisterType(ftf(tt), CustomStruct{})

		x := GetField(reflect.TypeOf(CustomStruct{}))
		So(x.Type, ShouldEqual, tt)
		So(x.Required, ShouldBeTrue)

		x = GetField(reflect.TypeOf(&CustomStruct{}))
		So(x.Type, ShouldEqual, tt)
		So(x.Required, ShouldBeFalse)

		type CustomStruct2 struct{}
		ttt := "adsfasfsdfs"
		RegisterType(ftf(ttt), &CustomStruct2{})
		x = GetField(reflect.TypeOf(&CustomStruct2{}))
		So(x.Type, ShouldEqual, ttt)

	})

}
