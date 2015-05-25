package common

import (
	"net/http"

	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
)

// This view has only support for HEAD method, useful for monitoring api server
type MonitorAPIView struct {
	views.APIView
}

func (c *MonitorAPIView) HEAD(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusOK).Write(w, r)
}
