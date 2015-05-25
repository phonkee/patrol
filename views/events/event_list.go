package events

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/views/mixins"
)

/*
EventListView
	list of events to eventgroup
*/
type EventListView struct {
	views.APIView
	mixins.EvenGroupDetailMixin
	mixins.ProjectMemberTypeMixin

	context *context.Context
}

func (p *EventListView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)
	if err = p.GetInstances(w, r); err != nil {
		return
	}

	// check membership in project
	if _, err = p.MemberType(p.context, r); err != nil {
		response.New().Status(http.StatusUnauthorized).Write(w, r)
		return
	}

	return
}

func (p *EventListView) GET(w http.ResponseWriter, r *http.Request) {
	glog.Info("this is eventgroup %+v\n", p.EventGroup)

	response.New(http.StatusOK).Write(w, r)
	return
}
