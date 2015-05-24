package teams

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamList(t *testing.T) {

	apitest.Setup()

	Convey("List teams for unauthorized user", t, func() {
		session := apitest.NewSession()
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAM_LIST)
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("List teams for authorized user", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
		})
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAM_LIST)
		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})

	Convey("List teams for superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(u *models.User) {
			u.IsSuperuser = true
			u.IsActive = true
		})
		request := session.Request("GET", settings.ROUTE_TEAMS_TEAM_LIST)
		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})

	Convey("List teams metadata for unauthorized user", t, func() {
		session := apitest.NewSession()
		request := session.Request("OPTIONS", settings.ROUTE_TEAMS_TEAM_LIST).Do()
		So(request.Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("List teams metadata for superuser", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
		})

		request := session.Request("OPTIONS", settings.ROUTE_TEAMS_TEAM_LIST).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		md, err := metadata.FromBytes(request.Response().Body.Bytes())
		So(err, ShouldBeNil)
		So(md.HasAction(metadata.ActionCreate), ShouldBeTrue)
		So(md.HasAction(metadata.ActionRetrieve), ShouldBeTrue)
	})

	Convey("List teams metadata for ordinary user", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
		})

		request := session.Request("OPTIONS", settings.ROUTE_TEAMS_TEAM_LIST).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		md, err := metadata.FromBytes(request.Response().Body.Bytes())
		So(err, ShouldBeNil)
		So(md.HasAction(metadata.ActionCreate), ShouldBeFalse)
		So(md.HasAction(metadata.ActionRetrieve), ShouldBeTrue)

	})
}
