package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

type UserDetailSerializer struct {
	ID          types.PrimaryKey `db:"id" json:"id"`
	Username    string           `db:"username" json:"username"`
	Email       string           `db:"email" json:"email" validator:"email"`
	Name        string           `db:"name" json:"name"`
	IsActive    bool             `db:"is_active" json:"is_active"`
	IsSuperuser bool             `db:"is_superuser" json:"is_superuser"`
	DateAdded   time.Time        `db:"date_added" json:"date_added"`
	LastLogin   time.Time        `db:"last_login" json:"last_login"`
}

func (u *UserDetailSerializer) FromUser(user *models.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.Email = user.Email
	u.Name = user.Name
	u.IsActive = user.IsActive
	u.IsSuperuser = user.IsSuperuser
	u.DateAdded = user.DateAdded
	u.LastLogin = user.LastLogin
}

type UserUpdateSerializer struct {
	Email string `db:"email" json:"email" validator:"email"`
	Name  string `db:"name" json:"name"`
	// warning only superuser can change these
	IsActive    bool `db:"is_active" json:"is_active"`
	IsSuperuser bool `db:"is_superuser" json:"is_superuser"`
}

func (u *UserUpdateSerializer) Clean() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
}

/*
	Validate
*/
func (u UserUpdateSerializer) Validate(context *context.Context) *validator.Result {
	val := validator.New()
	val["email"] = validator.ValidateEmail()
	val["name"] = models.ValidateUserName()
	return val.Validate(u)
}

/*
	Change password serializer
*/
type UserChangePasswordSerializer struct {
	Password string `json:"password" validator:"password"`
	Retype   string `json:"retype"`
}

/*
	Trim spaces
*/
func (u *UserChangePasswordSerializer) Clean() {
	u.Password = strings.TrimSpace(u.Password)
	u.Retype = strings.TrimSpace(u.Retype)
}

func (u UserChangePasswordSerializer) Validate(context *context.Context) (result *validator.Result) {
	val := validator.New()
	val["password"] = models.ValidatePassword()
	result = val.Validate(u)

	if !result.IsValid() {
		return
	}

	// validate equal
	if u.Password != u.Retype {
		result.AddUnboundError(fmt.Errorf("passwords_not_match"))
	}

	return
}

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
