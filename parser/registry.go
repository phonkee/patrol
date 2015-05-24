package parser

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/golang/glog"
)

var (
	Registry = NewEventParserRegistry()
)

func Versions() []string {
	return Registry.Versions()
}

func Register(version string, f EventParserFunc) error {
	return Registry.Register(version, f)
}

func Parse(body []byte, version string) ([]*RawEvent, error) { return Registry.Parse(body, version) }
func ParseRequest(r *http.Request, version string) ([]*RawEvent, error) {
	return Registry.ParseRequest(r, version)
}

func NewEventParserRegistry() *EventParserRegistry {
	return &EventParserRegistry{
		parsers: map[string]EventParserFunc{},
	}
}

type EventParserRegistry struct {
	parsers map[string]EventParserFunc
	mutex   sync.RWMutex
}

func (e *EventParserRegistry) Register(version string, f EventParserFunc) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.parsers[version]; ok {
		return ErrEventParserAlreadyRegistered
	}
	e.parsers[version] = f
	return nil
}

func (e *EventParserRegistry) GetEventParser(version string) (EventParserer, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	var (
		epf EventParserFunc
		ok  bool
	)
	if epf, ok = e.parsers[version]; !ok {
		return nil, ErrEventParserNotFound
	}
	return epf(), nil
}

/*
Parser http request
*/
func (e *EventParserRegistry) ParseRequest(r *http.Request, version string) (result []*RawEvent, err error) {

	var parser EventParserer

	if parser, err = e.GetEventParser(version); err != nil {
		return
	}

	var body []byte

	if body, err = parser.DecodeRequest(r); err != nil {
		return
	}

	return e.Parse(body, version)
}

func (e *EventParserRegistry) Parse(body []byte, version string) (result []*RawEvent, err error) {

	var ep EventParserer

	if ep, err = e.GetEventParser(version); err != nil {
		return
	}

	if result, err = ep.Parse(body); err != nil {
		return
	}

	return
}

func (e *EventParserRegistry) Versions() []string {
	result := make([]string, 0, len(e.parsers))
	for k := range Registry.parsers {
		result = append(result, k)
	}
	return result
}

/*
EventInterfaceParser registry
*/
func NewEventInterfaceParserRegistry() *EventInterfaceParserRegistry {
	return &EventInterfaceParserRegistry{
		interfaces: []*EventInterfaceParserRegistryItem{},
	}
}

type EventInterfaceParserRegistry struct {
	interfaces []*EventInterfaceParserRegistryItem
}

func (e *EventInterfaceParserRegistry) Len() int {
	return len(e.interfaces)
}
func (e *EventInterfaceParserRegistry) Swap(i, j int) {
	e.interfaces[i], e.interfaces[j] = e.interfaces[j], e.interfaces[i]
}
func (e *EventInterfaceParserRegistry) Less(i, j int) bool {
	return e.interfaces[i].score < e.interfaces[j].score
}

func (e *EventInterfaceParserRegistry) Register(f EventParserInterfacerFunc, id string, aliases []string, score int) error {
	e.interfaces = append(e.interfaces, &EventInterfaceParserRegistryItem{
		f:       f,
		aliases: aliases,
		id:      id,
		score:   score,
	})
	// reverse sort by score

	sort.Sort(sort.Reverse(e))

	return nil
}

func (e *EventInterfaceParserRegistry) Each(f func(item *EventInterfaceParserRegistryItem)) {
	for _, i := range e.interfaces {
		f(i)
	}
	return
}

func (e *EventInterfaceParserRegistry) Parse(values map[string]json.RawMessage) ([]EventParserInterfacer, error) {

	result := []EventParserInterfacer{}
	for _, iface := range e.interfaces {

		keys := append([]string{iface.id}, iface.aliases...)

		for _, key := range keys {
			if value, ok := values[key]; ok {
				tmp := iface.f()
				tmp.SetID(iface.id)
				tmp.SetScore(iface.score)
				delete(values, key)

				if err := json.Unmarshal(value, &tmp); err != nil {
					glog.Errorf("this is something %+v\n", err)
				}

				result = append(result, tmp)
			}
			continue
		}
	}
	return result, nil
}

type EventInterfaceParserRegistryItem struct {
	aliases []string
	f       EventParserInterfacerFunc
	id      string
	score   int
}
