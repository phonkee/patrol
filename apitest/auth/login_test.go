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

func TestAuthLogin(t *testing.T) {
	apitest.Setup()

	Convey("Test non existing Login", t, func() {
		session := apitest.NewSession(patrol.Context)
		body := `{"username":"nonexisting", "password": "also"}`
		request := session.Request("POST", settings.ROUTE_AUTH_LOGIN).StringBody(body)
		So(request.Do().Response().Code, ShouldEqual, http.StatusUnauthorized)
	})

	Convey("Test existing Login", t, func() {
		password := "password"
		session := apitest.NewSession(patrol.Context).WithNewUser(func(user *models.User) {
			user.IsActive = true
			user.SetPassword(password)
		})
		user := session.User()

		session = apitest.NewSession(patrol.Context)
		request := session.Request("POST", settings.ROUTE_AUTH_LOGIN).JSONBody(map[string]string{
			"username": user.Username,
			"password": password,
		})

		So(request.Do().Response().Code, ShouldEqual, http.StatusOK)
	})

	Convey("Test Login bad request", t, func() {
		session := apitest.NewSession(patrol.Context)
		request := session.Request("POST", settings.ROUTE_AUTH_LOGIN).StringBody(`{"username":"asdf", "password": "adsf"`)
		So(request.Do().Response().Code, ShouldEqual, http.StatusBadRequest)
	})
}
