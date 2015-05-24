package mixins

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/views"
)

type ProjectsProjectMixin struct{}

/*
	Returns project by request

	target object must be pointer
*/
func (p *ProjectsProjectMixin) GetProject(target interface{}, w http.ResponseWriter, r *http.Request, varname ...string) (err error) {
	var ctx *context.Context
	if ctx, err = context.Get(r); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return views.ErrInternalServerError
	}

	mux_param := "project_id"
	if len(varname) > 0 {
		mux_param = varname[0]
	}

	var pk types.PrimaryKey
	if pk, err = rest.GetMuxVarPrimaryKey(r, mux_param); err != nil {
		response.New(http.StatusBadRequest).Error(views.ErrInvalidParam).Write(w, r)
		return views.ErrInvalidParam
	}

	pm := models.NewProjectManager(ctx)

	if err = pm.GetByID(target, pk); err != nil {
		response.New(http.StatusNotFound).Error(views.ErrNotFound).Write(w, r)
		return views.ErrNotFound
	}

	return
}
