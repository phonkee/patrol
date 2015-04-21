package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"text/template"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPaging(t *testing.T) {
	pp := NewPagingParams("limit", "page")
	Convey("test min/max/default limit", t, func() {
		data := []struct {
			min      int
			max      int
			def      int
			value    int
			expected int
		}{
			{10, 20, 15, 30, 15},
			{10, 20, 18, -1, 18},
			{10, 20, 16, 12, 12},
		}

		for _, item := range data {
			p := NewPaging(item.min, item.max, item.def, pp)
			p.SetLimit(item.value)
			So(p.Limit, ShouldEqual, item.expected)
		}
	})

	Convey("test set page/count, get offset", t, func() {
		data := []struct {
			limit int
			count int
			page  int
		}{
			{5, 10, 20},
			{5, 10, 20},
		}
		for _, item := range data {
			paging := NewPaging(item.limit, item.limit, item.limit, pp)
			paging.SetCount(item.count)
			So(paging.Count, ShouldEqual, item.count)
			paging.SetPage(item.page)
			So(paging.Page, ShouldEqual, item.page)
			So(paging.Offset(), ShouldEqual, paging.Page*paging.Limit)
		}
	})

	Convey("test read request", t, func() {
		tpldata := struct {
			LimitParam string
			PageParam  string
		}{
			"limit",
			"page",
		}

		data := []struct {
			url      string
			minlimit int
			maxlimit int
			limit    int
			page     int
		}{
			{"/?{{ .LimitParam }}=1&{{ .PageParam }}=20", 10, 20, 10, 20},
			{"/?{{ .LimitParam }}=60&{{ .PageParam }}=-1", 10, 20, 10, 0},
			{"/?{{ .LimitParam }}=15&{{ .PageParam }}=20", 10, 20, 15, 20},
		}

		for _, item := range data {
			tmpl, err := template.New("test").Parse(item.url)
			So(err, ShouldBeNil)
			buf := bytes.NewBufferString("")
			errexec := tmpl.Execute(buf, tpldata)
			So(errexec, ShouldBeNil)
			So(err, ShouldBeNil)
			So(err, ShouldBeNil)

			r, err := http.NewRequest("GET", buf.String(), nil)
			So(err, ShouldBeNil)
			p := NewPaging(item.minlimit, item.maxlimit, item.minlimit, NewPagingParams(tpldata.LimitParam, tpldata.PageParam))
			p.ReadRequest(r)
			So(p.Limit, ShouldEqual, item.limit)
			So(p.Page, ShouldEqual, item.page)
		}
	})

	Convey("test api paging", t, func() {
		page := int(13)
		limit := int(6)
		p := NewPaging(10, 100, 10, pp)
		p.SetPage(page)
		p.SetLimit(limit)
		values := url.Values{}

		vqp := NewQueryParams(p.UpdateURLValues(values))

		So(vqp.GetInt(p.Params.LimitParam), ShouldEqual, p.Limit)
		So(vqp.GetInt(p.Params.PageParam), ShouldEqual, p.Page)

		dp := NewPaging(-1, -1, -1, pp)
		values2 := url.Values{}
		valuesResultDisabled := NewQueryParams(dp.UpdateURLValues(values2))

		So(valuesResultDisabled.GetInt(dp.Params.LimitParam, -1), ShouldEqual, -1)
		So(valuesResultDisabled.GetInt(dp.Params.PageParam, -1), ShouldEqual, -1)

	})

	Convey("test LimitOffset", t, func() {
		dp := NewPaging(-1, -1, -1, pp)
		So(dp.LimitOffset(), ShouldEqual, "")

		ap := NewPaging(10, 100, 10, pp)
		lo := ap.LimitOffset()
		So(lo, ShouldEqual, fmt.Sprintf("LIMIT %d", ap.Limit))

		ap.SetPage(3)
		lo2 := ap.LimitOffset()
		So(lo2, ShouldEqual, fmt.Sprintf("LIMIT %d OFFSET %d", ap.Limit, ap.Offset()))
	})

	Convey("test UpdateBuilder", t, func() {
		ap := NewPaging(10, 100, 10, pp)
		ap.SetPage(10)
		qb := QueryBuilder().Select("*").From("table")
		qb = ap.UpdateBuilder(qb)

		query, _, _ := qb.ToSql()

		So(query, ShouldContainSubstring, fmt.Sprintf("LIMIT %d", ap.Limit))
		So(query, ShouldContainSubstring, fmt.Sprintf("OFFSET %d", ap.Offset()))

		dp := NewPaging(-1, -1, -1, pp)
		dqb := QueryBuilder().Select("*").From("table")
		dqb = dp.UpdateBuilder(dqb)
		queryd, _, _ := dqb.ToSql()

		So(queryd, ShouldNotContainSubstring, "LIMIT")
		So(queryd, ShouldNotContainSubstring, "OFFSET")
	})

}
