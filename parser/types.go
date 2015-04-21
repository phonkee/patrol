/*
Package parser provides parsers for event message
Parsers can be registered by version string
*/
package parser

import (
	"compress/zlib"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
)

/*
event parser constructor function
returns new parser instance
*/
type EventParserFunc func() EventParserer

type EventParserInterfacerFunc func() EventParserInterfacer

/*
EventParserer interface
*/
type EventParserer interface {
	// decode body
	DecodeRequest(*http.Request) ([]byte, error)

	// parses body and returns events
	Parse(body []byte) ([]*RawEvent, error)

	// @TODO: every parser must return all interfaces it uses (for urls)
	EventInterfaceParserRegistry() *EventInterfaceParserRegistry
}

/*
Base EventParser
*/
type EventParser struct{}

// Decodes request
func (e *EventParser) DecodeRequest(r *http.Request) (body []byte, err error) {
	decoder := base64.NewDecoder(base64.StdEncoding, r.Body)
	var zlibr io.Reader
	if zlibr, err = zlib.NewReader(decoder); err != nil {
		return
	}
	body, err = ioutil.ReadAll(zlibr)
	return
}

/*
Return blank
*/
func (e *EventParser) EventInterfaceParserRegistry() *EventInterfaceParserRegistry {
	return NewEventInterfaceParserRegistry()
}

/*
EventParserInterfacer

interfaces as defined in sentry
*/
type EventParserInterfacer interface {
	// returns html template
	Template() string

	// return hash
	Hash() string

	// Sets id
	SetID(string)

	// returns id
	GetID() string

	// sets score
	SetScore(int)
}

type PatrolInterface struct {
	ID    string `json:"id,omitempty"`
	Score int    `json:"score"`
}

func (p *PatrolInterface) SetID(id string) {
	p.ID = id
}

func (p *PatrolInterface) GetID() string {
	return p.ID
}

func (p *PatrolInterface) SetScore(score int) {
	p.Score = score
}
