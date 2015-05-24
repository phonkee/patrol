package auth

import (
	"fmt"
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
)

/*
Constructor that returns New UserListAPIView
*/
func NewUserListAPIView() *UserListAPIView {
	return &UserListAPIView{}
}

/*
User list view
	GET - list of users
	POST - create new user
*/
type UserListAPIView struct {
	core.JSONView

	// store context
	context *context.Context

	user *models.User
}

// Before every method
func (u *UserListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	u.context = u.Context(r)

	u.user = models.NewUser()
	if err = u.user.Manager(u.context).GetAuthUser(u.user, r); err != nil {
		return
	}

	if !u.user.IsSuperuser {
		response.New(http.StatusForbidden).Write(w, r)
		return fmt.Errorf("forbidden")
	}

	return
}

// returns available metadata
func (u *UserListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {

	// @TODO: check permissions

	md := metadata.New("User list").SetDescription("User list endpoint")
	md.ActionRetrieve().From(UserListItemSerializer{})

	response.New(http.StatusOK).Metadata(md).Write(w, r)
}

/*
Retrieve list of registered users
*/
func (u *UserListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error
	manager := models.NewUserManager(u.context)

	paging := manager.NewPaging()

	// @TODO: check permissions

	list := []*UserListItemSerializer{}
	if err = manager.FilterPaged(&list, paging); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(list).Paging(paging).Write(w, r)
	return
}

/*
Create new user
*/
func (u *UserListAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	serializer := &UserCreateSerializer{}
	if err = u.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	// invalid
	if vr := serializer.Validate(u.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	// populate user with serializer data
	user := models.NewUser(func(us *models.User) {
		us.Username = serializer.Username
		us.Email = serializer.Email
		us.Name = serializer.Name
		us.IsActive = serializer.IsActive
		us.IsSuperuser = serializer.IsSuperuser
		us.SetPassword(serializer.Password)
	})

	if err = user.Insert(u.context); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	// prepare response
	result := &UserListItemSerializer{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Name:        user.Name,
		IsActive:    user.IsActive,
		IsSuperuser: user.IsSuperuser,
		DateAdded:   user.DateAdded,
		LastLogin:   user.LastLogin,
	}

	// clear password - no need to send it through wire
	response.New(http.StatusCreated).Result(result).Write(w, r)
	return
}
