package parser

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
)

var (
	interfacesV4 = NewEventInterfaceParserRegistry()
)

func init() {
	// register parser
	GetV4 := func() EventParserer {
		return &EventParserV4{}
	}
	Register(settings.EVENT_PARSER_PROTOCOL_V4, GetV4)

	// register interfaces parsers
	interfacesV4.Register(
		func() EventParserInterfacer { return &HttpInterfaceV4{} },
		"http", []string{"sentry.interfaces.Http"}, // id + aliases
		800, //score
	)
	interfacesV4.Register(
		func() EventParserInterfacer { return &ExceptionInterfaceV4{} },
		"exception", []string{"sentry.interfaces.Exception"}, // id + aliases
		900, //score
	)
}

/*
EventParserV4 - parser that implements sentry protocol version 4
*/
type EventParserV4 struct {
	EventParser
}

func (e *EventParserV4) EventInterfaceParserRegistry() *EventInterfaceParserRegistry {
	return interfacesV4
}

func (e *EventParserV4) Parse(body []byte) (events []*RawEvent, err error) {

	// fmt.Printf("body is this \n\n%+v\n\n", string(body))

	event := NewRawEvent()
	events = []*RawEvent{}

	values := map[string]json.RawMessage{}
	if err = json.Unmarshal(body, &values); err != nil {
		return
	}

	// parser functions
	ufs := map[string]func(value json.RawMessage) error{
		"event_id":    func(value json.RawMessage) error { return json.Unmarshal(value, &event.EventID) },
		"message":     func(value json.RawMessage) error { return json.Unmarshal(value, &event.Message) },
		"level":       func(value json.RawMessage) error { return json.Unmarshal(value, &event.Level) },
		"logger":      func(value json.RawMessage) error { return json.Unmarshal(value, &event.Logger) },
		"server_name": func(value json.RawMessage) error { return json.Unmarshal(value, &event.ServerName) },
		"culprit":     func(value json.RawMessage) error { return json.Unmarshal(value, &event.Culprit) },
		"platform":    func(value json.RawMessage) error { return json.Unmarshal(value, &event.Platform) },
		"project": func(value json.RawMessage) (err error) {
			var project string
			if err = json.Unmarshal(value, &project); err != nil {
				return
			}
			var pid int64
			if pid, err = strconv.ParseInt(project, 10, 0); err != nil {
				return
			}
			event.ProjectID = types.ForeignKey(pid)
			return
		},
		"tags": func(value json.RawMessage) (err error) {
			event.Tags = e.ParseTags(value)
			return
		},
		"timestamp": func(value json.RawMessage) (err error) {
			var timestamp string
			if err = json.Unmarshal(value, &timestamp); err != nil {
				return

			}
			event.Datetime, err = time.Parse(settings.SENTRY_TIMESTAMP_LAYOUT, timestamp)
			return
		},
		"extra": func(value json.RawMessage) (err error) {
			_ = json.Unmarshal(value, &event.Extra)
			return
		},
	}

	for key, f := range ufs {
		// if not found give blank RawMessage
		body, ok := values[key]
		if !ok {
			body = json.RawMessage{}
		}
		// if cannot parse field return error
		if err = f(body); err != nil {
			err = fmt.Errorf("cannot parse %s field, value %s", key, body)
			return
		}
		delete(values, key)
	}

	// add version
	event.Version = settings.EVENT_PARSER_PROTOCOL_V4

	var ifs []EventParserInterfacer
	if ifs, err = e.EventInterfaceParserRegistry().Parse(values); err != nil {
		return
	}

	// add iterfaces to data
	event.Data["interfaces"] = ifs

	// update checksum
	if len(ifs) > 0 {
		event.Checksum = ifs[0].Hash()
	} else {
		h := md5.New()
		io.WriteString(h, event.Message)
		event.Checksum = fmt.Sprintf("%x", h.Sum(nil))
	}

	// all other data will go to data
	for key, val := range values {
		if key == "interfaces" {
			continue
		}
		event.Data[key] = string(val)
	}

	events = append(events, event)

	return
}

