/*
	md := New("name").Description("description")
*/
package metadata

import "encoding/json"

// creates new metadata
func New(name string) *Metadata {
	md := &Metadata{
		Actions: map[string]Action{},
	}
	return md.SetName(name)
}

type Metadata struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Actions     map[string]Action `json:"actions"`
}

// if action does not exist it create new one and resturns
func (m *Metadata) Action(method string) (action Action) {
	method = cleanMethodName(method)

	var ok bool

	if action, ok = m.Actions[method]; ok {
		return
	}
	action = NewAction()
	m.Actions[method] = action
	return
}

// Aliases for common actions
func (m *Metadata) ActionCreate() Action {
	return m.Action("POST")
}

func (m *Metadata) ActionUpdate() Action {
	return m.Action("PUT")
}

func (m *Metadata) ActionRetrieve() Action {
	return m.Action("GET")
}

func (m *Metadata) ActionDelete() Action {
	return m.Action("DELETE")
}

// sets description and returns metadata for chaining calls
func (m *Metadata) SetDescription(description string) *Metadata {
	m.Description = description
	return m
}

// returns whether action is there
func (m *Metadata) HasAction(method string) bool {
	_, ok := m.Actions[cleanMethodName(method)]
	return ok
}

// returns available methods
func (m *Metadata) Methods() (result []string) {
	result = make([]string, 0, len(m.Actions))
	for k := range m.Actions {
		result = append(result, k)
	}
	return
}

// removes action
func (m *Metadata) RemoveAction(method string) *Metadata {
	delete(m.Actions, cleanMethodName(method))
	return m
}

// sets name and returns metadata for chaining calls
func (m *Metadata) SetName(name string) *Metadata {
	m.Name = name
	return m
}

// Returns marshalled
func (m *Metadata) String() string {
	result, _ := json.Marshal(m)
	return string(result)
}

func (m *Metadata) Pretty() string {
	result, _ := json.MarshalIndent(m, "", "    ")
	return string(result)
}
