/* Models for auth plugin
 */
package models

import (
	"net/http"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/dgrijalva/jwt-go"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

/* User model
 */
type User struct {
	Model
	Username    string            `db:"username" json:"username"`
	Email       string            `db:"email" json:"email" validator:"email"`
	Password    string            `db:"password" json:"-"`
	Name        string            `db:"name" json:"name"`
	IsActive    bool              `db:"is_active" json:"is_active"`
	IsSuperuser bool              `db:"is_superuser" json:"is_superuser"`
	DateAdded   time.Time         `db:"date_added" json:"date_added"`
	LastLogin   time.Time         `db:"last_login" json:"last_login"`
	Permissions types.StringSlice `db:"permissions" json:"permissions"`
}

// returns all columns except of primary key
func (u *User) Columns() []string {
	return []string{
		"username", "email", "password", "name", "is_active",
		"is_superuser", "date_added", "last_login", "permissions",
	}
}
func (u *User) Values() []interface{} {
	return []interface{}{
		u.Username, u.Email, u.Password, u.Name, u.IsActive,
		u.IsSuperuser, u.DateAdded, u.LastLogin, u.Permissions,
	}
}
func (u *User) String() string { return "auth:user:" + u.ID.String() }
func (u *User) Table() string  { return AUTH_USER_DB_TABLE }

/*
CRUD
*/
func (u *User) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, u); err != nil {
		return
	}

	return Cache(ctx, u.String(), u)
}

func (u *User) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	if changed, err = DBUpdate(ctx, u, fields...); err != nil {
		return
	}

	if err = Cache(ctx, u.String(), u); err != nil {
		return
	}

	return
}

func (u *User) Delete(ctx *context.Context) (err error) {
	cacheKey := u.String()
	if err = DBDelete(ctx, u); err != nil {
		return err
	}

	if err = RemoveCached(ctx, cacheKey); err != nil {
		return
	}
	return
}

// Cleans fields (trimspace and other)
func (u *User) Clean() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
}

func (u *User) Validate(ctx *context.Context) (result *validator.Result, err error) {
	// create validator
	v := validator.New().Add("email", validator.ValidateEmail())

	result = v.Validate(u)

	if u.ID == 0 {
		manager := u.Manager(ctx)
		user := manager.NewUser()

		if err = manager.Get(user, manager.QueryFilterEmail(u.Email)); err != nil {
			if err == ErrObjectDoesNotExists {
				// supress error this is ok state
				err = nil
			}
			// } else {
			// 	rseult.AddPostgresError(err)
			// }
		} else {
			result.AddFieldError("email", ErrObjectAlreadyExists)
		}
	}

	return result, nil
}

// sets password to user
func (u *User) SetPassword(password string) error {
	result, err := utils.HashPassword(password, settings.SETTINGS_SECRET_KEY, settings.SETTINGS_BCRYPT_COST)
	if err != nil {
		return err
	}
	u.Password = result
	return nil
}

// Verifies user password
func (u *User) VerifyPassword(password string) (bool, error) {
	return utils.VerifyHashedPassword(u.Password, password, settings.SETTINGS_SECRET_KEY)
}

// returns user manager
func (u *User) Manager(ctx *context.Context) *UserManager {
	return NewUserManager(ctx)
}

/*
 Migrations for auth module
*/
const (
	MIGRATION_AUTH_USER_INITIAL_ID = "initial-migration-auth-user"
	MIGRATION_AUTH_USER_INITIAL    = `CREATE TABLE ` + AUTH_USER_DB_TABLE + `
    (
        id bigserial PRIMARY KEY,
        username character varying(32) NOT NULL UNIQUE,
        email character varying(255) NOT NULL UNIQUE,
        password character varying(60) NOT NULL,
        name character varying(255) NOT NULL,
        is_active boolean,
        is_superuser boolean,
        date_added timestamp with time zone NOT NULL,
        last_login timestamp with time zone,
        permissions character varying(64) array
    )`
)

/*
Auth managers
*/
const (
	USER_ID_TOKEN_KEY    = "user_id"
	EXPIRATION_TOKEN_KEY = "exp"
)

// Constructor function for new UserManager
func NewUserManager(context *context.Context) *UserManager {
	um := &UserManager{context: context}
	return um
}

// UserManager handles methods for list of users
type UserManager struct {
	Manager
	context *context.Context
}

// Returns new User with default values
func NewUser(funcs ...func(*User)) (user *User) {
	user = &User{
		DateAdded: utils.NowTruncated(),
	}
	for _, f := range funcs {
		f(user)
	}
	return
}

// Returns new User with default values
func (u *UserManager) NewUser(funcs ...func(*User)) (user *User) {
	return NewUser(funcs...)
}

func (u *UserManager) NewUserList() []*User {
	return []*User{}
}

