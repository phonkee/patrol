package models

import (
	"time"

	sq "github.com/lann/squirrel"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

const (
	PROJECT_KEY_PUBLIC_KEY_LENGTH = 32
	PROJECT_KEY_SECRET_KEY_LENGTH = 32

	// project key migrations
	MIGRATION_PROJECT_KEY_INITIAL_ID = "initial-migration-project-key"
	MIGRATION_PROJECT_KEY_INITIAL    = `CREATE TABLE ` + PROJECTS_PROJECTKEY_DB_TABLE + `(
		id bigserial NOT NULL PRIMARY KEY,
		project_id bigint NOT NULL,
		public_key character varying(32),
		secret_key character varying(32),
		user_id integer,
		user_added_id integer,
		date_added timestamp with time zone
	)`
)

/*
	ProjectKey model
	Every project must have at least one ProjectKey
*/
type ProjectKey struct {
	Model
	ProjectID   types.ForeignKey `db:"project_id" json:"project_id"`
	PublicKey   string           `db:"public_key" json:"public_key"`
	SecretKey   string           `db:"secret_key" json:"secret_key"`
	UserID      types.ForeignKey `db:"user_id" json:"user_id"`
	UserAddedID types.ForeignKey `db:"user_added_id" json:"user_added_id"`
	DateAdded   time.Time        `db:"date_added" json:"date_added"`
}

// returns all columns except of primary key
func (p *ProjectKey) Columns() []string {
	return []string{
		"project_id", "public_key", "secret_key", "user_id", "user_added_id",
		"date_added",
	}
}
func (p *ProjectKey) Values() []interface{} {
	return []interface{}{
		p.ProjectID, p.PublicKey, p.SecretKey, p.UserID, p.UserAddedID,
		p.DateAdded,
	}
}
func (p *ProjectKey) String() string { return "projects:projectkey:" + p.PrimaryKey().String() }
func (p *ProjectKey) Table() string  { return PROJECTS_PROJECTKEY_DB_TABLE }

/*
CRUD operations
*/
func (p *ProjectKey) Insert(ctx *context.Context) error {
	return DBInsert(ctx, p)
}

func (p *ProjectKey) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	return DBUpdate(ctx, p, fields...)
}

func (p *ProjectKey) Delete(ctx *context.Context) error {
	return DBDelete(ctx, p)
}

func (p *ProjectKey) Manager(ctx *context.Context) *ProjectKeyManager {
	return NewProjectKeyManager(ctx)
}

/*
	Validates project key instance
*/
func (p *ProjectKey) Validate(context *context.Context) (result *validator.Result, err error) {
	result = validator.NewResult()
	return
}

/* Constructor function for ProjectKeyManager
 */
func NewProjectKeyManager(context *context.Context) *ProjectKeyManager {
	return &ProjectKeyManager{context: context}
}

/* ProjectKey manager
 */
type ProjectKeyManager struct {
	Manager
	context *context.Context
}

/* Returns newly prepared ProjectKey instance
 */
func NewProjectKey(funcs ...func(*ProjectKey)) (pk *ProjectKey) {
	pk = &ProjectKey{
		DateAdded: utils.NowTruncated(),
		PublicKey: utils.RandomString(PROJECT_KEY_PUBLIC_KEY_LENGTH),
		SecretKey: utils.RandomString(PROJECT_KEY_SECRET_KEY_LENGTH),
	}

	for _, f := range funcs {
		f(pk)
	}

	return
}

func NewProjectKeyList() (pk []*ProjectKey) {
	return []*ProjectKey{}
}

/* Returns newly prepared ProjectKey instance
 */
func (p *ProjectKeyManager) NewProjectKey(funcs ...func(*ProjectKey)) (pk *ProjectKey) {
	return NewProjectKey(funcs...)
}

/* Returns newly prepared ProjectKey instance
 */
func (p *ProjectKeyManager) NewProjectKeyList() (pk []*ProjectKey) {
	return NewProjectKeyList()
}

// select without paging
func (p *ProjectKeyManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*ProjectKey)
	return DBFilter(p.context, PROJECTS_PROJECTKEY_DB_TABLE+".*", PROJECTS_PROJECTKEY_DB_TABLE, !safe, target, qfs...)
}

/* Filters projects from database
qfs is list of QueryFuncs - that are functions that alter query builder
*/
func (p *ProjectKeyManager) FilterPaged(target interface{}, paging *utils.Paging, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(p.context, PROJECTS_PROJECTKEY_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, p.QueryFilterPaging(paging))

	_, safe := target.([]*ProjectKey)

	return DBFilter(p.context, PROJECTS_PROJECTKEY_DB_TABLE+".*", PROJECTS_PROJECT_DB_TABLE, !safe, target, qfs...)
}

/* Returns project by Query filter funcs
 */
func (p *ProjectKeyManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*ProjectKey)

	return DBGet(p.context, "*", PROJECTS_PROJECTKEY_DB_TABLE, !safe, target, qfs...)
}

func (p *ProjectKeyManager) GetByAuth(target interface{}, PublicKey, SecretKey string) (err error) {
	err = p.Get(target, p.QueryFilterWhere("secret_key = ? AND public_key = ?", SecretKey, PublicKey))
	return
}

func (p *ProjectKeyManager) QueryFilterProjectID(projectid int64) utils.QueryFunc {
	pk := p.NewProjectKey()
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(pk.Table()+".project_id = ?", projectid)
	}
}
