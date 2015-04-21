package plugins

import (
	"github.com/justinas/alice"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/views/teams"
)

func NewTeamsPlugin(context *context.Context) *TeamsPlugin {
	return &TeamsPlugin{context: context}
}

type TeamsPlugin struct {
	core.Plugin
	context *context.Context
}

func (t *TeamsPlugin) ID() string { return settings.TEAMS_PLUGIN_ID }
func (t *TeamsPlugin) URLViews() []*core.URLView {

	mids := []alice.Constructor{
		middlewares.AuthTokenValidMiddleware(),
	}
	return []*core.URLView{

		core.NewURLView("/api/teams/team/",
			func() core.Viewer {
				return &teams.TeamListAPIView{}
			}, settings.ROUTE_TEAMS_TEAM_LIST,
		).Middlewares(mids...),

		core.NewURLView("/api/teams/team/{team_id:[0-9]+}",
			func() core.Viewer {
				return &teams.TeamDetailAPIView{}
			},
			settings.ROUTE_TEAMS_TEAM_DETAIL,
		).Middlewares(mids...),
	}
}

func (t *TeamsPlugin) Migrations() []core.Migrationer {
	return []core.Migrationer{
		core.NewMigration(models.MIGRATION_TEAMS_TEAM_INITIAL_ID, []string{models.MIGRATION_TEAMS_TEAM_INITIAL}, []string{}),
		core.NewMigration(models.MIGRATION_TEAMS_TEAM_MEMBER_INITIAL_ID, []string{models.MIGRATION_TEAMS_TEAM_MEMBER_INITIAL}, []string{}),
	}
}
