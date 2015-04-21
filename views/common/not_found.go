package common

import (
	"net/http"

	"github.com/phonkee/patrol/rest/response"
)

// not found handler
// for proper request log this handler should be prepended by middleware
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.New(http.StatusNotFound).Write(w, r)
	}
}
