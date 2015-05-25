package plugins

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/signals"
	"github.com/phonkee/patrol/views/realtime"
)

const (
	REALTIME_PLUGIN_ID = "realtime"
	ROUTE_REALTIME     = "websocket"
)

func NewRealtimePlugin(ctx *context.Context, pr *core.PluginRegistry) *RealtimePlugin {
	return &RealtimePlugin{
		context: ctx,
		pr:      pr,
		onrtw:   []signals.OnRealtimeWebsocketSubscribeToSignalHandler{},
	}
}

type RealtimePlugin struct {
	core.Plugin

	context *context.Context
	pr      *core.PluginRegistry
	onrtw   []signals.OnRealtimeWebsocketSubscribeToSignalHandler
}

func (r *RealtimePlugin) ID() string {
	return REALTIME_PLUGIN_ID
}

func (r *RealtimePlugin) Init() (err error) {
	// register signal to return subscribe queues
	r.pr.Do(func(plugin core.Pluginer) error {
		if t, ok := plugin.(signals.OnRealtimeWebsocketSubscribeToSignalHandler); ok {
			glog.V(2).Infof("realtime signals: adding %T as OnRealtimeWebsocketSubscribeToSignalHandler.", plugin)
			r.onrtw = append(r.onrtw, t)
		}
		return nil
	})

	return
}

// list of urls
func (r *RealtimePlugin) URLs() []*views.URL {
	return []*views.URL{
		views.NewURL(
			"/api/realtime/websocket",
			func() views.Viewer {
				return realtime.NewWebsocketAPIView(r.getSubscribeQueues)
			},
		).Name(ROUTE_REALTIME),
	}
}

func (r *RealtimePlugin) getSubscribeQueues(user *models.User, request *http.Request) (result []string) {
	result = []string{}
	for _, item := range r.onrtw {
		for _, queue := range item.OnRealtimeWebsocketSubscribeTo(user, request) {
			result = append(result, queue)
		}
	}
	return
}

func (r *RealtimePlugin) OnRealtimeWebsocketSubscribeTo(u *models.User, req *http.Request) (result []string) {
	result = []string{}

	if u != nil {
		result = append(result, "reload-websocket"+u.ID.String())
	}

	return
}
