package projects

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
)

type ProjectKeyListAPIView struct {
	core.JSONView

	context *context.Context
}

func (p *ProjectKeyListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)
	return
}

// retrieve project keys
func (p *ProjectKeyListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	pm := models.NewProjectKeyManager(p.context)

	vars := mux.Vars(r)

	var (
		err error
		id  int64
	)

	if id, err = strconv.ParseInt(vars["project_id"], 10, 0); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	// create result list
	result := models.NewProjectKeyList()
	if err = pm.Filter(&result, pm.QueryFilterProjectID(id)); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(result).Write(w, r)
}

// OPTIONS - returns metadata
func (p *ProjectKeyListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	md := metadata.New("Project key list")
	md.ActionRetrieve().From(models.NewProjectKeyList())
	response.New(http.StatusOK).Metadata(md).Write(w, r)
}
