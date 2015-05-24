package apitest

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/utils"
)

func CreateProject(ctx *context.Context, user *models.User) (project *models.Project, err error) {

	team := models.NewTeam(func(t *models.Team) {
		t.Name = "test team " + utils.RandomString(10)
		t.OwnerID = user.ID.ToForeignKey()
	})

	if err = team.Insert(ctx); err != nil {
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
