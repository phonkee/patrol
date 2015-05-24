package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/views"
	"github.com/phonkee/patrol/views/mixins"
)

/*
	Rest endpoint for project members list
*/
type ProjectMemeberDetailAPIView struct {
	core.JSONView

	context *context.Context

	// mixins that provide shortcuts methods
	mixins.AuthUserMixin
	mixins.ProjectsProjectMixin
	mixins.ProjectMemberTypeMixin

	// model instances
	project *models.Project
	user    *models.User
	memtype models.MemberType
}

func (p *ProjectMemeberDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.Context(r)

	// GetProject writes response so we only need to return error
	p.project = models.NewProject()
	if err = p.GetProject(p.project, w, r); err != nil {
		return
	}

	p.user = models.NewUser()
	if err = p.GetAuthUser(p.user, w, r); err != nil {
		return
	}

	tmm := models.NewTeamMemberManager(p.context)

	if p.memtype, err = tmm.MemberTypeByProject(p.project, p.user); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return views.ErrUnauthorized
	}

	return
}

/*
Retrieve list of user
*/
func (p *ProjectMemeberDetailAPIView) GET(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusOK).Write(w, r)
	return
}

/*
	Returns metadata
*/
func (p *ProjectMemeberDetailAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {

}

/*
Add member to team
*/
func (p *ProjectMemeberDetailAPIView) POST(w http.ResponseWriter, r *http.Request) {

}
