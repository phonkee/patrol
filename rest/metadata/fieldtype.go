package metadata

import (
	"reflect"
	"time"
)

const (
	InvalidField         = "invalid"
	ArrayField           = "array"
	StructField          = "struct"
	MapField             = "map"
	IntegerField         = "integer"
	UnsignedIntegerField = "unsigned"
	StringField          = "string"
	BoolField            = "boolean"
	FloatField           = "float"
	DateTimeField        = "datetime"
)

// field type func returns Field by reflect value
type FieldTypeFunc func(reflect.Type) *Field

var (
	kinds = map[reflect.Kind]FieldTypeFunc{}

	// mapping of custom types
	types = map[reflect.Type]FieldTypeFunc{}
)

// register kinds
func RegisterKind(f FieldTypeFunc, kind ...reflect.Kind) {
	for _, k := range kind {
		kinds[k] = f
	}
}

// RegisterType provides register for custom types (e.g. time.Time)
func RegisterType(f FieldTypeFunc, values ...interface{}) (err error) {
	for _, val := range values {
		typ := reflect.TypeOf(val)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		types[typ] = f
	}
	return
}

// returns field by kind
func GetFieldByKind(typ reflect.Type) (field *Field) {

	typn := typ
	if typn.Kind() == reflect.Ptr {
		typn = typ.Elem()
	}

	if fn, ok := kinds[typn.Kind()]; ok {
		return fn(typ)
	}

	// if something is not implemented
	return NewField().SetType(InvalidField)
}

// returns field by value
func GetField(typ reflect.Type) (field *Field) {

	orig := typ
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// custom types
	if fn, ok := types[typ]; ok {
		return fn(orig)
	}

	return GetFieldByKind(orig)
}

func init() {
	// register kinds
	RegisterKind(ftf(IntegerField), reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
	RegisterKind(ftf(UnsignedIntegerField), reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)
	RegisterKind(ftf(StringField), reflect.String)
	RegisterKind(ftf(BoolField), reflect.Bool)
	RegisterKind(ftf(FloatField), reflect.Float32, reflect.Float64)

	// special types
	RegisterKind(ftfStruct, reflect.Struct)
	RegisterKind(ftfArray, reflect.Array, reflect.Slice)
	RegisterKind(ftfMap, reflect.Map)
	RegisterType(ftf(DateTimeField), time.Now())
}

// default ftf imlpementation
func ftf(fieldtype string) FieldTypeFunc {
	return func(typ reflect.Type) (result *Field) {
		result = NewField().SetType(fieldtype)
		if typ.Kind() == reflect.Ptr {
			result.SetRequired(false)
		}
		return
	}
}

// field type function for struct
func ftfStruct(typ reflect.Type) (result *Field) {
	result = NewField().SetType(StructField)

	if typ.Kind() == reflect.Ptr {
		result.SetRequired(false)
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i)
		tag, _ := ParseTag(ft.Tag.Get("json"))
		if tag == "-" {
			continue
		}
		name := ft.Name
		if tag != "" {
			name = tag
		}

		result.AddField(name, GetField(ft.Type))
	}

	return

}

// array type
func ftfArray(typ reflect.Type) (result *Field) {
	result = NewField().SetType(ArrayField)
	result.AddField("value", GetField(typ.Elem()))
	return
}

// map type
func ftfMap(typ reflect.Type) (result *Field) {
	result = NewField().SetType(MapField)
	result.AddField("key", GetField(typ.Key()))
	result.AddField("value", GetField(typ.Elem()))
	return
}
