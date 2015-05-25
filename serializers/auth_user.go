/*
Auth serializers

serializers for all auth api views
*/
package serializers

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

/*
AuthUserDetailSerializer
	is simplified user serializer which hides sensitive information such as
	password
*/
type AuthUserDetailSerializer struct {
	ID          types.PrimaryKey `db:"id"           json:"id"`
	Username    string           `db:"username"     json:"username"`
	Email       string           `db:"email"        json:"email" validator:"email"`
	Name        string           `db:"name"         json:"name"`
	IsActive    bool             `db:"is_active"    json:"is_active"`
	IsSuperuser bool             `db:"is_superuser" json:"is_superuser"`
	DateAdded   time.Time        `db:"date_added"   json:"date_added"`
	LastLogin   time.Time        `db:"last_login"   json:"last_login"`
}

/*
	Updates data from user
*/
func (u *AuthUserDetailSerializer) From(user *models.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.Email = user.Email
	u.Name = user.Name
	u.IsActive = user.IsActive
	u.IsSuperuser = user.IsSuperuser
	u.DateAdded = user.DateAdded
	u.LastLogin = user.LastLogin
}

/*
AuthUserUpdateSerializer
Serializer to update auth user

*/
type AuthUserUpdateSerializer struct {
	Email string `db:"email" json:"email" validator:"email"`
	Name  string `db:"name"  json:"name"`

	// only superuser can change these
	IsActive    bool `db:"is_active"    json:"is_active"`
	IsSuperuser bool `db:"is_superuser" json:"is_superuser"`
}

/*
Clean - trims space from strings
*/
func (u *AuthUserUpdateSerializer) Clean() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
}

/*
Validate
@TODO: check duplicate email
*/
func (u AuthUserUpdateSerializer) Validate(context *context.Context) *validator.Result {
	v := validator.New()
	v["email"] = validator.ValidateEmail()
	v["name"] = models.ValidateUserName()
	return v.Validate(u)
}

/*
Save updates user
*/
func (a *AuthUserUpdateSerializer) Save(context *context.Context, user *models.User, updater *models.User) (result *AuthUserDetailSerializer, err error) {
	user.Email = a.Email
	user.Name = a.Name

	// fields to update
	fields := []string{"email", "name"}

	if updater.IsSuperuser {
		user.IsActive = a.IsActive
		user.IsSuperuser = a.IsSuperuser
		fields = append(fields, "is_active")
		fields = append(fields, "is_superuser")
	}

	// update user
	if _, err = user.Update(context, fields...); err != nil {
		return
	}

	result = &AuthUserDetailSerializer{}
	result.From(user)
	return
}

/*
User create serializer
*/
type AuthUserCreateSerializer struct {
	Username       string `json:"username"        validator:"username"`
	Email          string `json:"email"           validator:"email"`
	Name           string `json:"name"            validator:"name"`
	IsActive       bool   `json:"is_active"`
	IsSuperuser    bool   `json:"is_superuser"`
	Password       string `json:"password"        validator:"password"`
	PasswordRetype string `json:"password_retype" validator:"password"`
}

/*
Cleans values (trims spaces)
*/
func (u *AuthUserCreateSerializer) Clean() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
	u.Password = strings.TrimSpace(u.Password)
	u.PasswordRetype = strings.TrimSpace(u.PasswordRetype)
}

/*
Validates create serializer
*/
func (u *AuthUserCreateSerializer) Validate(context *context.Context) *validator.Result {
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
	Saves user to database
*/
func (a *AuthUserCreateSerializer) Save(context *context.Context) (result *AuthUserDetailSerializer, err error) {
	// populate user with serializer data
	user := models.NewUser(func(us *models.User) {
		us.Username = a.Username
		us.Email = a.Email
		us.Name = a.Name
		us.IsActive = a.IsActive
		us.IsSuperuser = a.IsSuperuser
		us.SetPassword(a.Password)
	})

	if err = user.Insert(context); err != nil {
		return
	}

	result = &AuthUserDetailSerializer{}
	result.From(user)

	return
}

/*
AuthUserChangePasswordSerializer
	Change password serializer
*/
type AuthUserChangePasswordSerializer struct {
	Password string `json:"password" validator:"password"`
	Retype   string `json:"retype"`
}

/*
	Trim spaces
*/
func (u *AuthUserChangePasswordSerializer) Clean() {
	u.Password = strings.TrimSpace(u.Password)
	u.Retype = strings.TrimSpace(u.Retype)
}

func (u *AuthUserChangePasswordSerializer) Validate(context *context.Context) (result *validator.Result) {
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
	Save method changes password for given user
*/
func (a *AuthUserChangePasswordSerializer) Save(context *context.Context, user *models.User) (result *AuthUserDetailSerializer, err error) {
	user.SetPassword(a.Password)
	if _, err = user.Update(context, "password"); err != nil {
		return
	}

	result = &AuthUserDetailSerializer{}
	result.From(user)
	return
}

/*
AuthLoginSerializer
*/
type AuthLoginSerializer struct {
	Username string `json:"username" validator:"username"`
	Password string `json:"password"`
}

/*
Cleans data
*/
func (a *AuthLoginSerializer) Clean() {
	a.Username = strings.TrimSpace(a.Username)
	a.Password = strings.TrimSpace(a.Password)
}

/*
Validates input data
*/
func (a *AuthLoginSerializer) Validate(context *context.Context) *validator.Result {
	v := validator.New()
	v["username"] = models.ValidateUserName()

	return v.Validate(a)
}

/*
	Performs login and returns user and token
*/
func (a *AuthLoginSerializer) Login(context *context.Context) (user *models.User, token string, err error) {
	usermanager := models.NewUserManager(context)

	// find user in db by username
	user = usermanager.NewUser()
	if err = usermanager.Get(user, usermanager.QueryFilterUsername(a.Username)); err != nil {
		err = ErrUsernamePassword
		return
	}

	// check password
	var ok bool
	if ok, err = user.VerifyPassword(a.Password); err != nil {
		err = ErrInternalServerError
		return
	} else {
		if !ok {
			err = ErrUsernamePassword
			return
		}
	}

	// create token
	if token, err = usermanager.Login(user); err != nil {
		return
	}

	// update last login
	user.LastLogin = utils.NowTruncated()
	if _, err = user.Update(context, "last_login"); err != nil {
		return
	}

	return
}
