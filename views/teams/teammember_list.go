package teams

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/utils"
	"github.com/phonkee/patrol/views/mixins"
)

// Factory for new TeamMemberListAPIView
func NewTeamMemberListAPIView() core.Viewer {
	return &TeamMemberListAPIView{
		team: models.NewTeam(),
		user: models.NewUser(),
	}
}

/*
	Team member list endpoint
*/
type TeamMemberListAPIView struct {
	core.JSONView

	// userd mixins
	mixins.AuthUserMixin
	mixins.TeamsTeamMixin
	mixins.TeamsTeamMemberMixin

	context *context.Context

	membertype models.MemberType
	team       *models.Team
	user       *models.User
}

/*
	Basic checks
*/
func (t *TeamMemberListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	t.context = t.GetContext(r)

	if err = t.GetTeam(t.team, w, r); err != nil {
		return
	}

	if err = t.GetAuthUser(t.user, w, r); err != nil {
		return
	}

	manager := models.NewTeamMemberManager(t.context)

	// check if user is member of team
	if t.membertype, err = manager.MemberType(t.team, t.user); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	return
}

/*
	Retrieve list of all team members
*/
func (t *TeamMemberListAPIView) GET(w http.ResponseWriter, r *http.Request) {

	result := []*serializers.TeamsTeamMemberDetailSerializer{}

	if err := t.TeamsTeamMemberMixin.Filter(&result, t.context, utils.QueryFilterWhere("team_id = ?", t.team.ID)); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(result).ResultSize(len(result)).Write(w, r)
}

/*
Metadata request
*/
func (t *TeamMemberListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	md := metadata.New("team member list endpoint")
	md.ActionRetrieve().From(&serializers.TeamsTeamMemberDetailSerializer{})

	// member type admin and superuser can create new team members
	if t.membertype == models.MEMBER_TYPE_ADMIN {
		create := md.ActionCreate().From(&serializers.TeamsTeamMemberCreateSerializer{})
		create.Field("type").Choices.Add(models.MEMBER_TYPE_MEMBER, "member").Add(models.MEMBER_TYPE_ADMIN, "admin")
	}

	response.New(http.StatusOK).Metadata(md).Write(w, r)
}

/*
Add member to team
*/
func (t *TeamMemberListAPIView) POST(w http.ResponseWriter, r *http.Request) {

	var err error

	// only admin (and superuser) can add new members
	if t.membertype != models.MEMBER_TYPE_ADMIN {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	serializer := &serializers.TeamsTeamMemberCreateSerializer{}
	if err = t.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	if vr := serializer.Validate(t.context, t.team); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	var result *models.TeamMember

	if result, err = serializer.Save(t.context, t.team); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusCreated).Result(result).Write(w, r)
}
