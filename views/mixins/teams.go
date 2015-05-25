package mixins

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
	"github.com/phonkee/patrol/views"
)

type TeamsTeamMixin struct{}

/*
	Returns team
	if not found response is written and error is returned
*/
func (t *TeamsTeamMixin) GetTeam(target interface{}, w http.ResponseWriter, r *http.Request, muxvar ...string) (err error) {
	varname := "team_id"
	if len(muxvar) > 0 {
		varname = muxvar[0]
	}

	var pk types.PrimaryKey

	if pk, err = rest.GetMuxVarPrimaryKey(r, varname); err != nil {
		err = views.ErrInvalidParam
		response.New(http.StatusBadRequest).Error(err).Write(w, r)
		return
	}

	var ctx *context.Context
	if ctx, err = context.Get(r); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	manager := models.NewTeamManager(ctx)
	if err = manager.GetByID(target, pk); err != nil {
		switch err {
		default:
			response.New(http.StatusInternalServerError).Write(w, r)
		case models.ErrObjectDoesNotExists:
			response.New(http.StatusNotFound).Write(w, r)
		}
		return
	}

	return
}

/*
Teams team member mixin
provides helper methods for tem members
*/
type TeamsTeamMemberMixin struct{}

func (t *TeamsTeamMemberMixin) Filter(target *[]*serializers.TeamsTeamMemberDetailSerializer, ctx *context.Context, qfs ...utils.QueryFunc) (err error) {
	manager := models.NewTeamMemberManager(ctx)
	if err = manager.Filter(target, qfs...); err != nil {
		return
	}

	for _, item := range *target {
		if _, err = item.LoadUser(ctx); err != nil {
			return
		}
	}

	return
}

/*
Returns single team member instance, if error occures it writes response and returns error
*/
func (t *TeamsTeamMemberMixin) GetTeamMember(target interface{}, w http.ResponseWriter, r *http.Request, muxvar ...string) (err error) {
	varname := "teammember_id"
	if len(muxvar) > 0 {
		varname = muxvar[0]
	}

	var pk types.PrimaryKey

	if pk, err = rest.GetMuxVarPrimaryKey(r, varname); err != nil {
		err = views.ErrInvalidParam
		response.New(http.StatusBadRequest).Error(err).Write(w, r)
		return
	}

	var ctx *context.Context
	if ctx, err = context.Get(r); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	manager := models.NewTeamMemberManager(ctx)
	if err = manager.GetByID(target, pk); err != nil {
		switch err {
		case models.ErrObjectDoesNotExists:
			response.New(http.StatusNotFound).Write(w, r)
			return views.ErrNotFound
		default:
			response.New(http.StatusInternalServerError).Write(w, r)
			return views.ErrInternalServerError
		}
	}

	return
}
