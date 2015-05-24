package query_params

import (
	"net/url"
	"strconv"
	"strings"
)

var (
	boolvals = map[string]bool{
		"f": false, "false": false, "0": false, "off": false,
		"t": true, "true": true, "1": true, "on": true,
	}
)

func New(values url.Values) *QueryParams {
	return &QueryParams{
		values,
	}
}

type QueryParams struct{ url.Values }

// returns int64 value
func (q *QueryParams) GetInt(name string, def ...int) (result int) {
	var err error

	if result, err = strconv.Atoi(q.GetString(name)); err == nil {
		return
	}

	if len(def) > 0 {
		result = def[0]
	}

	return
}

// returns bool value
func (q *QueryParams) GetBool(name string, def ...bool) (result bool) {
	var ok bool

	if result, ok = boolvals[strings.ToLower(q.GetString(name))]; ok {
		return
	}

	if len(def) > 0 {
		result = def[0]
	}

	return
}

// returns string value
func (q *QueryParams) GetString(name string, def ...string) (result string) {
	result = strings.TrimSpace(q.Get(name))
	if result == "" && len(def) > 0 {
		result = def[0]
	}
	return
}

// returns float64 value
func (q *QueryParams) GetFloat(name string, def ...float64) (result float64) {
	var err error
	if result, err = strconv.ParseFloat(q.GetString(name), 0); err == nil {
		return
	}

	if len(def) > 0 {
		result = def[0]
	}

	return
}
