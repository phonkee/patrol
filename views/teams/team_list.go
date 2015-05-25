package teams

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/views/mixins"
)

/*
TeamListAPIView

Team list endpoint
Provides following methods
	Retrieve - GET - returns list of teams
	Create - POST - creates team
*/
type TeamListAPIView struct {
	views.APIView

	// mixins
	mixins.AuthUserMixin

	context *context.Context

	user *models.User
}

func (t *TeamListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	t.context = t.GetContext(r)

	t.user = models.NewUser()
	if err = t.GetAuthUser(t.user, w, r); err != nil {
		return
	}

	// check
	switch r.Method {
	case "POST":
		if !t.user.IsSuperuser {
			response.New(http.StatusForbidden).Write(w, r)
			return views.ErrBreakRequest
		}
	}

	return
}

/*
Returns options for current user
*/
func (t *TeamListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {

	md := metadata.New("Teams")

	if t.user.IsSuperuser {
		createAction := md.ActionCreate()
		createAction.Field("name").SetHelpText("team name").SetRequired(true).SetMax(200).SetMin(5)

		// add delete action
		md.ActionDelete()
	}

	retrieve := md.ActionRetrieve()
	retrieve.Field("name").SetHelpText("team name")

	// write metadata to response
	response.New(http.StatusOK).Metadata(md).Write(w, r)
}

// Retrieve tems list
func (t *TeamListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error

	teamm := models.NewTeamManager(t.context)
	teams := teamm.NewTeamList()

	// Filter teams for user

	if err = teamm.FilterByUser(&teams, t.user); err != nil {
		response.New(http.StatusUnauthorized).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(teams).Write(w, r)
}

// Create team
func (t *TeamListAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	serializer := serializers.TeamsTeamCreateSerializer{}
	if err = t.context.Bind(&serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	if vr := serializer.Validate(t.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	var team *models.Team

	// save serializer (team)
	if team, err = serializer.Save(t.context, t.user); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.New(http.StatusCreated).Result(team).Write(w, r)
}
