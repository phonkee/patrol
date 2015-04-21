package common

import (
	"net/http"

	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/rest/response"
)

// This view has only support for HEAD method, useful for monitoring api server
type MonitorAPIView struct {
	core.JSONView
}

func (c *MonitorAPIView) HEAD(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusOK).Write(w, r)
}
