package events

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/views/mixins"
)

func NewEventListView() views.Viewer {
	return &EventListView{
		eventgroup: models.NewEventGroup(),
		project:    models.NewProject(),
	}
}

/*
EventListView
	list of events to eventgroup
*/
type EventListView struct {
	views.APIView
	mixins.EventGroupMixin
	mixins.ProjectMemberTypeMixin
	mixins.ProjectsProjectMixin

	context *context.Context

	eventgroup *models.EventGroup
	project    *models.Project
}

func (p *EventListView) Before(w http.ResponseWriter, r *http.Request) (err error) {
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

func (p *EventListView) GET(w http.ResponseWriter, r *http.Request) {
	glog.Info("this is eventlist %+v\n", p.eventgroup)

	response.New(http.StatusOK).Write(w, r)
	return
}
