package auth

import (
	"fmt"
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/types"
)

type UserChangePasswordAPIView struct {
	core.JSONView
	context  *context.Context
	authuser *models.User
	user     *models.User
}

func (u *UserChangePasswordAPIView) Before(rw http.ResponseWriter, r *http.Request) (err error) {
	u.context = u.Context(r)

	um := models.NewUserManager(u.context)
	u.authuser = models.NewUser()
	if err = um.GetAuthUser(u.authuser, r); err != nil {
		return
	}

	var id *types.ForeignKey

	if id, err = u.GetMuxVarForeignKey(r, "user_id"); err != nil {
		response.New(http.StatusInternalServerError).Write(rw, r)
		return
	}

	u.user = models.NewUser()
	if err = um.GetByID(u.user, id); err != nil {
		response.New(http.StatusForbidden).Write(rw, r)
		return
	}

	// check
	if !u.authuser.IsSuperuser {
		if u.authuser.ID != u.user.ID {
			response.New(http.StatusForbidden).Write(rw, r)
			return fmt.Errorf("forbidden")
		}
	}

	return
}

/*
	change password
*/
func (u *UserChangePasswordAPIView) POST(rw http.ResponseWriter, r *http.Request) {
	var err error
	serializer := &UserChangePasswordSerializer{}
	if err = u.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(rw, r)
		return
	}

	vr := serializer.Validate(u.context)
	if !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(rw, r)
		return
	}

	u.user.SetPassword(serializer.Password)
	if _, err = u.user.Update(u.context, "password"); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(rw, r)
		return
	}

	response.New(http.StatusOK).Write(rw, r)
	return
}
