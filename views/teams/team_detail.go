package teams

import (
	"net/http"

	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/response"
)

type TeamDetailAPIView struct {
	core.JSONView

	// store team instance in Before method
	team *models.Team
}

func (t *TeamDetailAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	context := t.Context(r)

	um := models.NewUserManager(context)
	user := um.NewUser()
	if err = um.GetAuthUser(user, r); err != nil {
		response.New(http.StatusUnauthorized).Write(w, r)
		return core.ErrUnauthorized
	}

	// get team and store it to context for later usage
	tm := models.NewTeamManager(context)
	t.team = tm.NewTeam()
	if err = tm.GetFromRequest(t.team, r, "team_id"); err != nil {
		return
	}
	return
}

/*
Retrieve team
*/
func (t *TeamDetailAPIView) GET(w http.ResponseWriter, r *http.Request) {
	response.New(http.StatusOK).Result(t.team).Write(w, r)
}
