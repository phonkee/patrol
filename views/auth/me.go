package auth

import (
	"net/http"

	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/views/mixins"
)

/*
Auth Me view.
returns information about user by jwt
*/

type AuthMeAPIView struct {
	core.JSONView
	mixins.AuthUserMixin
}

func (a *AuthMeAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error

	user := models.NewUser()
	if err = a.GetAuthUser(user, w, r); err != nil {
		return
	}

	response.New(http.StatusOK).Result(user).Write(w, r)
}
