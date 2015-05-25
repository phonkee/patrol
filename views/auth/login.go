package auth

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/settings"
)

type AuthLoginAPIView struct {
	core.JSONView

	// store callback for signal
	LoginSignal func(user *models.User) error

	// context
	context *context.Context
}

func (a *AuthLoginAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	a.context = a.GetContext(r)
	return
}

// Options request
func (a *AuthLoginAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	md := metadata.New("Login user")
	md.Action("POST").From(serializers.AuthLoginSerializer{})
	response.New().Raw(md).Write(w, r)
}

/* POST method for login
 */
func (l *AuthLoginAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	serializer := serializers.AuthLoginSerializer{}

	// unmarshal request body to serializers.AuthLoginSerializer
	if err = l.context.Bind(&serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	vr := serializer.Validate(l.context)

	if !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	user := models.NewUser()
	token := ""
	if user, token, err = serializer.Login(l.context); err != nil {
		switch err {
		case serializers.ErrUsernamePassword:
			vr.AddUnboundError(err)
			response.New(http.StatusUnauthorized).Error(vr).Write(w, r)
		case serializers.ErrInternalServerError:
			fallthrough
		default:
			response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		}
		return
	}

	// send signal
	l.LoginSignal(user)

	response.New(http.StatusOK).Header(settings.AUTH_TOKEN_HEADER_NAME, token).Write(w, r)
}
