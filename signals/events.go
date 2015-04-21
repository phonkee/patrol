package signals

import (
	"net/http"

	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/parser"
)

/* OnEventRequestSignalHandler
This signal is called when patrol receives request with event.
Signal handler can write it's own response end return error. If error is returned
event is not added to message queue with events to process.
With this handler it's really easy to implement event throttling. With combination
with external API it can be used to bill customers.
*/
type OnEventRequestSignalHandler interface {
	// when event request is received
	OnEventRequest(event *parser.RawEvent, rw http.ResponseWriter, r *http.Request) error
}

/* OnEventSignalHandler
This signal is called when event is processed by event worker
*/
type OnEventSignalHandler interface {
	OnEvent(event *models.Event, eventgroup *models.EventGroup)
}
