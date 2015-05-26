package events

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/views/mixins"
)

/*
Mark eventgroup as resolved
send notification for frontend
*/
type EventGroupResolveAPIView struct {
	views.APIView
	context *context.Context
	user    *models.User

	// returns member type
	mixins.AuthUserMixin
	mixins.ProjectMemberTypeMixin
	mixins.ProjectsProjectMixin
	mixins.EventGroupMixin

	eventgroup *models.EventGroup
	project    *models.Project
}

/*
Before method retrieves eventgroup, project from datastore.
*/
func (p *EventGroupResolveAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)

	p.user = models.NewUser()
	if err = p.GetAuthUser(p.user, w, r); err != nil {
		return
	}

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
	Marks eventgroup as resolved
*/
func (p *EventGroupResolveAPIView) POST(w http.ResponseWriter, r *http.Request) {
	egm := models.NewEventGroupManager(p.context)
	if err := egm.Resolve(p.eventgroup, p.user); err != nil {
		response.New(http.StatusNotAcceptable).Write(w, r)
		return
	}

	response.New(http.StatusOK).Write(w, r)
	return
}
