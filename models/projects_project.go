package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	sq "github.com/lann/squirrel"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

/*
	Project model
*/
type Project struct {
	Model
	Name      string           `db:"name" json:"name"`
	DateAdded time.Time        `db:"date_added" json:"date_added"`
	Platform  string           `db:"platform" json:"platform"`
	TeamID    types.ForeignKey `db:"team_id" json:"team_id"`
}

// returns all columns except of primary key
func (p *Project) Columns() []string { return []string{"name", "date_added", "platform", "team_id"} }
func (p *Project) Values() []interface{} {
	return []interface{}{p.Name, p.DateAdded, p.Platform, p.TeamID}
}
func (p *Project) String() string { return "projects:project:" + p.PrimaryKey().String() }
func (p *Project) Table() string  { return PROJECTS_PROJECT_DB_TABLE }

/*
CRUD operations
*/
func (p *Project) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, p); err != nil {
		return
	}

	// cache instance
	return Cache(ctx, p.String(), p)
}

func (p *Project) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	if changed, err = DBUpdate(ctx, p, fields...); err != nil {
		return
	}
	// cache instance
	if err = Cache(ctx, p.String(), p); err != nil {
		return
	}

	return
}

func (p *Project) Delete(ctx *context.Context) (err error) {
	cacheKey := p.String()
	if err = DBDelete(ctx, p); err != nil {
		return
	}
	return RemoveCached(ctx, cacheKey)
}

/*
Validations
*/
func (p *Project) Clean(context *context.Context) {
	p.Name = strings.TrimSpace(p.Name)
}

func (p *Project) Validate(ctx *context.Context) (result *validator.Result, err error) {
	p.Clean(ctx)
	result = validator.NewResult()

	if p.TeamID == 0 {
		result.AddFieldError("team_id", errors.New("invalid_value"))
	}

	// validate team_id
	// validator.ValidateIntColumn(p.Table(), "team_id", p.TeamID, func(tid int64) (err error) {
	// 	teammanager := NewTeamManager(context)
	// 	if errGet := teammanager.GetByID(teammanager.NewTeam(), tid); err != ErrObjectDoesNotExists {
	// 		return errGet
	// 	}
	// 	return
	// })

	return
}

// returns team
func (p *Project) Team(target interface{}, manager *TeamManager) (err error) {
	if p.TeamID.Int64() == 0 {
		return ErrObjectDoesNotExists
	}
	return manager.GetByID(target, &p.TeamID)
}

func (p *Project) Manager(ctx *context.Context) *ProjectManager {
	return NewProjectManager(ctx)
}

const (
	// project migrations
	MIGRATION_PROJECT_INITIAL_ID = "initial-migration-project"
	MIGRATION_PROJECT_INITIAL    = `CREATE TABLE ` + PROJECTS_PROJECT_DB_TABLE + `
	(
		id bigserial NOT NULL PRIMARY KEY,
		name character varying(200) NOT NULL,
		date_added timestamp with time zone NOT NULL,
		platform character varying(32),
		team_id bigint REFERENCES ` + TEAMS_TEAM_DB_TABLE + ` ON DELETE RESTRICT
	)`
)

var (
	MIGRATION_PROJECT_INITIAL_DEPENDENCIES = []string{
		settings.TEAMS_PLUGIN_ID + ":" + MIGRATION_TEAMS_TEAM_INITIAL_ID,
		settings.AUTH_PLUGIN_ID + ":" + MIGRATION_AUTH_PERMISSION_INITIAL_ID,
	}
)

/*
	ProjectManager
*/
type ProjectManager struct {
	Manager
	context *context.Context
}

/*
	Constructor for project manager
*/
func NewProjectManager(context *context.Context, tx ...*sqlx.Tx) *ProjectManager {
	return &ProjectManager{context: context}
}

/*
	Returns blank new Project
*/
func NewProject(funcs ...func(*Project)) (project *Project) {
	project = &Project{
		DateAdded: utils.NowTruncated(),
	}
	for _, f := range funcs {
		f(project)
	}
	return
}

/*
	Returns blank new Project
*/
func (p *ProjectManager) NewProject(funcs ...func(*Project)) (project *Project) {
	return NewProject(funcs...)
}

