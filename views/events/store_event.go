package events

import (
	"net/http"
	"strconv"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/response"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/parser"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
)

type StoreEventAPIView struct {
	core.JSONView

	context *context.Context

	project *models.Project
}

func (s *StoreEventAPIView) GetProjectID(r *http.Request) (id int64, err error) {
	vars := mux.Vars(r)
	id, err = strconv.ParseInt(vars["project_id"], 10, 0)
	return
}

func (s *StoreEventAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	s.context = s.Context(r)
	var pid int64
	if pid, err = s.GetProjectID(r); err != nil {
		return
	}

	response := response.New()
	projectman := models.NewProjectManager(s.context)

	reqmanager := models.NewRequestManager(s.context)
	var auth map[string]string
	if auth, err = reqmanager.SentryAuthHeaders(r); err != nil {
		return
	}

	s.project = projectman.NewProject()
	if err = projectman.GetByAuth(s.project, auth[settings.SENTRY_AUTH_KEY], auth[settings.SENTRY_AUTH_SECRET]); err != nil {
		glog.Error(err)
		return
	}

	/*
		if project.PrimaryKey().Int64() != pid {
			response.Status(http.StatusNotFound).Write(w, r)
			return core.ErrBreakRequest
		}
	*/
	_ = response
	_ = pid

	return nil
}

func (s *StoreEventAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error
	response := response.New()

	reqmanager := models.NewRequestManager(s.context)

	var (
		auth = map[string]string{}
	)

	if auth, err = reqmanager.SentryAuthHeaders(r); err != nil {
		glog.Errorf("%v", err)
		return
	}

	version := auth[settings.SENTRY_AUTH_VERSION]

	var (
		events []*parser.RawEvent
	)
	events, err = parser.ParseRequest(r, version)
	if err != nil {
		glog.Errorf("this is parserequest error %v", err)
	}

	// no events returned and return is nil ? wohooo
	if len(events) == 0 {
		response.Status(http.StatusBadRequest).Write(w, r)
		return
	}

	// add project id to events
	for _, event := range events {
		event.ProjectID = types.ForeignKey(s.project.ID)
	}

	// Here we should send all events to queue
	result := map[string]string{
		"event_id": events[0].EventID,
	}

	raweventmanager := parser.NewRawEventManager(s.context)

	// push message to queue
	for _, event := range events {
		raweventmanager.PushRawEvent(event)
	}

	response.Status(http.StatusOK).Raw(result).Write(w, r)
}
