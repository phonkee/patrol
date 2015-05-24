package projects

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/views/projects"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProjectMemberCreate(t *testing.T) {

	apitest.Setup()

	sudo, errore := apitest.CreateUser(patrol.Context, func(u *models.User) {
		u.IsSuperuser = true
	})
	if errore != nil {
		t.FailNow()
	}

	Convey("Test create member unauthenticated user", t, func() {
		project, errcp := apitest.CreateProject(patrol.Context, sudo)
		So(errcp, ShouldBeNil)

		session := apitest.NewSession()
		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String())

		serializer := &projects.ProjectMemberCreate{
			Type:   models.MEMBER_TYPE_MEMBER,
			UserID: types.PrimaryKey(1).ToForeignKey(),
		}
		request.JSONBody(serializer).Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Test create member user", t, func() {
		project, errcp := apitest.CreateProject(patrol.Context, sudo)
		So(errcp, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()
		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String())

		serializer := &projects.ProjectMemberCreate{
			Type:   models.MEMBER_TYPE_MEMBER,
			UserID: types.PrimaryKey(1).ToForeignKey(),
		}
		request.JSONBody(serializer).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Create member", t, func() {
		project, errcp := apitest.CreateProject(patrol.Context, sudo)
		So(errcp, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()
		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String())

		tm := models.NewTeamMember(func(tm *models.TeamMember) {
			tm.TeamID = project.TeamID
			tm.UserID = session.User().ID.ToForeignKey()
			tm.Type = models.MEMBER_TYPE_MEMBER
		})
		err := tm.Insert(patrol.Context)
		So(err, ShouldBeNil)

		other, errcreate := apitest.CreateUser(patrol.Context)
		So(errcreate, ShouldBeNil)

		serializer := &projects.ProjectMemberCreate{
			Type:   models.MEMBER_TYPE_MEMBER,
			UserID: other.ID.ToForeignKey(),
		}
		request.JSONBody(serializer).Do()
		So(request.Response().Code, ShouldEqual, http.StatusForbidden)

	})

	Convey("Create - admin member", t, func() {
		project, errcp := apitest.CreateProject(patrol.Context, sudo)
		So(errcp, ShouldBeNil)

		session := apitest.NewSession().WithNewUser()

		tm := models.NewTeamMember(func(tm *models.TeamMember) {
			tm.TeamID = project.TeamID
			tm.UserID = session.User().ID.ToForeignKey()
			tm.Type = models.MEMBER_TYPE_ADMIN
		})
		err := tm.Insert(patrol.Context)
		So(err, ShouldBeNil)

		_, errcreate2 := apitest.CreateUser(patrol.Context)
		So(errcreate2, ShouldBeNil)
		newuser, errcreate := apitest.CreateUser(patrol.Context)
		So(errcreate, ShouldBeNil)

		pmcser := &projects.ProjectMemberCreate{
			Type:   models.MEMBER_TYPE_MEMBER,
			UserID: newuser.ID.ToForeignKey(),
		}
		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String())
		request.JSONBody(pmcser).Do()
		So(request.Response().Code, ShouldEqual, http.StatusCreated)

		requestget := session.Request("GET", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
		response := struct {
			Result []struct {
				ID     types.ForeignKey `json:"id"`
				UserID types.ForeignKey `json:"user_id"`
			} `json:"result"`
			ResultSize int `json:"result_size"`
		}{}

		requestget.Scan(&response)

		found := false
		for _, item := range response.Result {
			if item.UserID == pmcser.UserID {
				found = true
			}
		}

		So(found, ShouldBeTrue)

	})

	Convey("Create superuser", t, func() {
		project2, errcp := apitest.CreateProject(patrol.Context, sudo)
		So(errcp, ShouldBeNil)

		session := apitest.NewSession().WithNewUser(func(uu *models.User) {
			uu.IsSuperuser = true
		})

		other, errcreate := apitest.CreateUser(patrol.Context)
		So(errcreate, ShouldBeNil)

		serializer := &projects.ProjectMemberCreate{
			Type:   models.MEMBER_TYPE_MEMBER,
			UserID: other.ID.ToForeignKey(),
		}
		request := session.Request("POST", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project2.ID.String())
		request.JSONBody(serializer).Do()

		So(request.Response().Code, ShouldEqual, http.StatusCreated)

		requestget := session.Request("GET", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project2.ID.String()).Do()
		response := struct {
			Result []struct {
				ID     types.ForeignKey `json:"id"`
				UserID types.ForeignKey `json:"user_id"`
			} `json:"result"`
			ResultSize int `json:"result_size"`
		}{}

		requestget.Scan(&response)

		found := false
		for _, item := range response.Result {
			if item.UserID == serializer.UserID {
				found = true
			}
		}

		So(found, ShouldBeTrue)

	})

}
