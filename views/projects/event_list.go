package projects

import (
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
)

type ProjectEventGroupEventListView struct {
	core.JSONView

	project    *models.Project
	eventgroup *models.EventGroup
}

func (p *ProjectEventGroupEventListView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)

	var (
		projectid    int64
		eventgroupid int64
	)

	if projectid, err = strconv.ParseInt(vars["project_id"], 10, 0); err != nil {
		return
	}

	if eventgroupid, err = strconv.ParseInt(vars["eventgroup_id"], 10, 0); err != nil {
		return
	}

	glog.Infof("this is project id %+v and this is eventgroup id %+v", projectid, eventgroupid)

	return
}

func (p *ProjectEventGroupEventListView) GET(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusOK).Write(w, r)

	return
}
