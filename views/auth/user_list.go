package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

/*
User list item serializer
*/
type UserListItemSerializer struct {
	ID          types.PrimaryKey `db:"id" json:"id"`
	Username    string           `db:"username" json:"username"`
	Email       string           `db:"email" json:"email"`
	Name        string           `db:"name" json:"name"`
	IsActive    bool             `db:"is_active" json:"is_active"`
	IsSuperuser bool             `db:"is_superuser" json:"is_superuser"`
	DateAdded   time.Time        `db:"date_added" json:"date_added"`
	LastLogin   time.Time        `db:"last_login" json:"last_login"`
}

/*
User create serializer
*/
type UserCreateSerializer struct {
	Username       string `json:"username" validator:"username"`
	Email          string `json:"email" validator:"email"`
	Name           string `json:"name" validator:"name"`
	IsActive       bool   `json:"is_active"`
	IsSuperuser    bool   `json:"is_superuser"`
	Password       string `json:"password" validator:"password"`
	PasswordRetype string `json:"password_retype" validator:"password"`
}

/*
Cleans values (trims spaces)
*/
func (u *UserCreateSerializer) Clean() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
	u.Password = strings.TrimSpace(u.Password)
	u.PasswordRetype = strings.TrimSpace(u.PasswordRetype)
}

/*
Validates create serializer
*/
func (u *UserCreateSerializer) Validate(context *context.Context) *validator.Result {
	v := validator.New()
	v["name"] = models.ValidateUserName()
	v["email"] = validator.ValidateEmail()
	v["username"] = models.ValidateUserUsername()
	v["password"] = models.ValidatePassword()
	result := v.Validate(u)

	if u.Password != u.PasswordRetype {
		result.AddFieldError("password_retype", errors.New("passwords_dont_match"))
	}

	um := models.NewUserManager(context)
	// check duplicate username
	if !result.HasFieldErrors("username") {
		if err := um.Get(um.NewUser(), utils.QueryFilterWhere("username = ?", u.Username)); err == nil {
			result.AddFieldError("username", errors.New("already_exists"))
		}
	}

	// check duplicaate email
	if !result.HasFieldErrors("email") {
		if err := um.Get(um.NewUser(), utils.QueryFilterWhere("email = ?", u.Email)); err == nil {
			result.AddFieldError("email", errors.New("already_exists"))
		}
	}
	return result
}

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
