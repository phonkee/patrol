package mixins

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
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
