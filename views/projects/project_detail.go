package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/response"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/views/mixins"
)

/*
Project Detail view

	/api/projects/project/{project_id:[0-9]+}

*/
type ProjectDetailAPIView struct {
	core.JSONView

	// member type
	mixins.ProjectMemberTypeMixin

	context *context.Context
}

/*
Before function checks member type for auth user/project
*/
func (p *ProjectDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)
	if _, err = p.MemberType(p.context, r); err != nil {
		response.New().Status(http.StatusUnauthorized).Write(w, r)
		return err
	}
	return
}

/*
GET method handler
*/
func (p *ProjectDetailAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error

	// prepare blank response
	response := response.New()

	manager := models.NewProjectManager(p.context)

	project := manager.NewProject()
	err = manager.GetFromRequest(project, r)
	if err != nil {
		if err == models.ErrObjectDoesNotExists {
			response.Status(http.StatusNotFound).Write(w, r)
		} else {
			glog.Error(err)
			response.Status(http.StatusInternalServerError).Write(w, r)
		}
		return
	}

	response.Status(http.StatusOK).Result(project).Write(w, r)
}
