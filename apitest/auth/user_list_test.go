package auth

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthUserList(t *testing.T) {

	apitest.Setup()

	Convey("Test retrieve user list by non superuser", t, func() {
		session := apitest.NewSession(patrol.Context).WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
		})
		request := session.Request("GET", settings.ROUTE_AUTH_USER_LIST)
		response := request.Do().Response()
		So(response.Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Test retrieve user list by superuser", t, func() {
		session := apitest.NewSession(patrol.Context).WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
		})
		request := session.Request("GET", settings.ROUTE_AUTH_USER_LIST)
		response := request.Do().Response()
		So(response.Code, ShouldEqual, http.StatusOK)
	})

}
