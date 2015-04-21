package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProjects(t *testing.T) {
	if err := patrol.Setup(); err != nil {
		if err != patrol.ErrPatrolAlreadySetup {
			fmt.Printf("patrol: setup failed with error: %s", err)
		}
	}

	patrol.Run([]string{"migrate"})

	Convey("Test retrieve list of projects / insert project", t, func() {

		sudo := NewRequestsSession(patrol.Context).WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
			user.IsActive = true
		})

		So(sudo.Request("GET", settings.ROUTE_PROJECTS_PROJECT_LIST).Do().Response().Code, ShouldEqual, http.StatusOK)

		ser := serializers.TeamCreateSerializer{
			Name: "team test",
		}

		var result struct {
			Result models.Team `json:"result"`
		}

		// create new team
		req := sudo.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).JSONBody(ser).Do().Scan(&result)
		So(req.Error(), ShouldBeNil)
		So(req.Response().Code, ShouldEqual, http.StatusCreated)

		serializer := serializers.ProjectCreateSerializer{
			Name:     utils.RandomString(20),
			Platform: "go",
			TeamID:   result.Result.PrimaryKey().ToForeignKey(),
		}

		reqCreate := sudo.Request("POST", settings.ROUTE_PROJECTS_PROJECT_LIST).JSONBody(serializer).Do()
		So(reqCreate.Response().Code, ShouldEqual, http.StatusCreated)

		// unauthorized user cannot create project
		unauthorized := NewRequestsSession(patrol.Context)
		So(unauthorized.Request("POST", settings.ROUTE_PROJECTS_PROJECT_LIST).JSONBody(serializer).Do().Response().Code, ShouldEqual, http.StatusUnauthorized)

		// normal user without permissions cannot also
		normal := NewRequestsSession(patrol.Context).WithNewUser()
		So(normal.Request("POST", settings.ROUTE_PROJECTS_PROJECT_LIST).JSONBody(serializer).Do().Response().Code, ShouldEqual, http.StatusForbidden)

	})
}
