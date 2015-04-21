package common

import (
	"net/http"

	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/parser"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/settings"
)

// Returns version of patrol
type VersionAPIView struct {
	core.JSONView
}

func (v *VersionAPIView) GET(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{
		"version":   settings.VERSION,
		"protocols": parser.Versions(),
	}

	response.New(http.StatusOK).Result(result).Write(w, r)
}
