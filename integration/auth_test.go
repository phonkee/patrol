package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuth(t *testing.T) {
	if err := patrol.Setup(); err != nil {
		fmt.Printf("patrol: setup failed with error: %s", err)
	}

	patrol.Run([]string{"migrate"})

	Convey("Test non existing Login", t, func() {
		session := NewRequestsSession(patrol.Context)
		So(
			session.Request("POST", settings.ROUTE_AUTH_LOGIN).
				StringBody(`{"username":"nonexisting", "password": "also"}`).
				Do().
				Response().
				Code,
			ShouldEqual,
			http.StatusUnauthorized,
		)
	})

	Convey("Test existing Login", t, func() {
		password := "password"
		user, err := CreateUser(patrol.Context, func(user *models.User) {
			user.SetPassword(password)
		})
		_, _ = user, err

		body := `{"username":"` + user.Username + `", "password": "` + password + `"}`
		r := NewAPIRequest("POST", settings.ROUTE_AUTH_LOGIN, NewBody(body))
		w, _ := GetAPIResponse(r, nil)

		So(w.Code, ShouldEqual, http.StatusOK)

		token := w.Header().Get(settings.AUTH_TOKEN_HEADER_NAME)
		So(len(token), ShouldBeGreaterThan, 10)

		r2 := NewAPIRequest("GET", settings.ROUTE_AUTH_ME, nil)
		r2.Header.Set(settings.AUTH_TOKEN_HEADER_NAME, token)

		w2, _ := GetAPIResponse(r2, nil)
		So(w2.Code, ShouldEqual, http.StatusOK)

	})

	Convey("Test Login bad request", t, func() {
		body := `{"username":"asdf", "password": "adsf"`
		r := NewAPIRequest("POST", settings.ROUTE_AUTH_LOGIN, NewBody(body))
		w, _ := GetAPIResponse(r, nil)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})

}
