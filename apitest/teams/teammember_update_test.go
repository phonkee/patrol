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

func TestTeamMemberUpdate(t *testing.T) {

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

	Convey("Update Team Member - unauthenticated user", t, func() {
		session := apitest.NewSession()
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String())
		request.StringBody("{}").Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Update Team Member - authenticated user", t, func() {
		session := apitest.NewSession().WithNewUser()
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String())
		request.StringBody("{}").Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Update Team Member - member", t, func() {
		session := apitest.NewSession().WithUser(owner)
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String())
		request.StringBody("{}").Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Update Team Member - admin member", t, func() {
		session := apitest.NewSession().WithUser(admin)
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tm.ID.String())
		request.JSONBody(map[string]interface{}{
			"type": models.MEMBER_TYPE_ADMIN,
		}).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		request = session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result     []*serializers.TeamsTeamMemberDetailSerializer `json:"result"`
			ResultSize int                                            `json:"result_size"`
		}{}

		request.Scan(&response)

		for _, item := range response.Result {
			if item.UserID.ToPrimaryKey() == owner.ID {
				So(item.Type, ShouldEqual, models.MEMBER_TYPE_ADMIN)
			}
		}
	})

	Convey("Update Team Member - superuser", t, func() {
		// change admin to member
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
		})
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_DETAIL, "team_id", team.ID.String(), "teammember_id", tma.ID.String())
		request.JSONBody(map[string]interface{}{
			"type": models.MEMBER_TYPE_MEMBER,
		}).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		request = session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result     []*serializers.TeamsTeamMemberDetailSerializer `json:"result"`
			ResultSize int                                            `json:"result_size"`
		}{}

		request.Scan(&response)

		for _, item := range response.Result {
			if item.UserID.ToPrimaryKey() == admin.ID {
				So(item.Type, ShouldEqual, models.MEMBER_TYPE_MEMBER)
			}
		}
	})
}
