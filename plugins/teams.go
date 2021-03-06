package plugins

import (
	"github.com/justinas/alice"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/views"
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
func (t *TeamsPlugin) URLs() []*views.URL {

	mids := []alice.Constructor{
		middlewares.AuthTokenValidMiddleware(),
	}
	return []*views.URL{

		views.NewURL("/api/teams/team/",
			func() views.Viewer {
				return &teams.TeamListAPIView{}
			},
		).Name(settings.ROUTE_TEAMS_TEAM_LIST).Middlewares(mids...),

		views.NewURL("/api/teams/team/{team_id:[0-9]+}",
			func() views.Viewer {
				return &teams.TeamDetailAPIView{}
			},
		).Name(settings.ROUTE_TEAMS_TEAM_DETAIL).Middlewares(mids...),

		views.NewURL(
			"/api/teams/team/{team_id:[0-9]+}/member/",
			teams.NewTeamMemberListAPIView,
		).Name(settings.ROUTE_TEAMS_TEAMMEMBER_LIST).Middlewares(mids...),

		views.NewURL(
			"/api/teams/team/{team_id:[0-9]+}/member/{teammember_id:[0-9]+}",
			teams.NewTeamMemberDetailAPIView,
		).Name(settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL).Middlewares(mids...),
	}
}

func (t *TeamsPlugin) Migrations() []core.Migrationer {
	return []core.Migrationer{
		core.NewMigration(models.MIGRATION_TEAMS_TEAM_INITIAL_ID, []string{models.MIGRATION_TEAMS_TEAM_INITIAL}, []string{}),
		core.NewMigration(models.MIGRATION_TEAMS_TEAM_MEMBER_INITIAL_ID, []string{models.MIGRATION_TEAMS_TEAM_MEMBER_INITIAL}, []string{}),
	}
}
