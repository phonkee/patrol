package models

import "github.com/phonkee/patrol/context"

const (
	MESSAGE_TYPE_REFRESH_SOCKET = "refresh-socket"
)

/*
Realtime manager handles publishing of messages to websocket
*/
type RealtimeManager struct {
	Manager
	context *context.Context
}

func NewRealtimeManager(c *context.Context) *RealtimeManager {
	return &RealtimeManager{
		context: c,
	}
}

//  Publishes message to queue
func (r *RealtimeManager) Publish(message RealtimeMessage) (err error) {
	return
}
