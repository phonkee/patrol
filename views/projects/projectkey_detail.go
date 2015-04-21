package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
)

type ProjectKeyDetailAPIView struct {
	core.JSONView

	context *context.Context
}

func (p *ProjectKeyDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.Context(r)
	return
}

func (p *ProjectKeyDetailAPIView) GET(w http.ResponseWriter, r *http.Request) {
	pm := models.NewProjectKeyManager(p.context)

	_ = pm

	response.New(http.StatusOK).Write(w, r)
}
