package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/phonkee/patrol"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	if err := patrol.Setup(); err != nil {
		if err != patrol.ErrPatrolAlreadySetup {
			fmt.Printf("patrol: setup failed with error: %s", err)
		}
	}

	patrol.Run([]string{"migrate"})

	Convey("Test Version", t, func() {
		r := NewAPIRequest("GET", settings.ROUTE_COMMON_VERSION, nil)

		result := struct {
			Result struct {
				Version string `json:"version"`
			} `json:"result"`
		}{}
		w, err := GetAPIResponse(r, &result)

		So(err, ShouldBeNil)
		So(w.Code, ShouldEqual, http.StatusOK)
		So(result.Result.Version, ShouldEqual, settings.VERSION)
	})
}
