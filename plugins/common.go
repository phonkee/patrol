package plugins

import (
	"github.com/phonkee/patrol/commands"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/views/common"
)

func NewCommonPlugin(context *context.Context, pr *core.PluginRegistry) core.Pluginer {
	return &CommonPlugin{context: context, pr: pr}
}

/*Core plugin encapsulates core patrol functionality
 */
type CommonPlugin struct {
	core.Plugin
	context *context.Context
	pr      *core.PluginRegistry
}

func (p *CommonPlugin) ID() string { return settings.COMMON_PLUGIN_ID }

// Builtin commands
func (p *CommonPlugin) Commands() []core.Commander {
	return []core.Commander{
		commands.NewCommonListCommandsCommand(p.context, p.pr),
		commands.NewCommonMigrateCommand(p.context, p.pr),
		commands.NewCommonHttpServerCommand(p.context, p.pr),
		commands.NewCommonListRoutesCommand(p.context, p.pr),
	}
}

// list of urls
func (c *CommonPlugin) URLs() []*views.URL {
	return []*views.URL{
		views.NewURL("/api/version", func() views.Viewer { return &common.VersionAPIView{} }).Name(settings.ROUTE_COMMON_VERSION),
		views.NewURL("/api/monitor", func() views.Viewer { return &common.MonitorAPIView{} }).Name(settings.ROUTE_COMMON_MONITOR),
	}
}
