package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/metadata"
	. "github.com/smartystreets/goconvey/convey"
)

func GetTestContext() *context.Context {
	ctx, err := context.NewTest()
	if err != nil {
		panic("oops")
	}
	req, _ := http.NewRequest("OPTIONS", "/api/auth/login", nil)
	ctx = ctx.WithRequest(req)
	return ctx
}

func TestNewResponse(t *testing.T) {

	Convey("Test status", t, func() {

		ctx := GetTestContext()

		s := http.StatusConflict
		So(New().Status(s).Status_, ShouldEqual, s)
		So(New(s).Status_, ShouldEqual, s)

		rec := httptest.NewRecorder()

		New().Status(http.StatusTeapot).Result("hello world!").Write(rec, ctx.Request)
		So(rec.Code, ShouldEqual, http.StatusTeapot)

		// fmt.Printf("this is it %s", New().Status(http.StatusTeapot).Result([]struct{}{}).ResultSize(0))
		// fmt.Printf("this is it %s", New().Status(http.StatusTeapot).Error(errors.New("exception")))
	})

	Convey("TestOther", t, func() {
		ctx := GetTestContext()

		s := http.StatusGone

		type Product struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		md := metadata.New("hello")
		md.ActionCreate().FromStruct(Product{})

		rec := httptest.NewRecorder()
		New(s).Metadata(md).Write(rec, ctx.Request)

		fmt.Printf("this is recorder %+v", rec)

	})

}
