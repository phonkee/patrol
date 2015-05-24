package plugins

import (
	"github.com/justinas/alice"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/middlewares"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/views/projects"
)

func NewProjectsPlugin(context *context.Context) core.Pluginer {
	return &ProjectsPlugin{context: context}
}

/*
Projects plugin -
handles projects
*/
type ProjectsPlugin struct {
	core.Plugin
	context *context.Context
}

func (p *ProjectsPlugin) ID() string { return settings.PROJECTS_PLUGIN_ID }
func (p *ProjectsPlugin) URLViews() []*core.URLView {
	mids := []alice.Constructor{
		middlewares.AuthTokenValidMiddleware(),
	}
	return []*core.URLView{
		core.NewURLView(
			"/api/projects/project/", func() core.Viewer { return &projects.ProjectListAPIView{} },
		).Name(settings.ROUTE_PROJECTS_PROJECT_LIST).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}",
			func() core.Viewer {
				return &projects.ProjectDetailAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECT_DETAIL).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}/eventgroup",
			func() core.Viewer {
				return &projects.ProjectDetailEventGroupListAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECT_DETAIL_EVENTGROUP_LIST).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}/eventgroup/{eventgroup_id:[0-9]+}",
			func() core.Viewer {
				return &projects.ProjectDetailEventGroupDetailAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECT_DETAIL_EVENTGROUP_DETAIL).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}/eventgroup/{eventgroup_id:[0-9]+}/resolve",
			func() core.Viewer {
				return &projects.ProjectDetailEventGroupResolveAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECT_DETAIL_EVENTGROUP_RESOLVE).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}/key",
			func() core.Viewer {
				return &projects.ProjectKeyListAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECTKEY_LIST).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}/key/{projectkey_id:[0-9]+}",
			func() core.Viewer {
				return &projects.ProjectKeyDetailAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECTKEY_DETAIL).Middlewares(mids...),

		core.NewURLView("/api/projects/project/{project_id:[0-9]+}/eventgroup/{eventgroup_id:[0-9]+}/event/",
			func() core.Viewer {
				return &projects.ProjectDetailEventGroupListAPIView{}
			},
		).Name(settings.ROUTE_PROJECTS_PROJECT_DETAIL_EVENTGROUP_EVENT_LIST).Middlewares(mids...),
	}
}
func (p *ProjectsPlugin) Migrations() []core.Migrationer {
	return []core.Migrationer{
		core.NewMigration(
			models.MIGRATION_PROJECT_INITIAL_ID,           // migration id
			[]string{models.MIGRATION_PROJECT_INITIAL},    // migration queries
			models.MIGRATION_PROJECT_INITIAL_DEPENDENCIES, // migration dependencies
			p.PostInitialMigration,                        // post migration method
		),
		core.NewMigration(
			models.MIGRATION_PROJECT_KEY_INITIAL_ID,
			[]string{models.MIGRATION_PROJECT_KEY_INITIAL},
			[]string{},
		),
	}
}

func (p *ProjectsPlugin) PostInitialMigration(ctx *context.Context) (err error) {
	manager := models.NewPermissionManager(ctx)
	permission := manager.NewPermission(func(perm *models.Permission) {
		perm.Codename = settings.PERMISSION_TEAMS_TEAM_ADD
		perm.Name = "Can add new project"
	})
	if err = permission.Insert(ctx); err != nil {
		return
	}
	return
}
