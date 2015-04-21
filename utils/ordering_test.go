package utils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/phonkee/patrol/settings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOrdering(t *testing.T) {
	Convey("test ordering", t, func() {
		o := NewOrdering("order", "id", "username")
		o.Order("id")
		So(o.OrderingOrder, ShouldEqual, "id ASC")
		So(len(o.Allowed), ShouldEqual, 2)

		o.Allow()
		So(len(o.Allowed), ShouldEqual, 2)

		o.Allow("another")
		So(len(o.Allowed), ShouldEqual, 3)

		o.Order("")
		So(o.OrderingOrder, ShouldEqual, "")

		o.Order("-username")
		So(o.OrderingOrder, ShouldEqual, "username DESC")

		o.Order("-invalidfield")
		So(o.OrderingOrder, ShouldEqual, "")

		url := fmt.Sprintf("?%s=-id", settings.ORDERING_DEFAULT_PARAM_NAME)
		r, _ := http.NewRequest(settings.HTTP_GET, url, nil)

		o.ReadRequest(r)
		So(o.OrderingOrder, ShouldEqual, "id DESC")

		qb := QueryBuilder().Select("*").From("table")
		qb = ApplyQueryFuncs(qb, o.QueryFunc())

		q, _, _ := qb.ToSql()
		So(q, ShouldContainSubstring, o.OrderingOrder)

	})
}
