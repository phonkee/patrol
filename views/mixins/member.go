package mixins

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/views"
)

/*
Project detail mixin adds method to get project from mux vars
*/
type ProjectMemberTypeMixin struct{}

/*
returns member type for actual project
*/
func (b *ProjectMemberTypeMixin) MemberType(context *context.Context, r *http.Request, vars ...string) (mt models.MemberType, err error) {
	pm := models.NewProjectManager(context)
	project := pm.NewProject()

	// not found
	if err = pm.GetFromRequest(project, r, vars...); err != nil {
		return
	}

	um := models.NewUserManager(context)
	user := um.NewUser()
	if err = um.GetAuthUser(user, r); err != nil {
		return
	}

	teamm := models.NewTeamMemberManager(context)

	// get member type
	mt, err = teamm.MemberTypeByProject(project, user)

	return
}

/*
	Returns member type
*/
func (p *ProjectMemberTypeMixin) GetMemberType(project *models.Project, user *models.User, w http.ResponseWriter, r *http.Request) (mt models.MemberType, err error) {
	var ctx *context.Context
	if ctx, err = context.Get(r); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return mt, views.ErrInternalServerError
	}
	tmm := models.NewTeamMemberManager(ctx)
	if mt, err = tmm.MemberTypeByProject(project, user); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return mt, views.ErrForbidden
	}

	return
}
