package models

import (
	"time"

	"github.com/lann/squirrel"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/paginator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

const (
	// team member migrations
	MIGRATION_TEAMS_TEAM_MEMBER_INITIAL_ID = "initial-migration-teams-teammember"
	MIGRATION_TEAMS_TEAM_MEMBER_INITIAL    = `CREATE TABLE ` + TEAMS_TEAMMEMBER_DB_TABLE + `(
		id bigserial NOT NULL PRIMARY KEY,
		team_id bigint NOT NULL REFERENCES ` + TEAMS_TEAM_DB_TABLE + ` ON DELETE CASCADE,
		user_id bigint NOT NULL REFERENCES ` + AUTH_USER_DB_TABLE + ` ON DELETE CASCADE,
		date_added timestamp with time zone,
		type integer CHECK (type > 0),
		UNIQUE (team_id, user_id)
	)`
)

/*
TeamMember
*/
type TeamMember struct {
	Model
	TeamID    types.ForeignKey `db:"team_id" json:"team_id"`
	UserID    types.ForeignKey `db:"user_id" json:"user_id"`
	DateAdded time.Time        `db:"date_added" json:"date_added"`
	Type      MemberType       `db:"type" json:"type"`
}

// returns all columns except of primary key
func (t *TeamMember) Columns() []string {
	return []string{"team_id", "user_id", "date_added", "type"}
}
func (t *TeamMember) Values() []interface{} {
	return []interface{}{t.TeamID, t.UserID, t.DateAdded, t.Type}
}
func (t *TeamMember) String() string { return "teams:teammember:" + t.PrimaryKey().String() }
func (t *TeamMember) Table() string  { return TEAMS_TEAMMEMBER_DB_TABLE }

/*
CRUD
*/
func (t *TeamMember) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, t); err != nil {
		return
	}
	return
}

func (t *TeamMember) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	changed, err = DBUpdate(ctx, t, fields...)
	return
}

func (t *TeamMember) Delete(ctx *context.Context) (err error) {
	return DBDelete(ctx, t)
}

func (t *TeamMember) Manager(ctx *context.Context) *TeamMemberManager {
	return NewTeamMemberManager(ctx)
}

func (t *TeamMember) User(target interface{}, ctx *context.Context) (err error) {
	manager := NewUserManager(ctx)
	err = manager.GetByID(target, t.UserID)
	return
}

/*
TeamMember manager
*/
func NewTeamMemberManager(context *context.Context) *TeamMemberManager {
	return &TeamMemberManager{context: context}
}

type TeamMemberManager struct {
	Manager
	context *context.Context
}

// Returns new TeamMember with default values
func NewTeamMember(funcs ...func(*TeamMember)) (tm *TeamMember) {
	tm = &TeamMember{
		DateAdded: utils.NowTruncated(),
		Type:      MEMBER_TYPE_MEMBER,
	}
	for _, f := range funcs {
		f(tm)
	}
	return
}

func NewTeamMemberList() []*TeamMember {
	return []*TeamMember{}
}

// Returns new TeamMember with default values
func (t *TeamMemberManager) NewTeamMember(funcs ...func(*TeamMember)) (tm *TeamMember) {
	return NewTeamMember(funcs...)
}

// return blank list of TeamMemgers for filter
func (t *TeamMemberManager) NewTeamMemberList() []*TeamMember {
	return NewTeamMemberList()
}

// Filter team memebers from database with paging support
func (t *TeamMemberManager) FilterPaged(target interface{}, paging *paginator.Paginator, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(t.context, TEAMS_TEAMMEMBER_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, t.QueryFilterPaging(paging))

	_, safe := target.([]*TeamMember)

	return DBFilter(t.context, TEAMS_TEAMMEMBER_DB_TABLE+".*", TEAMS_TEAMMEMBER_DB_TABLE, !safe, target, qfs...)
}

// Filter team memebers from database
func (t *TeamMemberManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*TeamMember)

	return DBFilter(t.context, TEAMS_TEAMMEMBER_DB_TABLE+".*", TEAMS_TEAMMEMBER_DB_TABLE, !safe, target, qfs...)
}

func (t *TeamMemberManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*TeamMember)

	return DBGet(t.context, TEAMS_TEAMMEMBER_DB_TABLE+".*", TEAMS_TEAMMEMBER_DB_TABLE, !safe, target, qfs...)
}

func (t *TeamMemberManager) GetByID(target interface{}, id types.Keyer) (err error) {
	if id.Int64() == 0 {
		return ErrObjectDoesNotExists
	}

	if err = t.Get(target, t.QueryFilterID(id)); err != nil {
		return
	}

	return
}

/*
Adds member to team
*/
func (t *TeamMemberManager) SetTeamMemberType(team *Team, user *User, mt MemberType) (result *TeamMember, err error) {
	handleNilPointer(team)
	handleNilPointer(user)

	result = t.NewTeamMember()

	if err = t.Get(result, t.QueryFilterTeamUser(team, user)); err != nil {
		result.TeamID = types.ForeignKey(team.ID)
		result.UserID = types.ForeignKey(user.ID)
		result.Type = mt
		err = result.Insert(t.context)
		return
	}

	result.Type = mt
	_, err = result.Update(t.context, "type")

	return
}

/*
Removes team member
*/
func (t *TeamMemberManager) RemoveTeamMember(team *Team, user *User) (err error) {
	handleNilPointer(team)
	handleNilPointer(user)

	qb := utils.QueryBuilder().Delete(TEAMS_TEAMMEMBER_DB_TABLE).Where("team_id = ? AND user_id = ?", team.ID, user.ID)

	var (
		query string
		args  []interface{}
	)

	if query, args, err = qb.ToSql(); err != nil {
		return
	}

	_, err = t.context.DB.Exec(query, args...)
	return
}

/*
Returns member type for given project and user
*/
func (t *TeamMemberManager) MemberType(team *Team, user *User) (mt MemberType, err error) {
	// handle nil pointers
	handleNilPointer(team)
	handleNilPointer(team)

	// check
	if team.ID == 0 || user.ID == 0 {
		err = ErrObjectDoesNotExists
		return
	}

	// handle superuser correctly
	if user.IsSuperuser {
		return MEMBER_TYPE_ADMIN, nil
	}

	tm := t.NewTeamMember()

	// create QueryFilterUserTeam
	if err = t.Get(tm, t.QueryFilterTeamUser(team, user)); err != nil {
		return
	}

	return tm.Type, nil
}

/*
Returns member type for given project and user
*/
func (t *TeamMemberManager) MemberTypeByProject(project *Project, user *User) (mt MemberType, err error) {
	// handle nil pointers
	handleNilPointer(project)
	handleNilPointer(user)

	tm := NewTeamManager(t.context)
	team := tm.NewTeam()
	if err = project.Team(team, tm.context); err != nil {
		return
	}

	if mt, err = t.MemberType(team, user); err != nil {
		return
	}

	return
}

// query by user and team
func (t *TeamMemberManager) QueryFilterTeamUser(team *Team, user *User) utils.QueryFunc {
	handleNilPointer(team)
	handleNilPointer(user)

	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where("team_id = ? AND user_id = ?", team.ID, user.ID)
	}
}

func (t *TeamMemberManager) QueryFilterProject(project *Project) utils.QueryFunc {
	handleNilPointer(project)
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where("team_id = ?", project.TeamID)
	}
}
