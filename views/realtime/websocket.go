package realtime

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func NewWebsocketAPIView(g func(*models.User, *http.Request) []string) *WebsocketAPIView {
	return &WebsocketAPIView{
		getsubs: g,
	}
}

type WebsocketAPIView struct {
	core.JSONView

	context *context.Context

	// function that returns all subscribed channels
	getsubs func(*models.User, *http.Request) []string

	user *models.User
}

func (v *WebsocketAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	v.context = v.Context(r)

	user := models.NewUser()
	if e := user.Manager(v.context).GetAuthUser(user, r); e == nil {
		v.user = user
	}

	return
}

func (v *WebsocketAPIView) GET(w http.ResponseWriter, r *http.Request) {
	subs := v.getsubs(v.user, r)
	glog.V(2).Infof("this is what we will listen to %+v\n", subs)

	var (
		ws  *websocket.Conn
		err error
	)
	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	go v.handleWebsocket(ws)
}

func (v *WebsocketAPIView) handleWebsocket(ws *websocket.Conn) {
	// subscribe on all needed queues
	// send all messages to client
	// profit

	glog.Infof("this is user %+v\n", v.user)

	for {
		time.Sleep(time.Second)
		ws.WriteJSON(v.user)
	}

}
