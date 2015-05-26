package events

import (
	"errors"
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"

	"github.com/phonkee/patrol/views/mixins"
)

var (
	ErrInvalidParam = errors.New("invalid_url_param")
)

type EventGroupDetailAPIView struct {
	views.APIView

	// returns member type
	mixins.ProjectMemberTypeMixin
	mixins.ProjectsProjectMixin
	mixins.EventGroupMixin

	// context
	context *context.Context

	eventgroup *models.EventGroup
	project    *models.Project
}

/*
	Before method fetches project and eventgroup from storage.
	CHeck if eventgroup project is same as project
	Then check whether user has permissions to view this eventgroup
*/
func (p *EventGroupDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)

	if err = p.GetProject(p.project, w, r); err != nil {
		return
	}

	if err = p.GetEventGroup(p.eventgroup, w, r); err != nil {
		return
	}

	// check
	if p.eventgroup.ProjectID.ToPrimaryKey() != p.project.ID {
		response.New(http.StatusNotFound).Write(w, r)
		return views.ErrNotFound
	}

	// check membership in project
	if _, err = p.MemberType(p.context, r); err != nil {
		response.New().Status(http.StatusUnauthorized).Write(w, r)
		return
	}

	return
}

/*
	Retrieve method
	Since all loading is performed in Before method this is really simple -
	just make response
*/
func (p *EventGroupDetailAPIView) GET(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusOK).Result(p.eventgroup).Write(w, r)
}
