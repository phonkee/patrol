package metadata

import "reflect"

// NewAction acts as constructor and returns new action
func NewAction() Action {
	return Action{}
}

// action holds multiple fields
type Action map[string]*Field

func (a Action) Field(name string) *Field {
	if field, ok := a[name]; ok {
		return field
	}
	field := NewField()
	a[name] = field
	return field
}

// HasField checks if field exists
func (a Action) HasField(name string) bool {
	_, ok := a[name]
	return ok
}

// create action fields from target
func (a Action) From(v interface{}) Action {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic("Metadata.From supports only structs")
	}
	f := GetField(typ)
	for name, field := range f.Fields {
		a[name] = field
	}
	return a
}
