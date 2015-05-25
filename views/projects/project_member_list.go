package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/serializers"
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
	mixins.TeamsTeamMemberMixin

	// model instances
	project *models.Project
	user    *models.User
	memtype models.MemberType
}

func (p *ProjectMemeberListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)

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

	result := []*serializers.TeamsTeamMemberDetailSerializer{}

	if err := p.TeamsTeamMemberMixin.Filter(&result, p.context, tmm.QueryFilterProject(p.project)); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(result).ResultSize(len(result)).Write(w, r)
	return
}

/*
	Returns metadata
*/
func (p *ProjectMemeberListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	md := metadata.New("List of project members")
	md.ActionRetrieve().From(&serializers.TeamsTeamMemberDetailSerializer{})

	// admin has permissions to add new member
	if p.memtype == models.MEMBER_TYPE_ADMIN {
		create := md.ActionCreate().From(&serializers.TeamsTeamMemberCreateSerializer{})
		create.Field("type").Choices.Add(models.MEMBER_TYPE_MEMBER, "member").Add(models.MEMBER_TYPE_ADMIN, "admin")
	}

	response.New(http.StatusOK).Metadata(md).Write(w, r)
	return
}

/*
Add member to team
*/
func (p *ProjectMemeberListAPIView) POST(w http.ResponseWriter, r *http.Request) {

	var err error

	// only admin (and superuser) can add new members
	if p.memtype != models.MEMBER_TYPE_ADMIN {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	serializer := &serializers.TeamsTeamMemberCreateSerializer{}
	if err = p.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	team := models.NewTeam()
	if err = p.project.Team(team, p.context); err != nil {
		response.New(http.StatusNotFound).Error(err).Write(w, r)
		return
	}

	if vr := serializer.Validate(p.context, team); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	var result *models.TeamMember

	if result, err = serializer.Save(p.context, team); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusCreated).Result(result).Write(w, r)
}
