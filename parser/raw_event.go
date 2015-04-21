package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"time"

	"github.com/phonkee/ergoq"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
)

/*
RawEvent is event that is parsed from request and is sent through wire(MQ) to
event worker.

This model is not saved to database it's just parsed data
*/
func NewRawEvent() *RawEvent {
	return &RawEvent{
		Tags: map[string]string{},
		Data: types.GzippedMap{},
	}
}

// Parsed event
type RawEvent struct {
	Culprit    string                 `json:"culprit"`
	Checksum   string                 `json:"checksum"`
	Datetime   time.Time              `json:"date_time"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
	EventID    string                 `json:"event_id"`
	Level      string                 `json:"level"`
	Logger     string                 `json:"logger"`
	Message    string                 `json:"message"`
	ProjectID  types.ForeignKey       `json:"project_id"`
	Platform   string                 `json:"platform"`
	Release    string                 `json:"release"`
	ServerName string                 `json:"server_name"`
	Version    string                 `json:"version"`
	Tags       map[string]string      `json:"tags"`
	Data       types.GzippedMap       `json:"data"`
}

/*
Manager

*/
type RawEventManager struct {
	context *context.Context
}

func NewRawEventManager(context *context.Context) *RawEventManager {
	return &RawEventManager{context: context}
}

func (rem *RawEventManager) NewRawEvent() *RawEvent {
	return &RawEvent{}
}

func (rem *RawEventManager) NewRawEventList() []*RawEvent {
	return []*RawEvent{}
}

/*
Push raw event to queue
*/
func (r *RawEventManager) PushRawEvent(e interface{}) (err error) {
	var body []byte
	if body, err = json.Marshal(e); err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	var writer *gzip.Writer
	if writer, err = gzip.NewWriterLevel(buffer, settings.RAW_EVENT_COMPRESSION_LEVEL); err != nil {
		return
	}
	writer.Write(body)
	writer.Close()
	return r.context.Queue.Push(settings.EVENT_QUEUE_ID, buffer.Bytes())
}

func (r *RawEventManager) PopRawEvent() (e *RawEvent, message ergoq.QueueMessage, err error) {
	if message, err = r.context.Queue.Pop(settings.EVENT_QUEUE_ID); err == nil {
		// add decompression here
		buffer := bytes.NewReader(message.Message())
		var gr *gzip.Reader
		if gr, err = gzip.NewReader(buffer); err != nil {
			return
		}
		gr.Close()

		result := new(bytes.Buffer)
		io.Copy(result, gr)

		e = r.NewRawEvent()
		err = json.Unmarshal(result.Bytes(), &e)
	}

	return
}
