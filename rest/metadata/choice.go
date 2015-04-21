package metadata

func NewChoices() *Choices {
	return &Choices{}
}

// type choices
type Choices []*Choice

// adds choice
func (c *Choices) Add(value interface{}, display string) *Choices {
	*c = append(*c, NewChoice(value, display))
	return c
}

func (c *Choices) Flush() *Choices {
	*c = *NewChoices()
	return c
}

// returns length of choices
func (c *Choices) Len() int {
	return len(*c)
}

// returns new choice
func NewChoice(value interface{}, display string) *Choice {
	return &Choice{
		Value:   value,
		Display: display,
	}
}

// choice type
type Choice struct {
	Value   interface{} `json:"value"`
	Display string      `json:"display"`
}