/*
Parse tags from raw json message
*/
func (e *EventParserV4) ParseTags(body json.RawMessage) (tags map[string]string) {
	tags = make(map[string]string)

	var err error

	// first try to directly parse map[string]string
	if err = json.Unmarshal(body, &tags); err == nil {
		return
	}

	// first we try to parse tags as [][]string
	tagslist := [][]string{}
	if err = json.Unmarshal(body, &tagslist); err == nil {
		for _, item := range tagslist {
			if len(item) != 2 {
				continue
			}
			tags[item[0]] = item[1]
		}
		return
	}
	return
}

/*
Event Interfaces
*/
type HttpInterfaceV4 struct {
	PatrolInterface
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	QueryString string            `json:"query_string"`
	Headers     map[string]string `json:"headers'"`
	Env         map[string]string `json:"env'"`
}

/*
Returns hash (checksum) for http
*/
func (h *HttpInterfaceV4) Hash() string {
	hash := md5.New()
	io.WriteString(hash, h.URL)
	io.WriteString(hash, h.Method)
	io.WriteString(hash, h.QueryString)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (h *HttpInterfaceV4) String() string {
	return fmt.Sprintf("http %v with checksum %s", h, h.Hash())
}

func (h *HttpInterfaceV4) Template() string {
	return "this is {{ interface.URL }}"
}

/*
Exception
*/
type ExceptionInterfaceV4 struct {
	PatrolInterface
	Value      string                 `json:"value"`
	Type       string                 `json:"type"`
	Stacktrace *StacktraceInterfaceV4 `json:"stacktrace,omitempty"`
}

func (e *ExceptionInterfaceV4) Hash() string {

	if e.Stacktrace == nil {
		return ""
	}

	return e.Stacktrace.Hash()
}
func (e *ExceptionInterfaceV4) String() string   { return e.Hash() }
func (e *ExceptionInterfaceV4) Template() string { return "this is template for exception" }

// func (e *ExceptionInterfaceV4) UnmarshalJSON(body []byte) error {
// 	// fmt.Printf("Really?????????????? %s", string(body))
// 	return nil
// }

/*
Stacktrace
*/
type StacktraceInterfaceV4 struct {
	PatrolInterface
	Frames []StacktraceFrameV4 `json:"frames,omitempty"`
}

func (e *StacktraceInterfaceV4) Hash() string {
	h := md5.New()
	for _, frame := range e.Frames {
		io.WriteString(h, frame.Hash())
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
func (e *StacktraceInterfaceV4) String() string   { return "ehm" }
func (e *StacktraceInterfaceV4) Template() string { return "ehm" }

// func (e *StacktraceInterfaceV4) UnmarshalJSON(body []byte) error {
// 	fmt.Printf("Really?????????????? %s", string(body))
// 	return nil
// }

type StacktraceFrameV4 struct {
	Filename    string   `json:"filename"`
	Function    string   `json:"function"`
	Lineno      int      `json:"lineno,omitempty"`
	Module      string   `json:"module"`
	AbsPath     string   `json:"abs_path"`
	PreContext  []string `json:"pre_context"`
	ContextLine string   `json:"context_line,omitempty"`
	PostContext []string `json:"post_context"`
}

func (s *StacktraceFrameV4) IsURL() bool {
	return strings.Contains(s.Filename, "://")
}

func (s *StacktraceFrameV4) Hash() string {

	parts := []string{}

	if s.Module != "" {
		parts = append(parts, s.Module)
	}

	if s.Filename != "" && !s.IsURL() {
		parts = append(parts, s.Filename)
	}

	canUseContext := true

	if s.ContextLine == "" {
		canUseContext = false
	} else if len(s.ContextLine) > 128 {
		canUseContext = false
	} else if s.Function != "" && strings.HasPrefix(s.Function, "[Anonymous") {
		canUseContext = true
	} else {
		canUseContext = false
	}

	if canUseContext {
		parts = append(parts, s.ContextLine)
	} else if len(parts) == 0 {
		// If we were unable to achieve any context at this point
		// (likely due to a bad JavaScript error) we should just
		// bail on recording this frame
		return ""
	} else if s.Function != "" {
		parts = append(parts, s.Function)
	} else if s.Lineno != 0 {
		parts = append(parts, strconv.Itoa(s.Lineno))
	}

	h := md5.New()
	io.WriteString(h, strings.Join(parts, ":"))
	return fmt.Sprintf("%x", h.Sum(nil))
}
