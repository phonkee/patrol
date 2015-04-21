package plugins

import (
	"github.com/golang/glog"
	"github.com/phonkee/patrol/commands"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/signals"
	"github.com/phonkee/patrol/views/auth"
)

func NewAuthPlugin(context *context.Context, pr *core.PluginRegistry) core.Pluginer {
	return &AuthPlugin{context: context, pr: pr}
}

type AuthPlugin struct {
	core.Plugin
	context                         *context.Context
	pr                              *core.PluginRegistry
	OnSuccessfulLoginSignalHandlers []signals.OnSuccessfulLoginSignalHandler
}

// Plugin identifier
func (p *AuthPlugin) ID() string { return settings.AUTH_PLUGIN_ID }
func (a *AuthPlugin) Init() error {
	a.OnSuccessfulLoginSignalHandlers = []signals.OnSuccessfulLoginSignalHandler{}
	a.pr.Do(func(plugin core.Pluginer) error {
		if t, ok := plugin.(signals.OnSuccessfulLoginSignalHandler); ok {
			glog.V(2).Infof("event signals: adding %T as OnLoginSignalHandler.", plugin)
			a.OnSuccessfulLoginSignalHandlers = append(a.OnSuccessfulLoginSignalHandlers, t)
		}
		return nil
	})
	return nil
}

// Return all auth views
func (a *AuthPlugin) URLViews() []*core.URLView {
	urls := []*core.URLView{
		core.NewURLView(
			"/api/auth/login", func() core.Viewer {
				return &auth.AuthLoginAPIView{
					LoginSignal: a.SendSuccessfulLoginSignal,
				}
			},
			settings.ROUTE_AUTH_LOGIN,
		),

		core.NewURLView(
			"/api/auth/me", func() core.Viewer {
				return &auth.AuthMeAPIView{}
			},
			settings.ROUTE_AUTH_ME,
		).Middlewares(middlewares.AuthTokenValidMiddleware()),

		core.NewURLView(
			"/api/auth/user/", func() core.Viewer {
				return auth.NewUserListAPIView()
			},
			settings.ROUTE_AUTH_USER_LIST,
		).Middlewares(middlewares.AuthTokenValidMiddleware()),
	}
	return urls
}

// Returns migrations for auth plugin
func (a *AuthPlugin) Migrations() []core.Migrationer {
	return []core.Migrationer{
		core.NewMigration(models.MIGRATION_AUTH_USER_INITIAL_ID, []string{models.MIGRATION_AUTH_USER_INITIAL}, []string{}),
		core.NewMigration(models.MIGRATION_AUTH_PERMISSION_INITIAL_ID, []string{models.MIGRATION_AUTH_PERMISSION_INITIAL}, []string{}),
	}
}

// send succesfull login signal
func (a *AuthPlugin) SendSuccessfulLoginSignal(user *models.User) error {
	glog.V(2).Infof("signal: sending successful login signal")
	for _, sh := range a.OnSuccessfulLoginSignalHandlers {
		func() {
			sh.OnSuccessfulLogin(user)
			defer func() {
				if err := recover(); err != nil {
					glog.Errorf("OnLoginSignalHandler %v panicked", sh)
				}
			}()
		}()
	}
	return nil
}

func (a *AuthPlugin) Commands() []core.Commander {
	return []core.Commander{
		commands.NewAuthCreateSuperuserCommand(a.context),
	}
}
