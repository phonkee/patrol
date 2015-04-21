package models

import "github.com/phonkee/patrol/types"

/*
RealtimeMessage is interface for sending messages to websockets
*/
type RealtimeMessage interface {
	Identifier() string
	Type() string
}

type UserLoggedRealtimeMessage struct {
	ID       types.PrimaryKey `json:"id"`
	Username string           `json:"username"`
}

func (u *UserLoggedRealtimeMessage) Identifier() string {
	return "auth:user_logged:" + u.ID.String()
}
