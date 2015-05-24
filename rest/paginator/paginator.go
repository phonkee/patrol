package paginator

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lann/squirrel"
	"github.com/phonkee/patrol/rest/query_params"
)

// Returns new Paginator instance
func New(minlimit, maxlimit, def int, params *PaginatorParams) *Paginator {
	Paginator := &Paginator{
		MinLimit:     minlimit,
		MaxLimit:     maxlimit,
		DefaultLimit: def,
		Params:       params,
	}
	Paginator.SetLimit(def)
	return Paginator
}

// Paginator implementation
type Paginator struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Count int `json:"count"`

	// limits
	MinLimit     int `json:"-"`
	MaxLimit     int `json:"-"`
	DefaultLimit int `json:"-"`

	// Paginator params
	Params *PaginatorParams `json:"-"`
}

func (p *Paginator) ReadRequest(r *http.Request) *Paginator {
	qp := query_params.New(r.URL.Query())
	p.SetLimit(qp.GetInt(p.Params.LimitParam, -1))
	p.SetPage(qp.GetInt(p.Params.PageParam, -1))

	// p.SetLimit(p.Params.GetLimit(r.URL.Query()))
	// p.SetPage(p.Params.GetPage(r.URL.Query()))
	return p
}

func (p *Paginator) SetLimit(limit int) *Paginator {
	if limit > p.MaxLimit || limit < p.MinLimit {
		p.Limit = p.DefaultLimit
		return p
	}
	p.Limit = limit
	return p
}

func (p *Paginator) SetPage(page int) *Paginator {
	if page < 0 {
		p.Page = 0
	} else {
		p.Page = page
	}
	return p
}

func (p *Paginator) SetCount(count int) *Paginator {
	p.Count = count
	return p
}

func (p *Paginator) Offset() int {
	return p.Limit * p.Page
}

// returns limit/offset clause sql query
func (p *Paginator) LimitOffset() (result string) {
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

func (p *Paginator) UpdateBuilder(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if p.Limit <= 0 {
		return builder
	}
	return builder.Limit(uint64(p.Limit)).Offset(uint64(p.Offset()))
}

// updates url values with Paginator values
func (p *Paginator) UpdateURLValues(values url.Values) url.Values {
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
	PaginatorParams
*/
type PaginatorParams struct {
	LimitParam string
	PageParam  string
}

func NewParams(LimitParam, PageParam string) *PaginatorParams {
	return &PaginatorParams{
		LimitParam: LimitParam,
		PageParam:  PageParam,
	}
}
