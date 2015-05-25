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

func TestTeamMemberDelete(t *testing.T) {

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

	Convey("Delete member - unauthenticated user", t, func() {
		session := apitest.NewSession()
		request := session.Request("DELETE", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Delete member - authenticated user", t, func() {
		session := apitest.NewSession().WithNewUser()
		request := session.Request("DELETE", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Delete member - member", t, func() {
		session := apitest.NewSession().WithUser(owner)
		request := session.Request("DELETE", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Delete invalid member - member admin", t, func() {
		session := apitest.NewSession().WithUser(owner)
		request := session.Request("DELETE", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", "0").Do()
		So(request.Response().Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("Delete valid member - member admin", t, func() {
		session := apitest.NewSession().WithUser(admin)
		request := session.Request("DELETE", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		request = session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result     []*serializers.TeamsTeamMemberDetailSerializer `json:"result"`
			ResultSize int                                            `json:"result_size"`
		}{}

		request.Scan(&response)

		found := false
		for _, item := range response.Result {
			if item.UserID.ToPrimaryKey() == owner.ID {
				found = true
			}
		}

		So(found, ShouldBeFalse)
	})

	Convey("Delete valid member - superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
		})
		request := session.Request("DELETE", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tma.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		request = session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result     []*serializers.TeamsTeamMemberDetailSerializer `json:"result"`
			ResultSize int                                            `json:"result_size"`
		}{}

		request.Scan(&response)

		found := false
		for _, item := range response.Result {
			if item.UserID.ToPrimaryKey() == admin.ID {
				found = true
			}
		}

		So(found, ShouldBeFalse)
	})

}