// select without paging
func (u *UserManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*User)

	return DBFilter(u.context, AUTH_USER_DB_TABLE+".*", AUTH_USER_DB_TABLE, !safe, target, qfs...)
}

/* Filters projects from database
qfs is list of utils.QueryFuncs - that are functions that alter query builder
*/
func (u *UserManager) FilterPaged(target interface{}, paging *utils.Paging, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(u.context, AUTH_USER_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, u.QueryFilterPaging(paging))

	_, safe := target.([]*User)

	return DBFilter(u.context, AUTH_USER_DB_TABLE+".*", AUTH_USER_DB_TABLE, !safe, target, qfs...)
}

/* Returns project by Query filter funcs
 */
func (u *UserManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*User)

	return DBGet(u.context, "*", AUTH_USER_DB_TABLE, !safe, target, qfs...)
}

/* Returns project by ID
@TODO: cache only if target is *User
*/
func (u *UserManager) GetByID(target interface{}, id types.Keyer) (err error) {
	handleNilPointer(target)

	cacheKey := u.NewUser(func(user *User) { user.SetPrimaryKey(id) }).String()

	// found cached
	if err = GetCached(u.context, cacheKey, target); err == nil {
		return
	}

	if err = u.Get(target, u.QueryFilterID(id)); err != nil {
		return
	}

	// if it's user we can safely cache
	if _, ok := target.(*User); ok {
		if err = Cache(u.context, cacheKey, target); err != nil {
			return
		}
	}

	return
}

const (
	AUTH_TOKEN_CONTEXT_KEY = "AUTH:AUTH_TOKEN"
	AUTH_USER_CONTEXT_KEY  = "AUTH:AUTH_USER"
)

// Returns jwt Token from request
func (u *UserManager) GetAuthToken(r *http.Request) (token *jwt.Token, err error) {
	tokenString := r.Header.Get("Authorization")

	if len(tokenString) > 7 {
		tokenString = tokenString[7:]
	}

	secretKey := u.context.Get(context.SECRET_KEY).(string)

	if token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}); err != nil {
		return
	}

	return
}

// Returns User from request
func (u *UserManager) GetAuthUser(user *User, r *http.Request) (err error) {
	handleNilPointer(user)

	var token *jwt.Token

	// get token first
	if token, err = u.GetAuthToken(r); err != nil {
		return
	}

	uid := token.Claims[USER_ID_TOKEN_KEY]
	id := types.PrimaryKey(uid.(float64))

	if err = u.GetByID(user, id); err != nil {
		return err
	}

	return
}

// logs in user and returns token
func (u *UserManager) Login(user *User) (t string, err error) {
	if user.ID == 0 {
		err = ErrObjectDoesNotExists
		return
	}
	// Inactive user cannot be logged in
	if !user.IsActive {
		err = ErrCannotLoginUser
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims[USER_ID_TOKEN_KEY] = user.ID
	token.Claims[EXPIRATION_TOKEN_KEY] = time.Now().Add(time.Hour * 24).Unix()

	// get secret key from context
	secretKey := u.context.Get(context.SECRET_KEY).(string)
	t, err = token.SignedString([]byte(secretKey))
	return
}

func (u *UserManager) QueryFilterEmail(email string) utils.QueryFunc {
	return utils.QueryFilterWhere("email ILIKE ?", email)
}

func (u *UserManager) QueryFilterUsername(username string) utils.QueryFunc {
	return utils.QueryFilterWhere("username ILIKE ?", username)
}

func (u *UserManager) QueryFilterIsActive() utils.QueryFunc {
	return utils.QueryFilterWhere("is_active = ?", false)
}

/*
Various test methods
*/

func (u *UserManager) MakeTestUser(password string) func(*User) {
	return func(user *User) {
		user.Email = utils.RandomString(20) + randomdata.Email()
		user.Username = utils.RandomString(20)
		user.Name = randomdata.FullName(randomdata.Male)
		user.SetPassword(password)
		user.IsActive = true
	}
}

func (u *UserManager) GetTestUnsavedUser(callbacks ...func(*User)) (user *User, err error) {
	user = u.NewUser(func(us *User) {
		us.Email = utils.RandomString(20) + randomdata.Email()
		us.Username = utils.RandomString(20)
		us.Name = randomdata.FullName(randomdata.Male)
		us.SetPassword("")
	})
	for _, callback := range callbacks {
		callback(user)
	}
	return
}

func (u *UserManager) GetTestUser(callbacks ...func(*User)) (user *User, err error) {
	if user, err = u.GetTestUnsavedUser(callbacks...); err != nil {
		return
	}

	err = user.Insert(u.context)
	return
}

/*
Returns user with jwt token
*/
func (u *UserManager) GetTestUserWithToken(callbacks ...func(*User)) (user *User, token string, err error) {
	if user, err = u.GetTestUser(callbacks...); err != nil {
		return
	}
	if token, err = u.Login(user); err != nil {
		return
	}

	return
}
