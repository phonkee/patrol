package apitest

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/utils"
)

/*
Creates new project

First it creates team with user as owner, and then creates project with given team.
*/
func CreateProject(ctx *context.Context, user *models.User) (project *models.Project, err error) {

	var team *models.Team
	if team, err = CreateTeam(ctx, user); err != nil {
		return
	}

	project = models.NewProject(func(p *models.Project) {
		p.Name = "test project " + utils.RandomString(10)
		p.Platform = "any"
		p.TeamID = team.ID.ToForeignKey()
	})

	if err = project.Insert(ctx); err != nil {
		return
	}

	return
}
