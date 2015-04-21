package teams

import (
	"net/http"

	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/validator"
)

type TeamCreateSerializer struct {
	Name string `json:"name"`
}

/*
TeamListAPIView

Team list endpoint
Provides following methods
	Retrieve - GET - returns list of teams
	Create - POST - creates team
*/
type TeamListAPIView struct {
	core.JSONView

	user *models.User
}

func (t *TeamListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	context := t.Context(r)
	t.user = models.NewUser()

	// user somehow cannot be read from request
	if err = t.user.Manager(context).GetAuthUser(t.user, r); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return core.ErrUnauthorized
	}

	// check
	switch r.Method {
	case "POST":
		if !t.user.IsSuperuser {
			response.New(http.StatusForbidden).Write(w, r)
			return core.ErrBreakRequest
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
		createAction.Field("name").SetHelpText("team name").SetRequired(true).SetMax(135)

		// add delete action
		md.ActionDelete()
	}

	retrieve := md.ActionRetrieve()
	retrieve.Field("name").SetHelpText("team name")

	response.New(http.StatusOK).Raw(md).Write(w, r)
	return
}

// Retrieve tems list
func (t *TeamListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error

	context := t.Context(r)

	authm := models.NewUserManager(context)
	teamm := models.NewTeamManager(context)

	// GetAuthUser - returns user from request
	user := authm.NewUser()
	if err = authm.GetAuthUser(user, r); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	teams := teamm.NewTeamList()

	// Filter teams for user

	if err = teamm.FilterByUser(&teams, user); err != nil {
		response.New(http.StatusUnauthorized).Error(err).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(teams).Write(w, r)
}

// Create team
func (t *TeamListAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	context := t.Context(r)

	cs := TeamCreateSerializer{}
	if err = context.Bind(&cs); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	// GetAuthUser - returns user from request
	authm := models.NewUserManager(context)
	user := authm.NewUser()
	if err = authm.GetAuthUser(user, r); err != nil {
		response.New(http.StatusForbidden).Write(w, r)
		return
	}

	teamm := models.NewTeamManager(context)
	team := teamm.NewTeam(func(team *models.Team) {
		team.Name = cs.Name
		team.OwnerID = user.PrimaryKey().ToForeignKey()
	})

	var result *validator.Result
	if result, err = team.Validate(context); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	if !result.IsValid() {
		response.New(http.StatusBadRequest).Error(result).Write(w, r)
		return
	}

	if err = team.Insert(context); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.New(http.StatusCreated).Result(team).Write(w, r)
}
