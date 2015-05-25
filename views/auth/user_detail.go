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

/*
	User detail api view
*/
type UserDetailAPIView struct {
	views.APIView

	mixins.AuthUserMixin

	// store context
	context *context.Context

	// loaded user
	user *models.User

	// authenticated user
	authuser *models.User
}

/*
	Preload user for all methods

	@TODO: check permissions
		1. logged user is the same user as detail
		2. user is superuser
		3. edited user belongs to team that logged user is admin
*/
func (u *UserDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	u.context = u.GetContext(r)

	var id types.PrimaryKey

	if id, err = rest.GetMuxVarPrimaryKey(r, "user_id"); err != nil {
		return
	}

	usermanager := models.NewUserManager(u.context)
	u.user = models.NewUser()
	if err = u.user.Manager(u.context).GetByID(u.user, id); err != nil {
		return
	}

	// get auth user to check permissions
	u.authuser = usermanager.NewUser()
	if err = u.GetAuthUser(u.authuser, w, r); err != nil {
		return
	}

	// check permissions
	if !u.authuser.IsSuperuser {
		if u.authuser.ID != u.user.ID {
			response.New(http.StatusForbidden).Write(w, r)
			return fmt.Errorf("forbidden")
		}
	}

	return nil
}

/*
	GET (retrieve) method
*/
func (u *UserDetailAPIView) GET(rw http.ResponseWriter, r *http.Request) {
	serializer := &serializers.AuthUserDetailSerializer{}
	serializer.From(u.user)
	response.New().Status(http.StatusOK).Result(serializer).Write(rw, r)
}

/*
	POST (update) request
*/
func (u *UserDetailAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error
	serializer := &serializers.AuthUserUpdateSerializer{}

	if err = u.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	// validate struct
	if vr := serializer.Validate(u.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	var result *serializers.AuthUserDetailSerializer

	if result, err = serializer.Save(u.context, u.user, u.authuser); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(result).Write(w, r)
}
