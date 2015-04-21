package plugins

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"github.com/phonkee/patrol/commands"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/settings"
)

var (
	staticDir = flag.String("static_dir", "", "Static dir to serve static")
)

func NewStaticPlugin(context *context.Context, pr *core.PluginRegistry) core.Pluginer {
	return &StaticPlugin{context: context, pr: pr}
}

/*Core plugin encapsulates core patrol functionality
 */
type StaticPlugin struct {
	core.Plugin
	context *context.Context
	pr      *core.PluginRegistry
}

func (p *StaticPlugin) ID() string { return settings.STATIC_PLUGIN_ID }

func (p *StaticPlugin) OnHttpServerStart() {
	glog.Infof("server started, listening in static plugin")

	// TODO: serving static from within binary (go-bindata)
	// http.Handle("/", static.HttpBindata())

	// serving static dir from custom location
	if *staticDir != "" {
		glog.Infof("Serving static from custom location \"%s\"", *staticDir)
		fs := http.FileServer(
			http.Dir(*staticDir),
		)

		// serve static
		p.context.Router.PathPrefix("/").Handler(
			alice.New(
				middlewares.ContextMiddleware(p.context),
				middlewares.RequestLogMiddleware(),
				middlewares.RecoveryMiddleware(),
			).Then(fs),
		)
	}
}

// Builtin commands
func (p *StaticPlugin) Commands() []core.Commander {
	return []core.Commander{
		commands.NewStaticDumpCommand(p.context),
	}
}
