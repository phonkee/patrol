package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/alice"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/settings"
)

/*
Middlewares
*/
// Recovery middleware for panic
func RecoveryMiddleware() alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					if settings.DEBUG {
						stack := debug.Stack()
						fmt.Printf("Request panicked `%v` with %s\n", err, string(stack))
					}

					response.New(http.StatusInternalServerError).Error(err).Write(w, r)

					return
				}
			}()

			// server request
			h.ServeHTTP(w, r)
		})
	}
}

// Recovery middleware for panic
func RequestLogMiddleware() alice.Constructor {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			cont, _ := context.Get(r)
			manager := models.NewRequestManager(cont)

			h.ServeHTTP(w, r)

			// Log request
			manager.LogRequest(r, start)
		})
	}
}

func ContextMiddleware(con *context.Context) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// create context
			newc := con.WithRequest(r)
			models.NewRequestManager(newc).Status(http.StatusOK)

			h.ServeHTTP(w, r)

			// clear context
			context.Clear(r)
		})
	}
}
