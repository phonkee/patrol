package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/response"

	"github.com/gorilla/mux"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/views/mixins"
)

type ProjectDetailEventGroupListAPIView struct {
	core.JSONView

	// returns member type
	mixins.ProjectMemberTypeMixin

	// context
	context *context.Context
}

// check if user is member of project
func (p *ProjectDetailEventGroupListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.Context(r)
	if _, err = p.MemberType(p.context, r); err != nil {
		response.New().Status(http.StatusUnauthorized).Write(w, r)
		return
	}

	// @TODO: get project and store to view

	return
}

/*
Retrieve list of event groups for given project
*/
func (p *ProjectDetailEventGroupListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error
	response := response.New()
	egm := models.NewEventGroupManager(p.context)
	egl := egm.NewEventGroupList()
	vars := mux.Vars(r)

	pm := models.NewProjectManager(p.context)
	project := pm.NewProject()

	// something bad happened
	if err = pm.GetFromRequest(project, r); err != nil {
		response.Status(http.StatusInternalServerError).Write(w, r)
		return
	}

	// filter event groups for given project
	// @TODO: add query param filtering
	if err = egm.Filter(&egl, egm.QueryFilterWhere("project_id = ?", vars["project_id"])); err != nil {
		response.Status(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.Status(http.StatusOK).Result(egl).Write(w, r)
}
