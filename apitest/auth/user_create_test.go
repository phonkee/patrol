package auth

import (
	"net/http"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
	"github.com/phonkee/patrol/views/auth"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthCreateUser(t *testing.T) {

	apitest.Setup()

	Convey("Create user - unauthorized user", t, func() {
		session := apitest.NewSession()
		request := session.Request("POST", settings.ROUTE_AUTH_USER_LIST).StringBody("{}")
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Create user - authorized user", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
			user.IsActive = true
		})
		request := session.Request("POST", settings.ROUTE_AUTH_USER_LIST).StringBody("{}")
		So(request.Do().Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Create user, valid data - superuser", t, func() {
		password := utils.RandomString(20)
		serializer := auth.UserCreateSerializer{
			Email:          utils.RandomString(20) + randomdata.Email(),
			Username:       utils.RandomString(20),
			Name:           randomdata.FullName(randomdata.Male),
			Password:       password,
			PasswordRetype: password,
			IsSuperuser:    false,
			IsActive:       true,
		}

		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = true
			user.IsActive = true
		})

		request := session.Request("POST", settings.ROUTE_AUTH_USER_LIST).JSONBody(serializer).Do()
		So(request.Response().Code, ShouldEqual, http.StatusCreated)

		rerequest := session.Request("POST", settings.ROUTE_AUTH_USER_LIST).JSONBody(serializer).Do()
		So(rerequest.Response().Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("Create user, invalid data - superuser", t, func() {
		// @TODO: add invalid data
	})

}
