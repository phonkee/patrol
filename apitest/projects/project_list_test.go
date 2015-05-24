package projects

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProjectList(t *testing.T) {
	apitest.Setup()

	Convey("Test list of projects for non logged user", t, func() {
		session := apitest.NewSession()
		request := session.Request("GET", settings.ROUTE_PROJECTS_PROJECT_LIST)
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Test list of projects", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
		})
		request := session.Request("GET", settings.ROUTE_PROJECTS_PROJECT_LIST)
		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})
}
