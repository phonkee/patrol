package auth

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthUserChangePassword(t *testing.T) {

	apitest.Setup()

	Convey("Test change own password", t, func() {

		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
			user.IsActive = true
		})
		password := utils.RandomString(32)
		request := session.Request("POST", settings.ROUTE_AUTH_USER_CHANGE_PASSWORD, "user_id", session.User().ID.String())
		request.JSONBody(map[string]interface{}{
			"password": password,
			"retype":   password,
		}).Do()
		So(request.Response().Code, ShouldEqual, http.StatusOK)

		loginsession := apitest.NewSession()
		request = loginsession.Request("POST", settings.ROUTE_AUTH_LOGIN).JSONBody(map[string]interface{}{
			"username": session.User().Username,
			"password": password,
		})
		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})

	Convey("Test change own password - invalid", t, func() {
		session := apitest.NewSession().WithNewUser(func(user *models.User) {
			user.IsSuperuser = false
			user.IsActive = true
		})
		request := session.Request("POST", settings.ROUTE_AUTH_USER_CHANGE_PASSWORD, "user_id", session.User().ID.String())
		request.JSONBody(map[string]interface{}{
			"password": "password",
			"retype":   "retype",
		})
		So(request.Do().Response().Code, ShouldEqual, http.StatusBadRequest)

		request = session.Request("POST", settings.ROUTE_AUTH_USER_CHANGE_PASSWORD, "user_id", session.User().ID.String())
		request.JSONBody(map[string]interface{}{
			"password": "a",
			"retype":   "a",
		})
		So(request.Do().Response().Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("Test change others password - unauthenticated user", t, func() {
		user, err := apitest.CreateUser(patrol.Context)
		So(err, ShouldBeNil)

		session := apitest.NewSession(patrol.Context)
		request := session.Request("POST", settings.ROUTE_AUTH_USER_CHANGE_PASSWORD, "user_id", user.ID.String())
		request.JSONBody(map[string]interface{}{
			"password": "password",
			"retype":   "password",
		})
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Test change others password - normal user", t, func() {
		other, err := apitest.CreateUser(patrol.Context, func(u *models.User) {
			u.IsSuperuser = false
			u.IsActive = true
		})
		So(err, ShouldBeNil)

		session := apitest.NewSession().WithNewUser(func(u *models.User) {
			u.IsSuperuser = false
			u.IsActive = true
		})

		password := utils.RandomString(32)
		request := session.Request("POST", settings.ROUTE_AUTH_USER_CHANGE_PASSWORD, "user_id", other.ID.String())
		request.JSONBody(map[string]interface{}{
			"password": password,
			"retype":   password,
		})

		So(request.Do().Response().Code, ShouldEqual, http.StatusForbidden)
	})

	Convey("Test change others password - superuser", t, func() {
		other, err := apitest.CreateUser(patrol.Context, func(u *models.User) {
			u.IsSuperuser = false
			u.IsActive = true
		})
		So(err, ShouldBeNil)

		session := apitest.NewSession().WithNewUser(func(u *models.User) {
			u.IsSuperuser = true
			u.IsActive = true
		})

		password := utils.RandomString(32)
		request := session.Request("POST", settings.ROUTE_AUTH_USER_CHANGE_PASSWORD, "user_id", other.ID.String())
		request.JSONBody(map[string]interface{}{
			"password": password,
			"retype":   password,
		})

		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})
}
