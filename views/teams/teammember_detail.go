package teams

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/views"
	"github.com/phonkee/patrol/views/mixins"
)

/*
	Factory function to create view
*/
func NewTeamMemberDetailAPIView() core.Viewer {
	return &TeamMemberDetailAPIView{
		team:       models.NewTeam(),
		teammember: models.NewTeamMember(),
		user:       models.NewUser(),
	}
}

/*
	Team member detail rest endpoint
*/
type TeamMemberDetailAPIView struct {
	core.JSONView

	// used mixins
	mixins.AuthUserMixin
	mixins.TeamsTeamMixin
	mixins.TeamsTeamMemberMixin

	context *context.Context

	// stored instances
	membertype models.MemberType
	team       *models.Team
	teammember *models.TeamMember
	user       *models.User
}

/*
	Basic checks
*/
func (t *TeamMemberDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	if err = t.GetAuthUser(t.user, w, r); err != nil {
		return
	}

	if err = t.GetTeam(t.team, w, r); err != nil {
		return
	}

	t.context = t.GetContext(r)

	manager := models.NewTeamMemberManager(t.context)

	// check if user is member of team
	if t.membertype, err = manager.MemberType(t.team, t.user); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	if err = t.TeamsTeamMemberMixin.GetTeamMember(t.teammember, w, r); err != nil {
		return
	}

	// check ...
	if t.teammember.TeamID.ToPrimaryKey() != t.team.ID {
		response.New(http.StatusNotFound).Write(w, r)
		return views.ErrNotFound
	}

	return
}

/*
	Retrieve teammember detail
*/
func (t *TeamMemberDetailAPIView) GET(w http.ResponseWriter, r *http.Request) {
	result := &serializers.TeamsTeamMemberDetailSerializer{}
	result.TeamMember = *t.teammember
	result.LoadUser(t.context)

	response.New(http.StatusOK).Result(result).Write(w, r)
}

/*
Delete team member
*/
func (t *TeamMemberDetailAPIView) DELETE(w http.ResponseWriter, r *http.Request) {

	// check permissions - only admin (and superuser) can remove member of team
	if t.membertype != models.MEMBER_TYPE_ADMIN {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	user := models.NewUser()
	if err := t.teammember.User(user, t.context); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	manager := models.NewTeamMemberManager(t.context)
	if err := manager.RemoveTeamMember(t.team, user); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Write(w, r)
}

/*
Updates
	updates team member
*/
func (t *TeamMemberDetailAPIView) POST(w http.ResponseWriter, r *http.Request) {

	var err error

	// check permissions - only admin (and superuser) can remove member of team
	if t.membertype != models.MEMBER_TYPE_ADMIN {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	serializer := &serializers.TeamsTeamMemberUpdateSerializer{}
	if err = t.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	if vr := serializer.Validate(t.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	if err = serializer.Update(t.context, t.teammember); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Write(w, r)
}
