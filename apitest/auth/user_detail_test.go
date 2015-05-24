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

func TestAuthUserDetail(t *testing.T) {

	apitest.Setup()

	Convey("Test user detail for non superuser", t, func() {
		session := apitest.NewSession(patrol.Context).WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
		})

		user, err := apitest.CreateUser(patrol.Context)
		So(err, ShouldBeNil)
		request := session.Request("GET", settings.ROUTE_AUTH_USER_DETAIL, "user_id", user.ID.String())
		So(request.Do().Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Test user detail for superuser", t, func() {
		session := apitest.NewSession(patrol.Context).WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
		})

		user, err := apitest.CreateUser(patrol.Context)
		So(err, ShouldBeNil)
		request := session.Request("GET", settings.ROUTE_AUTH_USER_DETAIL, "user_id", user.ID.String())
		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})

}
