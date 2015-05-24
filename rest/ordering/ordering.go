/*
Ordering object
Usable for ordering sql queries. Supports allowed fields.
Method QueryFunc returns function directly usable as QueryQueryFunc in DBFilter.
*/
package ordering

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lann/squirrel"
	"github.com/phonkee/patrol/utils"
)

type orderdir int

const (
	ORDER_ASC = iota
	ORDER_DESC
)

func New(paramName string, allowed ...string) *Ordering {
	return &Ordering{
		Allowed:   allowed,
		paramName: paramName,
	}
}

type Ordering struct {
	Allowed       []string `json:"allowed"`
	OrderingOrder string   `json:"order"`
	paramName     string
}

func (o *Ordering) Allow(allowed ...string) *Ordering {
	if len(allowed) > 0 {
		o.Allowed = append(o.Allowed, allowed...)
	}
	return o
}

func (o *Ordering) Order(orderstring string) (result *Ordering) {
	var (
		field string
		order string
	)

	order = "ASC"
	orderstring = strings.TrimSpace(orderstring)

	result = o

	if orderstring == "" {
		o.OrderingOrder = ""
		return
	}

	if fmt.Sprintf("%c", orderstring[0]) == "-" {
		field, order = fmt.Sprint(orderstring[1:]), "DESC"
	} else {
		field = orderstring
	}

	o.OrderingOrder = ""
	for _, af := range o.Allowed {
		if af == field {
			o.OrderingOrder = fmt.Sprintf("%s %s", field, order)
			return
		}
	}

	return
}

func (o *Ordering) QueryFunc() utils.QueryFunc {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		if o.OrderingOrder != "" {
			builder = builder.OrderBy(o.OrderingOrder)
		}
		return builder
	}
}

func (o *Ordering) ReadRequest(r *http.Request) *Ordering {
	return o.Order(r.URL.Query().Get(o.paramName))
}
