package mixins

import (
	"errors"
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/types"
)

var (
	ErrInvalidParam = errors.New("invalid_url_param")
	ErrNotFound     = errors.New("not_found")
)

/*
EventGroupDetailMixin is mixin for event group details. Fetches data from db and
store them
*/
type EventGroupMixin struct{}

/*
Loads eventgroup from storage
*/
func (e *EventGroupMixin) GetEventGroup(target interface{}, w http.ResponseWriter, r *http.Request, muxvar ...string) (err error) {
	var ctx *context.Context
	// get context
	if ctx, err = context.Get(r); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return views.ErrInternalServerError
	}

	varname := "eventgroup_id"
	if len(muxvar) > 0 {
		varname = muxvar[0]
	}

	// read eventgroup id from mux vars
	var pk types.PrimaryKey

	if pk, err = rest.GetMuxVarPrimaryKey(r, varname); err != nil {
		err = views.ErrInvalidParam
		response.New(http.StatusBadRequest).Error(err).Write(w, r)
		return
	}

	manager := models.NewEventGroupManager(ctx)

	if err = manager.GetByID(target, pk); err != nil {
		switch err {
		case models.ErrObjectDoesNotExists:
			response.New(http.StatusNotFound).Write(w, r)
		default:
			response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		}
		return
	}

	return
}
