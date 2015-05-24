/*
Various http helpers

GetMuxVar...  - functions to get variable from mux router
*/
package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/phonkee/patrol/types"
)

var (
	ErrMuxVarNotFound = errors.New("mux var not found")
)

/*
	Returns mux variable as string
*/
func GetMuxVarString(r *http.Request, name string) (value string, err error) {
	vars := mux.Vars(r)
	stringvar, ok := vars[name]
	if !ok {
		return "", ErrMuxVarNotFound
	}
	return stringvar, nil
}

/*
returns mux variable as int64
*/
func GetMuxVarInt64(r *http.Request, name string) (value int64, err error) {
	var stringvar string
	if stringvar, err = GetMuxVarString(r, name); err != nil {
		return
	}

	if value, err = strconv.ParseInt(stringvar, 10, 0); err != nil {
		return
	}

	return
}

/*
	Returns mux var as primary key
*/
func GetMuxVarPrimaryKey(r *http.Request, name string) (value types.PrimaryKey, err error) {
	var intvar int64

	if intvar, err = GetMuxVarInt64(r, name); err != nil {
		return
	}
	return types.PrimaryKey(intvar), nil
}
