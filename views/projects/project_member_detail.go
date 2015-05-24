package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
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

	// check if user is member of project
	if p.memtype, err = p.GetMemberType(p.project, p.user, w, r); err != nil {
		return
	}

	return
}

/*
	Delete member from team
*/
func (p *ProjectMemeberDetailAPIView) DELETE(w http.ResponseWriter, r *http.Request) {
	// only admin can delete member
	if p.memtype != models.MEMBER_TYPE_ADMIN {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}
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