func (p *ProjectManager) NewProjectList() []*Project {
	return []*Project{}
}

// select without paging
func (p *ProjectManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*Project)
	return DBFilter(p.context, PROJECTS_PROJECT_DB_TABLE+".*", PROJECTS_PROJECT_DB_TABLE, !safe, target, qfs...)
}

/* Filters projects from database
qfs is list of QueryFuncs - that are functions that alter query builder
*/
func (p *ProjectManager) FilterPaged(target interface{}, paging *utils.Paging, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(p.context, PROJECTS_PROJECT_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, p.QueryFilterPaging(paging))

	_, safe := target.([]*Project)

	return DBFilter(p.context, PROJECTS_PROJECT_DB_TABLE+".*", PROJECTS_PROJECT_DB_TABLE, !safe, target, qfs...)
}

/* Returns project by Query filter funcs
 */
func (p *ProjectManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*Project)

	return DBGet(p.context, "*", PROJECTS_PROJECT_DB_TABLE, !safe, target, qfs...)
}

/* Returns project by ID
 */
func (p *ProjectManager) GetByID(target interface{}, id types.Keyer) (err error) {
	// var cached bool

	cacheKey := p.NewProject(func(pr *Project) { pr.SetPrimaryKey(id) }).String()
	if err = GetCached(p.context, cacheKey, target); err == nil {
		// cache hit
		return
	}

	if err = p.Get(target, p.QueryFilterID(id)); err != nil {
		return
	}

	//cache instance
	if _, ok := target.(*Project); ok {
		if err = Cache(p.context, cacheKey, target); err != nil {
			return
		}
	}

	return
}

func (p *ProjectManager) GetByAuth(target interface{}, PublicKey, SecretKey string) (err error) {
	pkm := NewProjectKeyManager(p.context)
	pk := pkm.NewProjectKey()
	err = pkm.GetByAuth(pk, PublicKey, SecretKey)
	if err != nil {

		return err
	}

	// get by id
	err = p.GetByID(target, pk.ProjectID)

	return
}

/*
Returns project from mux var
muxvar is optional with default "project_id"

*/
func (p *ProjectManager) GetFromRequest(target interface{}, r *http.Request, muxvar ...string) (err error) {
	v := "project_id"
	if len(muxvar) > 0 {
		v = muxvar[0]
	}
	var id int64
	vars := mux.Vars(r)
	if id, err = strconv.ParseInt(vars[v], 10, 0); err != nil {
		return
	}

	return p.GetByID(target, types.PrimaryKey(id))
}

/*
Queries projects visible to given user
*/
func (p *ProjectManager) QueryFilterUser(user *User) utils.QueryFunc {
	handleNilPointer(user)

	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		// superuser can see all projects
		if user.IsSuperuser {
			return builder
		}

		// inner join returns all of them
		joinTeam := fmt.Sprintf(
			"%s ON (%s.id = %s.team_id)",
			TEAMS_TEAM_DB_TABLE,
			TEAMS_TEAM_DB_TABLE,
			PROJECTS_PROJECT_DB_TABLE,
		)
		joinTeamMember := fmt.Sprintf(
			"%s ON (%s.team_id = %s.id)",
			TEAMS_TEAMMEMBER_DB_TABLE,
			TEAMS_TEAMMEMBER_DB_TABLE,
			TEAMS_TEAM_DB_TABLE,
		)
		builder = builder.Join(joinTeam).Join(joinTeamMember)
		return builder.Where(TEAMS_TEAMMEMBER_DB_TABLE+".user_id = ?", user.ID)
	}
}

/*
Returns member type for given user and project

For now we check only team member, in the future here is the place to add Organisation
member type.
*/
func (p *ProjectManager) MemberType(project *Project, user *User) (mt MemberType, err error) {
	tmManager := NewTeamMemberManager(p.context)

	// super user acts as admin
	if user.IsSuperuser {
		mt = MEMBER_TYPE_ADMIN
		return
	}

	// check member type
	if mt, err = tmManager.MemberTypeByProject(project, user); err != nil {
		return
	}

	return
}
