package mixins

import (
	"errors"
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest"
	"github.com/phonkee/patrol/rest/response"
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
type EvenGroupDetailMixin struct {
	EventGroup *models.EventGroup
	Project    *models.Project
}

/*
Loads objects from request
*/
func (e *EvenGroupDetailMixin) GetInstances(w http.ResponseWriter, r *http.Request) (err error) {

	// get context
	ctx, _ := context.Get(r)

	egm := models.NewEventGroupManager(ctx)
	pm := models.NewProjectManager(ctx)

	// read project id
	var projectid types.PrimaryKey
	if projectid, err = rest.GetMuxVarPrimaryKey(r, "project_id"); err != nil {
		response.New(http.StatusBadRequest).Error(ErrInvalidParam).Write(w, r)
		return ErrInvalidParam
	}

	// read eventgroup id from mux vars
	var eventgroupid types.PrimaryKey
	if eventgroupid, err = rest.GetMuxVarPrimaryKey(r, "eventgroup_id"); err != nil {
		response.New(http.StatusBadRequest).Error(ErrInvalidParam).Write(w, r)
		return ErrInvalidParam
	}

	// get project from database
	e.Project = models.NewProject()
	if err = pm.GetByID(e.Project, projectid); err != nil {
		response.New(http.StatusNotFound).Write(w, r)
		return ErrNotFound
	}

	// get eventgroup from database
	e.EventGroup = models.NewEventGroup()
	if err = egm.GetByID(e.EventGroup, eventgroupid); err != nil {
		response.New(http.StatusNotFound).Write(w, r)
		return ErrNotFound
	}

	// check if eventgroup belongs to project
	if e.EventGroup.ProjectID.ToPrimaryKey() != e.Project.ID {
		response.New(http.StatusNotFound).Write(w, r)
		return ErrNotFound
	}

	return
}
