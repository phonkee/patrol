package teams

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamCreate(t *testing.T) {

	apitest.Setup()

	Convey("Create team for unauthorized user", t, func() {
		session := apitest.NewSession()

		serializer := serializers.TeamsTeamCreateSerializer{
			Name: "test team" + utils.RandomString(10),
		}

		request := session.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).JSONBody(serializer)
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Create team for non superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
			user.IsActive = true
		})

		serializer := serializers.TeamsTeamCreateSerializer{
			Name: "test team" + utils.RandomString(10),
		}

		request := session.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).JSONBody(serializer)
		So(request.Do().Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Create team for superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
			user.IsActive = true
		})

		serializer := serializers.TeamsTeamCreateSerializer{
			Name: "test team" + utils.RandomString(10),
		}

		request := session.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).JSONBody(serializer)
		So(request.Do().Response().Code, ShouldEqual, http.StatusCreated)
	})

}
