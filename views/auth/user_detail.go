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

/*
	User detail api view
*/
type UserDetailAPIView struct {
	core.JSONView

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
func (u *UserDetailAPIView) Before(rw http.ResponseWriter, r *http.Request) (err error) {
	u.context = u.Context(r)

	var id *types.ForeignKey

	if id, err = u.GetMuxVarForeignKey(r, "user_id"); err != nil {
		return
	}

	u.user = models.NewUser()
	if err = u.user.Manager(u.context).GetByID(u.user, id); err != nil {
		return
	}

	// get auth user to check permissions
	usermanager := models.NewUserManager(u.context)
	u.authuser = usermanager.NewUser()
	if err = usermanager.GetAuthUser(u.authuser, r); err != nil {
		response.New(http.StatusUnauthorized).Write(rw, r)
		return
	}

	// check permissions
	if !u.authuser.IsSuperuser {
		if u.authuser.ID != u.user.ID {
			response.New(http.StatusForbidden).Write(rw, r)
			return fmt.Errorf("forbidden")
		}
	}

	return nil
}

/*
	GET (retrieve) method
*/
func (u *UserDetailAPIView) GET(rw http.ResponseWriter, r *http.Request) {
	serializer := &UserDetailSerializer{
		ID:          u.user.ID,
		Username:    u.user.Username,
		Email:       u.user.Email,
		Name:        u.user.Name,
		IsActive:    u.user.IsActive,
		IsSuperuser: u.user.IsSuperuser,
	}
	response.New().Status(http.StatusOK).Result(serializer).Write(rw, r)
}

/*
	POST (update) request
*/
func (u *UserDetailAPIView) POST(rw http.ResponseWriter, r *http.Request) {
	serializer := &UserUpdateSerializer{}

	if err := u.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(rw, r)
		return
	}

	// validate struct
	if vr := serializer.Validate(u.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(rw, r)
		return
	}

	var err error

	u.user.Email = serializer.Email
	u.user.Name = serializer.Name

	// fields to update
	fields := []string{"email", "name"}

	if u.authuser.IsSuperuser {
		u.user.IsActive = serializer.IsActive
		u.user.IsSuperuser = serializer.IsSuperuser
		fields = append(fields, "is_active")
		fields = append(fields, "is_superuser")
	}

	// update user
	if _, err = u.user.Update(u.context, fields...); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(rw, r)
		return
	}

	seruser := &UserDetailSerializer{
		ID:          u.user.ID,
		Username:    u.user.Username,
		Email:       u.user.Email,
		Name:        u.user.Name,
		IsActive:    u.user.IsActive,
		IsSuperuser: u.user.IsSuperuser,
	}

	response.New(http.StatusOK).Result(seruser).Write(rw, r)
	return
}
