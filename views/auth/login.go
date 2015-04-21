package auth

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
)

// Serializers used in login
type LoginSerializer struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginAPIView struct {
	core.JSONView

	// store callback for signal
	LoginSignal func(user *models.User) error

	// context
	context *context.Context
}

func (a *AuthLoginAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	a.context = a.Context(r)
	return
}

// Options request
func (a *AuthLoginAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	md := metadata.New("Login user")
	md.Action("POST").From(LoginSerializer{})
	response.New().Raw(md).Write(w, r)
}

/* POST method for login
 */
func (l *AuthLoginAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	login := LoginSerializer{}

	context := l.Context(r)

	// prepare blank response
	response := response.New()

	// unmarshal request body to LoginSerializer
	if err = l.Unmarshal(r.Body, &login); err != nil {
		response.Status(http.StatusBadRequest).Write(w, r)
		return
	}

	usermanager := models.NewUserManager(context)

	// find user in db by username
	user := usermanager.NewUser()
	if err = usermanager.Get(user, usermanager.QueryFilterUsername(login.Username)); err != nil {
		response.Status(http.StatusUnauthorized).Error(err).Write(w, r)
		return
	}

	// check password
	var ok bool
	if ok, err = user.VerifyPassword(login.Password); err != nil {
		response.Status(http.StatusInternalServerError).Write(w, r)
		return
	} else {
		if !ok {
			response.Status(http.StatusUnauthorized).Write(w, r)
			return
		}
	}

	// create token
	var token string
	if token, err = usermanager.Login(user); err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	// update last login
	user.LastLogin = utils.NowTruncated()
	user.Update(context, "last_login")

	// send signal
	l.LoginSignal(user)

	response.Status(http.StatusOK).Header(settings.AUTH_TOKEN_HEADER_NAME, token).Write(w, r)
}
