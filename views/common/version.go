package common

import (
	"net/http"

	"github.com/phonkee/patrol/parser"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/settings"
)

// Returns version of patrol
type VersionAPIView struct {
	views.APIView
}

func (v *VersionAPIView) GET(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{
		"version":   settings.VERSION,
		"protocols": parser.Versions(),
	}

	response.New(http.StatusOK).Result(result).Write(w, r)
}
