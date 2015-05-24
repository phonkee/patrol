package common

import (
	"net/http"
	"testing"

	"github.com/phonkee/patrol/apitest"
	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {

	apitest.Setup()

	Convey("Test Version", t, func() {

		session := apitest.NewSession()
		request := session.Request("GET", settings.ROUTE_COMMON_VERSION).Do()

		result := struct {
			Result struct {
				Version string `json:"version"`
			} `json:"result"`
		}{}

		request.Scan(&result)
		So(request.Error(), ShouldBeNil)
		So(request.Response().Code, ShouldEqual, http.StatusOK)
		So(result.Result.Version, ShouldEqual, settings.VERSION)
	})
}
