package patrol

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/phonkee/patrol/context"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"github.com/mgutz/ansi"
	"github.com/phonkee/patrol/core"
	ms "github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/plugins"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/views/common"
)

var (

	// Patrol context
	Context *context.Context

	// chain of middlewares
	middlewares alice.Chain

	// plugin registry
	pluginRegistry *core.PluginRegistry

	// errors
	ErrNoArgs             = errors.New("no run arguments given.")
	ErrPatrolAlreadySetup = errors.New("patrol already setup.")
	issetup               bool

	// alertcolor
	alertcolor = ansi.ColorFunc("red+h:black")
)

/*
	Patrol bootstrap
*/
func init() {
	var err error

	// @TODO: solve this one
	if os.Getenv("TESTING") != "TRUE" {
		// parse flags
		flag.Parse()
	}

	// Create patrol context
	if Context, err = context.New(
		settings.SETTINGS_DATABASE_DSN,
		settings.SETTINGS_MESSAGE_QUEUE_DSN,
		settings.SETTINGS_CACHE_DSN,
	); err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	// setup core middlewares
	middlewares = alice.New(
		ms.ContextMiddleware(Context),
		ms.RequestLogMiddleware(),
		ms.RecoveryMiddleware(),
	)

	// initialize plugin registry
	pluginRegistry = core.NewPluginRegistry()

	// initialize router

	Context.Set(context.SECRET_KEY, settings.SETTINGS_SECRET_KEY)

	// register not found handler
	Context.Router.NotFoundHandler = middlewares.ThenFunc(common.NotFoundHandler())

	runtime.GOMAXPROCS(settings.SETTINGS_GOMAXPROCS)
	glog.V(2).Infof("patrol: hello user, welcome to patrol with GOMAXPROCS=%d.", settings.SETTINGS_GOMAXPROCS)

	if err = bootstrap(); err != nil {
		glog.Errorf("patrol: bootstrap error %s.", err)
		os.Exit(1)
	}
}

/*
	Bootstraps patrol application
*/
func bootstrap() error {
	glog.Info(settings.PATROL_LOGO)
	if settings.DEBUG {
		glog.Infof("patrol: running in %s mode.", alertcolor("DEBUG"))
	}

	// list of patrol builtin plugins to register
	plugins := []core.Pluginer{
		plugins.NewCommonPlugin(Context, pluginRegistry),
		plugins.NewAuthPlugin(Context, pluginRegistry),
		plugins.NewEventsPlugin(Context, pluginRegistry),
		plugins.NewProjectsPlugin(Context),
		plugins.NewTeamsPlugin(Context),
		plugins.NewStaticPlugin(Context, pluginRegistry),
		plugins.NewRealtimePlugin(Context, pluginRegistry),
	}
	for _, p := range plugins {
		if err := pluginRegistry.RegisterPlugin(p); err != nil {
			return err
		}
		// add core plugin ids to RESTRICTED_PLUGIN_IDS
		settings.RESTRICTED_PLUGIN_IDS = append(settings.RESTRICTED_PLUGIN_IDS, p.ID())
	}

	return nil
}

/*
	Calls Init on all plugins
*/
func initPlugins() error {
	return pluginRegistry.InitPlugins()
}

// Iterate over all plugins
func PluginsDo(f func(plugin core.Pluginer) error) error {
	return pluginRegistry.Do(f)
}

// returns plugin by name
func Plugin(id string) (core.Pluginer, error) {
	return pluginRegistry.Plugin(id)
}

/*
	Setups handlers for router from all URLs
*/
func setupHandlers() error {
	return pluginRegistry.Do(func(plugin core.Pluginer) error {
		for _, url := range plugin.URLs() {
			if err := url.Register(Context.Router, middlewares); err != nil {
				return err
			}

		}
		return nil
	})
}

/*
	Sets up patrol application.
	Registers all URLs
*/
func Setup() error {
	if issetup {
		return ErrPatrolAlreadySetup
	}
	issetup = true
	glog.V(2).Infoln("patrol: init plugins.")
	if err := initPlugins(); err != nil {
		return err
	}

	glog.V(2).Infoln("patrol: setup handlers.")
	if err := setupHandlers(); err != nil {
		return err
	}
	return nil
}

/*
	Registers plugin to patrol
*/
func RegisterPlugin(plugin core.Pluginer) error {
	for _, rpn := range settings.RESTRICTED_PLUGIN_IDS {
		if rpn == plugin.ID() {
			return fmt.Errorf("plugin %s cannot be registered because it's in RESTRICTED_PLUGIN_IDS (%v).", plugin.ID(), settings.RESTRICTED_PLUGIN_IDS)
		}
	}

	return pluginRegistry.RegisterPlugin(plugin)
}

// runs command
func Run(args []string) (err error) {
	if len(args) == 0 {
		err = ErrNoArgs
		return
	}

	var command core.Commander
	if command, err = pluginRegistry.Command(args[0]); err != nil {
		return err
	}

	// parse command line positional args for command
	if err = command.ParseArgs(args[1:]); err != nil {
		return
	}

	// runs command
	if err = command.Run(); err != nil {
		return
	}

	return nil
}
