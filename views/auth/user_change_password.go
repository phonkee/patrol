package auth

import (
	"fmt"
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/views/mixins"
)

type UserChangePasswordAPIView struct {
	views.APIView
	mixins.AuthUserMixin

	context  *context.Context
	authuser *models.User
	user     *models.User
}

func (u *UserChangePasswordAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	u.context = u.GetContext(r)

	u.authuser = models.NewUser()
	if err = u.GetAuthUser(u.authuser, w, r); err != nil {
		return
	}

	var id types.PrimaryKey

	if id, err = rest.GetMuxVarPrimaryKey(r, "user_id"); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	u.user = models.NewUser()
	um := models.NewUserManager(u.context)
	if err = um.GetByID(u.user, id); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	// check
	if !u.authuser.IsSuperuser {
		if u.authuser.ID != u.user.ID {
			response.New(http.StatusForbidden).Write(w, r)
			return fmt.Errorf("forbidden")
		}
	}

	return
}

/*
	change password
*/
func (u *UserChangePasswordAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	serializer := &serializers.AuthUserChangePasswordSerializer{}
	if err = u.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	if vr := serializer.Validate(u.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	var result *serializers.AuthUserDetailSerializer

	if result, err = serializer.Save(u.context, u.user); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(result).Write(w, r)
	return
}
