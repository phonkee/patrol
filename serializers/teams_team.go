package serializers

import (
	"strings"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
)

/*
TeamsTeamCreateSerializer
	serializer for creating new team
*/
type TeamsTeamCreateSerializer struct {
	Name string `json:"name"`
}

/*
Clean
	cleans values *such as trim string spaces
	it's run before actual validation
*/
func (t *TeamsTeamCreateSerializer) Clean() {
	t.Name = strings.TrimSpace(t.Name)
}

/*
Validate
	validates serializer and returns validator result.
*/
func (s *TeamsTeamCreateSerializer) Validate(ctx *context.Context) *validator.Result {
	v := validator.New()
	v["name"] = models.ValidateTeamName()
	return v.Validate(s)
}

/*
Saves team to database and returns it
*/
func (s *TeamsTeamCreateSerializer) Save(ctx *context.Context, owner *models.User) (team *models.Team, err error) {
	team = models.NewTeam(func(team *models.Team) {
		team.Name = s.Name
		team.OwnerID = owner.PrimaryKey().ToForeignKey()
	})

	err = team.Insert(ctx)
	return
}
