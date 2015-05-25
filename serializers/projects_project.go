package serializers

import (
	"strings"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
)

/*
ProjectsProjectCreateSerializer
	serializer for creating new project

*/
type ProjectsProjectCreateSerializer struct {
	Name     string           `json:"name"    validator:"name"`
	Platform string           `json:"platform"`
	TeamID   types.ForeignKey `json:"team_id" validator:"team_id"`
}

/*
	Cleans data in serializer
*/
func (p *ProjectsProjectCreateSerializer) Clean() {
	p.Name = strings.TrimSpace(p.Name)
	p.Platform = strings.TrimSpace(p.Platform)
}

/*
Validate
	validates new project
*/
func (p *ProjectsProjectCreateSerializer) Validate(context *context.Context) *validator.Result {
	validator := validator.New()
	validator["name"] = models.ValidateProjectName()
	validator["team_id"] = models.ValidateTeamID(context)
	return validator.Validate(p)
}

/*
	Saves new project to database
*/
func (p *ProjectsProjectCreateSerializer) Save(context *context.Context, team *models.Team, author *models.User) (project *models.Project, err error) {
	project = models.NewProject(func(proj *models.Project) {
		proj.Name = p.Name
		proj.Platform = p.Platform
		proj.TeamID = team.ID.ToForeignKey()
	})

	if err = project.Insert(context); err != nil {
		return
	}

	// create new project key
	projectkey := models.NewProjectKey(func(projectKey *models.ProjectKey) {
		projectKey.UserID = author.ID.ToForeignKey()
		projectKey.UserAddedID = author.ID.ToForeignKey()
		projectKey.ProjectID = project.ID.ToForeignKey()
	})

	if err = projectkey.Insert(context); err != nil {
		return
	}

	return
}
