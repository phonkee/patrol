package projects

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/views/mixins"
)

/*
Mark eventgroup as resolved
send notification for frontend
*/
type ProjectDetailEventGroupResolveAPIView struct {
	core.JSONView
	context *context.Context
	user    *models.User

	// returns member type
	mixins.ProjectMemberTypeMixin
	mixins.EvenGroupDetailMixin
}

/*
Before method retrieves eventgroup, project from datastore.
*/
func (p *ProjectDetailEventGroupResolveAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.Context(r)

	p.user = models.NewUser()
	if err = p.user.Manager(p.context).GetAuthUser(p.user, r); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

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

/*
	Marks eventgroup as resolved
*/
func (p *ProjectDetailEventGroupResolveAPIView) POST(w http.ResponseWriter, r *http.Request) {
	egm := models.NewEventGroupManager(p.context)
	if err := egm.Resolve(p.EventGroup, p.user); err != nil {
		response.New(http.StatusNotAcceptable).Write(w, r)
		return
	}

	response.New(http.StatusOK).Write(w, r)
	return
}
