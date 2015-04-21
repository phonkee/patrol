package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/phonkee/patrol/context"

	"github.com/gorilla/mux"
	"github.com/lann/squirrel"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

// Team model
type Team struct {
	Model
	Name      string           `db:"name" json:"name" validator:"name"`
	OwnerID   types.ForeignKey `db:"owner_id" json:"owner_id"`
	DateAdded time.Time        `db:"date_added" json:"date_added"`
	Status    TeamStatus       `db:"status" json:"status"`
}

// returns all columns except of primary key
func (t *Team) Columns() []string {
	return []string{"name", "owner_id", "date_added", "status"}
}
func (t *Team) Values() []interface{} {
	return []interface{}{t.Name, t.OwnerID, t.DateAdded, t.Status}
}
func (t *Team) String() string { return "teams:team:" + t.PrimaryKey().String() }
func (t *Team) Table() string  { return TEAMS_TEAM_DB_TABLE }

// Insert Team instance to database
func (t *Team) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, t); err != nil {
		return
	}
	if err = Cache(ctx, t.String(), t); err != nil {
		return
	}
	return
}

// Updates team in database
func (t *Team) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	changed, err = DBUpdate(ctx, t, fields...)

	if err = Cache(ctx, t.String(), t); err != nil {
		return
	}

	return
}

// Deletes team from database
func (t *Team) Delete(ctx *context.Context) (err error) {
	cacheKey := t.String()
	if err = DBDelete(ctx, t); err != nil {
		return
	}
	if err = RemoveCached(ctx, cacheKey); err != nil {
		return
	}
	return

}

// cleans values
func (t *Team) Clean() {
	t.Name = strings.TrimSpace(t.Name)
}

// returns cache key

// validates data
func (t *Team) Validate(context *context.Context) (result *validator.Result, err error) {
	// clean values

	v := validator.New()
	v["name"] = validator.Any(
		validator.ValidateStringMinLength(4),
		validator.ValidateStringMaxLength(255),
	)
	result = v.Validate(t)

	// validate status
	if !t.Status.IsValid() {
		result.AddFieldError("status", ErrInvalidChoice)
	}
	return
}

func (t *Team) Manager(ctx *context.Context) *TeamManager {
	return NewTeamManager(ctx)
}

const (
	// teams migrations
	MIGRATION_TEAMS_TEAM_INITIAL_ID = "initial-migration-teams-team"
	MIGRATION_TEAMS_TEAM_INITIAL    = `CREATE TABLE ` + TEAMS_TEAM_DB_TABLE + `
	(
		id bigserial NOT NULL PRIMARY KEY,
		name character varying(200) NOT NULL,
		owner_id integer REFERENCES ` + AUTH_USER_DB_TABLE + ` ON DELETE RESTRICT,
		date_added timestamp with time zone NOT NULL,
		status integer CHECK (status > 0)
	)`
)

/*
Team manager
*/
func NewTeamManager(context *context.Context) *TeamManager {
	return &TeamManager{context: context}
}

type TeamManager struct {
	Manager
	context *context.Context
}

// returns new team
func NewTeam(funcs ...func(*Team)) (team *Team) {
	team = &Team{
		DateAdded: utils.NowTruncated(),
		Status:    TEAM_STATUS_VISIBLE,
	}
	for _, f := range funcs {
		f(team)
	}
	return
}

// returns new team
func (t *TeamManager) NewTeam(funcs ...func(*Team)) (team *Team) {
	return NewTeam(funcs...)
}

func (t *TeamManager) NewTeamList() []*Team {
	return []*Team{}
}

func (t *TeamManager) FilterPaged(target interface{}, paging *utils.Paging, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(t.context, TEAMS_TEAM_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, t.QueryFilterPaging(paging))

	_, safe := target.([]*Team)

	return DBFilter(t.context, TEAMS_TEAM_DB_TABLE+".*", TEAMS_TEAM_DB_TABLE, !safe, target, qfs...)
}

func (t *TeamManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*Team)

	return DBFilter(t.context, TEAMS_TEAM_DB_TABLE+".*", TEAMS_TEAM_DB_TABLE, !safe, target, qfs...)
}

// returns teams for given user
func (t *TeamManager) FilterByUser(target interface{}, user *User) (err error) {
	qfs := []utils.QueryFunc{}

	if !user.IsSuperuser {
		qfs = append(qfs, func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
			join := fmt.Sprintf(
				"%s ON (%s.team_id = %s.id)",
				TEAMS_TEAMMEMBER_DB_TABLE,
				TEAMS_TEAMMEMBER_DB_TABLE,
				TEAMS_TEAM_DB_TABLE,
			)

			return builder.Join(join).Where(fmt.Sprintf("%s.id = ?", TEAMS_TEAMMEMBER_DB_TABLE), user.ID)
		})
	}

	return t.Filter(target, qfs...)
}

func (t *TeamManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*Team)

	return DBGet(t.context, TEAMS_TEAM_DB_TABLE+".*", TEAMS_TEAM_DB_TABLE, !safe, target, qfs...)
}

func (t *TeamManager) GetByID(target interface{}, id types.Keyer) (err error) {
	if id.Int64() == 0 {
		return ErrObjectDoesNotExists
	}

	cacheKey := t.NewTeam(func(team *Team) { team.SetPrimaryKey(id) }).String()

	if err = GetCached(t.context, cacheKey, target); err == nil {
		// cache hit
		return
	}

	if err = t.Get(target, t.QueryFilterID(id)); err != nil {
		return
	}

	// it's safe to cache
	if _, ok := target.(*Team); ok {
		if err = Cache(t.context, cacheKey, target); err != nil {
			return
		}
	}

	return
}

func (t *TeamManager) GetFromRequest(target interface{}, r *http.Request, muxvar ...string) (err error) {
	v := "team_id"
	if len(muxvar) > 0 {
		v = muxvar[0]
	}
	var primaryKey types.PrimaryKey
	vars := mux.Vars(r)

	if err = primaryKey.Parse(vars[v]); err != nil {
		return
	}

	if err = t.GetByID(target, primaryKey); err != nil {
		return
	}

	return
}

/*
QueryFilterUser

returns teams that user can see
if user.IsSuperuser, all teams are available
if not we make join to teammember table
*/
func (t *TeamManager) QueryFilterUser(user *User) utils.QueryFunc {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		// superuser can do anything
		if user.IsSuperuser {
			return builder
		}
		join := fmt.Sprintf(
			"%s ON (%s.team_id = %s.id)",
			TEAMS_TEAMMEMBER_DB_TABLE,
			TEAMS_TEAMMEMBER_DB_TABLE,
			TEAMS_TEAM_DB_TABLE,
		)

		return builder.Join(join).Where(TEAMS_TEAMMEMBER_DB_TABLE+".user_id = ?", user.ID)
	}
}
