package mixins

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/views"
)

/*
AuthUserMixin provides shortcut methods for auth user
*/
type AuthUserMixin struct{}

/*
	Gets auth user

	If any error occures, writes response and returns error
*/
func (a *AuthUserMixin) GetAuthUser(user *models.User, w http.ResponseWriter, r *http.Request) (err error) {
	var ctx *context.Context

	if ctx, err = context.Get(r); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return views.ErrInternalServerError
	}

	if err = models.NewUserManager(ctx).GetAuthUser(user, r); err != nil {
		response.New(http.StatusUnauthorized).Write(w, r)
		return views.ErrUnauthorized
	}

	return
}
