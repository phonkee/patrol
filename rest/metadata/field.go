package metadata

import "encoding/json"

// returns new field
func NewField() *Field {
	return &Field{
		Fields:   map[string]*Field{},
		Choices:  Choices{},
		Required: true,
	}
}

// Field struct
type Field struct {
	HelpText string            `json:"help_text,omitempty"`
	Label    string            `json:"label,omitempty"`
	Min      *int              `json:"min,omitempty"`
	Max      *int              `json:"max,omitempty"`
	Required bool              `json:"required"`
	Type     string            `json:"type,omitempty"`
	Fields   map[string]*Field `json:"fields,omitempty"`
	Choices  Choices           `json:"choices,omitempty"`
}

/*
Possibility to add callbacks
*/
func (f *Field) Update(funcs ...func(field *Field)) *Field {
	for _, fn := range funcs {
		fn(f)
	}
	return f
}

// sets help text
func (f *Field) SetHelpText(ht string) *Field {
	f.HelpText = ht
	return f
}

// setter for label
func (f *Field) SetLabel(label string) *Field {
	f.Label = label
	return f
}

// sets required
func (f *Field) SetRequired(required bool) *Field {
	f.Required = required
	return f
}

func (f *Field) SetType(typ string) *Field {
	f.Type = typ
	return f
}

func (f *Field) SetMax(m int) *Field {
	f.Max = &m
	return f
}

func (f *Field) RemoveMax(m int) *Field {
	f.Max = nil
	return f
}

func (f *Field) SetMin(m int) *Field {
	f.Min = &m
	return f
}

func (f *Field) RemoveMin(m int) *Field {
	f.Min = nil
	return f
}

// adds field
func (f *Field) AddField(name string, field *Field) *Field {
	f.Fields[name] = field
	return field
}

// check if field has field
func (f *Field) HasField(name string, names ...string) bool {
	field, ok := f.Fields[name]
	if !ok {
		return ok
	}

	for _, n := range names {
		if !field.HasField(n) {
			return false
		}
		field = field.Field(n)
	}

	return true
}

// add sub field
func (f *Field) Field(name string) *Field {
	if field, ok := f.Fields[name]; ok {
		return field
	}
	field := NewField()
	f.Fields[name] = field
	return field
}

// pretty prints
func (f *Field) Pretty() string {
	result, _ := json.MarshalIndent(f, "", "    ")
	return string(result)
}
