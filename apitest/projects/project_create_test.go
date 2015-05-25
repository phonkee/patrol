package projects

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProjectCreate(t *testing.T) {
	apitest.Setup()

	user, _ := apitest.CreateUser(patrol.Context, func(user *models.User) {
		user.IsSuperuser = true
		user.IsActive = true
	})
	team := models.NewTeam(func(t *models.Team) {
		t.Name = "test team " + utils.RandomString(10)
		t.OwnerID = user.ID.ToForeignKey()
	})
	team.Insert(patrol.Context)

	Convey("Test create project for nonauthenticated user", t, func() {
		session := apitest.NewSession()
		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECT_LIST).Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Test create project for normal user", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsActive = true
		})

		serializer := serializers.ProjectsProjectCreateSerializer{
			Name:     "test project " + utils.RandomString(10),
			Platform: "any",
			TeamID:   team.ID.ToForeignKey(),
		}

		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECT_LIST).JSONBody(serializer).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Test create project for authenticated superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
			user.IsActive = true
		})

		serializer := serializers.ProjectsProjectCreateSerializer{
			Name:     "test project " + utils.RandomString(10),
			Platform: "any",
			TeamID:   team.ID.ToForeignKey(),
		}

		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECT_LIST).JSONBody(serializer)
		So(request.Do().Response().Code, ShouldEqual, http.StatusCreated)
	})
}
