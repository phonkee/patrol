package events

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEventList(t *testing.T) {

	apitest.Setup()

	user, erruser := apitest.CreateUser(patrol.Context)
	if erruser != nil {
		t.FailNow()
	}
	project, errproject := apitest.CreateProject(patrol.Context, user)
	if errproject != nil {
		t.FailNow()
	}

	team := models.NewTeam()
	errteam := project.Team(team, patrol.Context)
	if errteam != nil {
		t.FailNow()
	}

	// add as team member
	tmm := models.NewTeamMemberManager(patrol.Context)
	_, errmt := tmm.SetTeamMemberType(team, user, models.MEMBER_TYPE_ADMIN)
	if errmt != nil {
		t.FailNow()
	}

	eventgroup, erreg := apitest.CreateEventGroup(patrol.Context, project)
	if erreg != nil {
		t.FailNow()
	}

	events, errevents := apitest.CreateEvents(patrol.Context, eventgroup, 100)
	if errevents != nil {
		t.FailNow()
	}

	_ = events

	Convey("List events - unauthenticated user", t, func() {
		session := apitest.NewSession()
		request := session.Request("GET", settings.ROUTE_EVENTS_EVENT_LIST, "project_id", project.ID.String(), "eventgroup_id", eventgroup.ID.String())
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("List events - authenticated user", t, func() {
		session := apitest.NewSession().WithNewUser()
		request := session.Request("GET", settings.ROUTE_EVENTS_EVENT_LIST, "project_id", project.ID.String(), "eventgroup_id", eventgroup.ID.String())
		So(request.Do().Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("List events - member", t, func() {
		session := apitest.NewSession().WithUser(user)
		request := session.Request("GET", settings.ROUTE_EVENTS_EVENT_LIST, "project_id", project.ID.String(), "eventgroup_id", eventgroup.ID.String())
		request.SetValue("page", "2").Do()
		apitest.PrettyPrint("this is response %+v", request.Response().Body.String())

		So(request.Response().Code, ShouldEqual, http.StatusOK)
	})

}
