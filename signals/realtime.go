package signals

import (
	"net/http"

	"github.com/phonkee/patrol/models"
)

/* OnRealtimeWebsocketSubscribeToSignalHandler
This signal is called on websocket and should return list of
*/
type OnRealtimeWebsocketSubscribeToSignalHandler interface {
	// when event request is received
	OnRealtimeWebsocketSubscribeTo(user *models.User, r *http.Request) []string
}
