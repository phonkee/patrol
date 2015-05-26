package plugins

import (
	"github.com/golang/glog"
	"github.com/justinas/alice"
	"github.com/phonkee/patrol/commands"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/signals"
	"github.com/phonkee/patrol/views/events"
)

func NewEventsPlugin(context *context.Context, pr *core.PluginRegistry) core.Pluginer {
	return &EventsPlugin{context: context, pr: pr}
}

type EventsPlugin struct {
	core.Plugin
	context *context.Context
	pr      *core.PluginRegistry

	// signal handlers
	onEventRequestHandlers []signals.OnEventRequestSignalHandler
	onEventHandlers        []signals.OnEventSignalHandler
}

func (e *EventsPlugin) ID() string { return settings.EVENTS_PLUGIN_ID }
func (e *EventsPlugin) Init() (err error) {
	// add on event signal handlers
	if err = e.pr.Do(func(plugin core.Pluginer) error {
		if t, ok := plugin.(signals.OnEventSignalHandler); ok {
			glog.V(2).Infof("event signals: adding %T as OnEventHandler.", plugin)
			e.onEventHandlers = append(e.onEventHandlers, t)
		}
		return nil
	}); err != nil {
		return err
	}

	// add on event request signal handlers
	if err = e.pr.Do(func(plugin core.Pluginer) error {
		if t, ok := plugin.(signals.OnEventRequestSignalHandler); ok {
			glog.V(2).Infof("event signals: adding %T as OnEventRequestSignalHandler.", plugin)
			e.onEventRequestHandlers = append(e.onEventRequestHandlers, t)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// send OnEvent
func (e *EventsPlugin) SendOnEventSignal(event *models.Event, eventgroup *models.EventGroup) {
	for _, sh := range e.onEventHandlers {
		func() {
			sh.OnEvent(event, eventgroup)
			if err := recover(); err != nil {
				glog.Errorf("signal handler panicked %+v", err)
			}
		}()
	}
}

func (e *EventsPlugin) URLs() []*views.URL {
	mids := []alice.Constructor{
		middlewares.AuthTokenValidMiddleware(),
	}
	result := []*views.URL{
		views.NewURL("/api/{project_id:[0-9]+}/store/",
			func() views.Viewer {
				return &events.EventStoreAPIView{}
			},
		).Name(settings.ROUTE_EVENTS_EVENT_STORE),
		views.NewURL("/api/projects/project/{project_id:[0-9]+}/eventgroup",
			func() views.Viewer {
				return &events.EventGroupListAPIView{}
			},
		).Name(settings.ROUTE_EVENTS_EVENTGROUP_LIST).Middlewares(mids...),

		views.NewURL("/api/projects/project/{project_id:[0-9]+}/eventgroup/{eventgroup_id:[0-9]+}",
			func() views.Viewer {
				return &events.EventGroupDetailAPIView{}
			},
		).Name(settings.ROUTE_EVENTS_EVENTGROUP_DETAIL).Middlewares(mids...),

		views.NewURL("/api/projects/project/{project_id:[0-9]+}/eventgroup/{eventgroup_id:[0-9]+}/resolve",
			func() views.Viewer {
				return &events.EventGroupResolveAPIView{}
			},
		).Name(settings.ROUTE_EVENTS_EVENTGROUP_RESOLVE).Middlewares(mids...),

		views.NewURL(
			"/api/projects/project/{project_id:[0-9]+}/eventgroup/{eventgroup_id:[0-9]+}/event/",
			events.NewEventListView,
		).Name(settings.ROUTE_EVENTS_EVENT_LIST).Middlewares(mids...),
	}

	return result
}

func (e *EventsPlugin) Commands() []core.Commander {
	return []core.Commander{
		commands.NewEventWorkerCommand(e.context, e.SendOnEventSignal),
	}
}

func (e *EventsPlugin) Migrations() []core.Migrationer {
	return []core.Migrationer{
		core.NewMigration(
			models.MIGRATION_EVENTS_EVENTGROUP_INITIAL_ID,
			[]string{models.MIGRATION_EVENTS_EVENTGROUP_INITIAL},
			models.MIGRATION_EVENTS_EVENTGROUP_INITIAL_DEPENDENCIES,
		),
		core.NewMigration(
			models.MIGRATION_EVENTS_EVENT_INITIAL_ID,
			[]string{models.MIGRATION_EVENTS_EVENT_INITIAL},
			[]string{},
		),
	}
}

// signal handler to send event to frontend
func (e *EventsPlugin) OnEvent(event *models.Event, eventgroup *models.EventGroup) {
	glog.Infof("got signal %+v %+v", event, eventgroup)
}
