package apitest

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/utils"
)

func CreateTeam(ctx *context.Context, owner *models.User) (team *models.Team, err error) {

	team = models.NewTeam(func(t *models.Team) {
		t.Name = "test team " + utils.RandomString(10)
		t.OwnerID = owner.ID.ToForeignKey()
	})

	if err = team.Insert(ctx); err != nil {
		return
	}

	return
}
