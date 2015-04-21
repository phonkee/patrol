package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lann/squirrel"
)

// Returns new paging instance
func NewPaging(minlimit, maxlimit, def int, params *PagingParams) *Paging {
	paging := &Paging{
		MinLimit:     minlimit,
		MaxLimit:     maxlimit,
		DefaultLimit: def,
		Params:       params,
	}
	paging.SetLimit(def)
	return paging
}

// Paging implementation
type Paging struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Count int `json:"count"`

	// limits
	MinLimit     int `json:"-"`
	MaxLimit     int `json:"-"`
	DefaultLimit int `json:"-"`

	// paging params
	Params *PagingParams `json:"-"`
}

func (p *Paging) ReadRequest(r *http.Request) *Paging {
	qp := NewQueryParams(r.URL.Query())
	p.SetLimit(qp.GetInt(p.Params.LimitParam, -1))
	p.SetPage(qp.GetInt(p.Params.PageParam, -1))

	// p.SetLimit(p.Params.GetLimit(r.URL.Query()))
	// p.SetPage(p.Params.GetPage(r.URL.Query()))
	return p
}

func (p *Paging) SetLimit(limit int) *Paging {
	if limit > p.MaxLimit || limit < p.MinLimit {
		p.Limit = p.DefaultLimit
		return p
	}
	p.Limit = limit
	return p
}

func (p *Paging) SetPage(page int) *Paging {
	if page < 0 {
		p.Page = 0
	} else {
		p.Page = page
	}
	return p
}

func (p *Paging) SetCount(count int) *Paging {
	p.Count = count
	return p
}

func (p *Paging) Offset() int {
	return p.Limit * p.Page
}

// returns limit/offset clause sql query
func (p *Paging) LimitOffset() (result string) {
	if p.Limit <= 0 {
		result = ""
	} else {
		buffer := bytes.NewBufferString("")
		buffer.WriteString(fmt.Sprintf("LIMIT %d", p.Limit))
		if offset := p.Offset(); offset != 0 {
			buffer.WriteString(fmt.Sprintf(" OFFSET %d", offset))
		}
		result = buffer.String()
	}
	return
}

func (p *Paging) UpdateBuilder(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if p.Limit <= 0 {
		return builder
	}
	return builder.Limit(uint64(p.Limit)).Offset(uint64(p.Offset()))
}

// updates url values with paging values
func (p *Paging) UpdateURLValues(values url.Values) url.Values {
	if p.Limit <= 0 {
		return values
	}
	values.Set(p.Params.LimitParam, fmt.Sprintf("%d", p.Limit))
	if p.Page > 0 {
		values.Set(p.Params.PageParam, fmt.Sprintf("%d", p.Page))
	}
	return values
}

/*
	PagingParams
*/
type PagingParams struct {
	LimitParam string
	PageParam  string
}

func NewPagingParams(LimitParam, PageParam string) *PagingParams {
	return &PagingParams{
		LimitParam: LimitParam,
		PageParam:  PageParam,
	}
}
