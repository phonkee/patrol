package middlewares

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
)

/*
Middlewares
*/
func AuthTokenValidMiddleware() alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				err error
			)

			ctx, _ := context.Get(r)

			manager := models.NewUserManager(ctx)

			user := manager.NewUser()
			if err = manager.GetAuthUser(user, r); err != nil {
				response.New(http.StatusUnauthorized).Error(err).Write(w, r)
				return
			}

			// inactive user sorry for this!
			// @TODO: is Unauthorized status appropriate?
			if !user.IsActive {
				response.New(http.StatusUnauthorized).Write(w, r)
				return
			}

			// serve next
			h.ServeHTTP(w, r)
		})
	}
}
