package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/views/auth"
	"github.com/phonkee/patrol/views/mixins"
)

/*
	Rest endpoint for project members list
*/
type ProjectMemeberListAPIView struct {
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

func (p *ProjectMemeberListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
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

	if p.memtype, err = p.GetMemberType(p.project, p.user, w, r); err != nil {
		return
	}

	return
}

/*
Retrieve list of user
*/
func (p *ProjectMemeberListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	tmm := models.NewTeamMemberManager(p.context)
	memberlist := models.NewTeamMemberList()
	if err := tmm.Filter(&memberlist, tmm.QueryFilterProject(p.project)); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	/*
		Iterate over team members, load user info and add
	*/

	var result = make([]*ProjectMemberListItem, len(memberlist))

	for i, item := range memberlist {
		member := &ProjectMemberListItem{
			ID:     item.ID,
			Type:   item.Type,
			UserID: item.UserID,
			User:   &auth.UserDetailSerializer{},
		}
		result[i] = member

		user := models.NewUser()
		if err := user.Manager(p.context).GetByID(user, item.UserID); err != nil {
			continue
		}
		result[i].User.FromUser(user)
	}

	response.New(http.StatusOK).Result(result).ResultSize(len(result)).Write(w, r)
	return
}

/*
	Returns metadata
*/
func (p *ProjectMemeberListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	md := metadata.New("List of project members")
	md.ActionRetrieve().From(&ProjectMemberListItem{})

	response.New(http.StatusOK).Metadata(md).Write(w, r)
	return
}

/*
Add member to team
*/
func (p *ProjectMemeberListAPIView) POST(w http.ResponseWriter, r *http.Request) {

}
