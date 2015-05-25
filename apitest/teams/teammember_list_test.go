package teams

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamMemberList(t *testing.T) {

	apitest.Setup()

	Convey("TeamMember list - unauthenticated user", t, func() {
		superuser, errsuperuser := apitest.CreateUser(patrol.Context, func(u *models.User) {
			u.IsSuperuser = true
		})
		So(errsuperuser, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, superuser)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession()
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("TeamMember list - authenticated user", t, func() {
		superuser, errsuperuser := apitest.CreateUser(patrol.Context, func(u *models.User) {
			u.IsSuperuser = true
		})
		So(errsuperuser, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, superuser)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("TeamMember list - member", t, func() {
		superuser, errsuperuser := apitest.CreateUser(patrol.Context, func(u *models.User) {
			u.IsSuperuser = true
		})
		So(errsuperuser, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, superuser)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()
		user := session.User()

		// add as team member
		tmm := models.NewTeamMemberManager(patrol.Context)
		_, errmt := tmm.SetTeamMemberType(team, user, models.MEMBER_TYPE_MEMBER)
		So(errmt, ShouldBeNil)

		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result []struct {
			} `json:"result"`
			ResultSize int `json:"result_size"`
		}{}

		request.Scan(&response)
		So(response.ResultSize, ShouldEqual, 1)
	})

	Convey("TeamMember list - superuser", t, func() {
		superuser, errsuperuser := apitest.CreateUser(patrol.Context, func(u *models.User) {
			u.IsSuperuser = true
		})
		So(errsuperuser, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, superuser)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithUser(superuser)

		request := session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result []struct {
			} `json:"result"`
			ResultSize int `json:"result_size"`
		}{}

		request.Scan(&response)
		So(response.ResultSize, ShouldEqual, 0)

	})

}
