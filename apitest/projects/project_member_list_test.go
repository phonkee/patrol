package projects

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/phonkee/patrol"
// 	"github.com/phonkee/patrol/apitest"
// 	"github.com/phonkee/patrol/models"
// 	"github.com/phonkee/patrol/rest/metadata"
// 	"github.com/phonkee/patrol/settings"
// 	"github.com/phonkee/patrol/types"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestProjectMemberList(t *testing.T) {

// 	apitest.Setup()

// 	sudo, errore := apitest.CreateUser(patrol.Context, func(u *models.User) {
// 		u.IsSuperuser = true
// 	})
// 	if errore != nil {
// 		t.FailNow()
// 	}

// 	Convey("Test unauthenticated user", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession()
// 		request := session.Request("GET", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
// 		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
// 	})

// 	Convey("Test authenticated, nonmember", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession().WithNewUser()
// 		request := session.Request("GET", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
// 		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
// 	})

// 	Convey("Test authenticated, member", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession().WithNewUser(func(user *models.User) {
// 			user.IsSuperuser = true
// 		})

// 		tm := models.NewTeamMember(func(tm *models.TeamMember) {
// 			tm.TeamID = project.TeamID
// 			tm.UserID = sudo.ID.ToForeignKey()
// 			tm.Type = models.MEMBER_TYPE_MEMBER
// 		})
// 		err := tm.Insert(patrol.Context)
// 		So(err, ShouldBeNil)
// 		request := session.Request("GET", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()

// 		response := struct {
// 			Result []struct {
// 				ID types.ForeignKey `json:"id"`
// 			} `json:"result"`
// 			ResultSize int `json:"result_size"`
// 		}{}

// 		request.Scan(&response)
// 		So(response.ResultSize, ShouldEqual, 1)
// 		So(len(response.Result), ShouldEqual, 1)

// 	})

// 	Convey("Test metadata - unauthenticated user", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession()
// 		request := session.Request("OPTIONS", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
// 		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
// 	})

// 	Convey("Test metadata - authenticated non member", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession().WithNewUser()
// 		request := session.Request("OPTIONS", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
// 		So(request.Response().Code, ShouldEqual, http.StatusForbidden)
// 	})

// 	Convey("Test metadata - authenticated member", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession().WithNewUser()

// 		tm := models.NewTeamMember(func(tm *models.TeamMember) {
// 			tm.TeamID = project.TeamID
// 			tm.UserID = sudo.ID.ToForeignKey()
// 			tm.Type = models.MEMBER_TYPE_MEMBER
// 		})
// 		err := tm.Insert(patrol.Context)
// 		So(err, ShouldBeNil)

// 		request := session.Request("OPTIONS", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
// 		So(request.Response().Code, ShouldEqual, http.StatusOK)

// 		md, errmd := metadata.FromBytes(request.Response().Body.Bytes())
// 		So(errmd, ShouldBeNil)

// 		So(md.HasAction(metadata.ActionRetrieve), ShouldBeTrue)
// 		So(md.HasAction(metadata.ActionCreate), ShouldBeFalse)
// 	})

// 	Convey("Test list members metadata - admin member", t, func() {
// 		project, errcp := apitest.CreateProject(patrol.Context, sudo)
// 		So(errcp, ShouldBeNil)

// 		session := apitest.NewSession().WithNewUser()

// 		tm := models.NewTeamMember(func(tm *models.TeamMember) {
// 			tm.TeamID = project.TeamID
// 			tm.UserID = session.User().ID.ToForeignKey()
// 			tm.Type = models.MEMBER_TYPE_ADMIN
// 		})
// 		err := tm.Insert(patrol.Context)
// 		So(err, ShouldBeNil)

// 		request := session.Request("OPTIONS", settings.ROUTE_PROJECTS_PROJECTMEMBER_LIST, "project_id", project.ID.String()).Do()
// 		So(request.Response().Code, ShouldEqual, http.StatusOK)

// 		md, errmd := metadata.FromBytes(request.Response().Body.Bytes())
// 		So(errmd, ShouldBeNil)

// 		So(md.HasAction(metadata.ActionRetrieve), ShouldBeTrue)
// 		So(md.HasAction(metadata.ActionCreate), ShouldBeTrue)
// 	})
// }
