package teams

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamMemberCreate(t *testing.T) {

	apitest.Setup()

	Convey("Create member - unauthorized user", t, func() {
		sudo, errsudo := apitest.CreateUser(patrol.Context, func(user *models.User) {
			user.IsSuperuser = true
		})
		So(errsudo, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, sudo)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession()
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String())
		request.StringBody("{}").Do()

		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)

	})

	Convey("Create member - authorized user", t, func() {
		sudo, errsudo := apitest.CreateUser(patrol.Context, func(user *models.User) {
			user.IsSuperuser = true
		})
		So(errsudo, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, sudo)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()
		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String())
		request.StringBody("{}").Do()

		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Create member - member", t, func() {
		sudo, errsudo := apitest.CreateUser(patrol.Context, func(user *models.User) {
			user.IsSuperuser = true
		})
		So(errsudo, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, sudo)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()

		// add as team member
		tmm := models.NewTeamMemberManager(patrol.Context)
		_, errmt := tmm.SetTeamMemberType(team, session.User(), models.MEMBER_TYPE_MEMBER)
		So(errmt, ShouldBeNil)

		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String())
		request.StringBody("{}").Do()

		So(request.Response().Code, ShouldEqual, http.StatusForbidden)

	})

	Convey("Create member - member admin", t, func() {
		sudo, errsudo := apitest.CreateUser(patrol.Context, func(user *models.User) {
			user.IsSuperuser = true
		})
		So(errsudo, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, sudo)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()

		// add as team member
		tmm := models.NewTeamMemberManager(patrol.Context)
		_, errmt := tmm.SetTeamMemberType(team, session.User(), models.MEMBER_TYPE_ADMIN)
		So(errmt, ShouldBeNil)

		other, errother := apitest.CreateUser(patrol.Context)
		So(errother, ShouldBeNil)

		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String())
		request.JSONBody(map[string]interface{}{
			"user_id": other.ID,
			"type":    models.MEMBER_TYPE_MEMBER,
		}).Do()

		So(request.Response().Code, ShouldEqual, http.StatusCreated)

		// find other user in array
		request = session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result []struct {
				UserID types.PrimaryKey `json:"user_id"`
			} `json:"result"`
			ResultSize int `json:"result_size"`
		}{}

		request.Scan(&response)

		found := false
		for _, item := range response.Result {
			if item.UserID == other.ID {
				found = true
			}
		}

		So(found, ShouldBeTrue)
	})

	Convey("Create member - superuser", t, func() {
		sudo, errsudo := apitest.CreateUser(patrol.Context, func(user *models.User) {
			user.IsSuperuser = true
		})
		So(errsudo, ShouldBeNil)

		team, errteam := apitest.CreateTeam(patrol.Context, sudo)
		So(errteam, ShouldBeNil)

		session := apitest.NewSession().WithUser(sudo)

		other, errother := apitest.CreateUser(patrol.Context)
		So(errother, ShouldBeNil)

		request := session.Request("POST", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String())
		request.JSONBody(map[string]interface{}{
			"user_id": other.ID,
			"type":    models.MEMBER_TYPE_MEMBER,
		}).Do()

		So(request.Response().Code, ShouldEqual, http.StatusCreated)

		// find other user in array
		request = session.Request("GET", settings.ROUTE_TEAMS_TEAMMEMBER_LIST, "team_id", team.ID.String()).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		response := struct {
			Result []struct {
				UserID types.PrimaryKey `json:"user_id"`
			} `json:"result"`
			ResultSize int `json:"result_size"`
		}{}

		request.Scan(&response)

		found := false
		for _, item := range response.Result {
			if item.UserID == other.ID {
				found = true
			}
		}

		So(found, ShouldBeTrue)

	})

}
