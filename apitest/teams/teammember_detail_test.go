package teams

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamMemberDetail(t *testing.T) {

	apitest.Setup()

	owner, errowner := apitest.CreateUser(patrol.Context)
	if errowner != nil {
		t.FailNow()
	}
	admin, errowner := apitest.CreateUser(patrol.Context)
	if errowner != nil {
		t.FailNow()
	}

	team, errteam := apitest.CreateTeam(patrol.Context, owner)
	if errteam != nil {
		t.FailNow()
	}

	tmm := models.NewTeamMemberManager(patrol.Context)

	tm, errtmm := tmm.SetTeamMemberType(team, owner, models.MEMBER_TYPE_MEMBER)
	if errtmm != nil {
		t.FailNow()
	}

	tma, errtmma := tmm.SetTeamMemberType(team, admin, models.MEMBER_TYPE_ADMIN)
	if errtmma != nil {
		t.FailNow()
	}
	_ = tma

	Convey("TeamMember Detail - unauthenticated", t, func() {
		session := apitest.NewSession()
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("TeamMember Detail - authenticated", t, func() {
		session := apitest.NewSession().WithNewUser()
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("TeamMember Detail - member", t, func() {
		session := apitest.NewSession().WithUser(owner)
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)
		response := struct {
			Result serializers.TeamsTeamMemberDetailSerializer `json:"result"`
		}{}
		request.Scan(&response)
		So(response.Result.User.ID, ShouldEqual, owner.ID)
	})

	Convey("TeamMember Detail - member admin", t, func() {
		session := apitest.NewSession().WithUser(admin)
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)
		response := struct {
			Result serializers.TeamsTeamMemberDetailSerializer `json:"result"`
		}{}
		request.Scan(&response)
		So(response.Result.User.ID, ShouldEqual, owner.ID)

	})

	Convey("TeamMember Detail - superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
		})
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)
		response := struct {
			Result serializers.TeamsTeamMemberDetailSerializer `json:"result"`
		}{}
		request.Scan(&response)
		So(response.Result.User.ID, ShouldEqual, owner.ID)
	})

}
