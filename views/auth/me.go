package auth

import (
	"net/http"

	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
)

/*
Auth Me view.
returns information about user by jwt
*/

type AuthMeAPIView struct {
	core.JSONView
}

func (a *AuthMeAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error

	context := a.Context(r)

	um := models.NewUserManager(context)
	user := um.NewUser()
	if err = um.GetAuthUser(user, r); err != nil {
		switch err {
		case models.ErrObjectDoesNotExists:
			response.New(http.StatusNotFound).Write(w, r)
			return
		}
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(user).Write(w, r)
}
