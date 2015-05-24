package models

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/paginator"
	"github.com/phonkee/patrol/utils"
)

/*
Migrations
*/
const (
	MIGRATION_AUTH_PERMISSION_INITIAL_ID = "initial-migration-auth-permissions"
	MIGRATION_AUTH_PERMISSION_INITIAL    = `CREATE TABLE ` + AUTH_PERMISSION_DB_TABLE + `
	(
		id bigserial PRIMARY KEY,
		codename character varying(64) NOT NULL UNIQUE,
		name character varying(255) NOT NULL
	)`
)

/*
Permissions
*/
type Permission struct {
	Model
	Codename string `db:"codename" json:"codename"`
	Name     string `db:"name" json:"name"`
}

// returns all columns except of primary key
func (p *Permission) Columns() []string     { return []string{"codename", "name"} }
func (p *Permission) Values() []interface{} { return []interface{}{p.Codename, p.Name} }
func (p *Permission) String() string        { return "auth:permission:" + p.ID.String() }
func (p *Permission) Table() string         { return AUTH_PERMISSION_DB_TABLE }

/*
CRUD
*/
func (p *Permission) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, p); err != nil {
		return
	}
	return
}

func (p *Permission) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	changed, err = DBUpdate(ctx, p, fields...)
	return
}

func (p *Permission) Delete(ctx *context.Context) (err error) {
	return DBDelete(ctx, p)
}

func (p *Permission) Manager(ctx *context.Context) *PermissionManager {
	return NewPermissionManager(ctx)
}

type PermissionManager struct {
	Manager
	context *context.Context
}

func NewPermissionManager(context *context.Context) *PermissionManager {
	pm := &PermissionManager{
		context: context,
	}
	return pm
}
func NewPermission(funcs ...func(permission *Permission)) (result *Permission) {
	result = &Permission{}
	for _, f := range funcs {
		f(result)
	}
	return
}
func (p *PermissionManager) NewPermission(funcs ...func(permission *Permission)) (result *Permission) {
	return NewPermission(funcs...)
}
func (p *PermissionManager) NewPermissionList() []*Permission { return []*Permission{} }

// select without paging
func (p *PermissionManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*Permission)
	return DBFilter(p.context, AUTH_PERMISSION_DB_TABLE+".*", AUTH_PERMISSION_DB_TABLE, !safe, target, qfs...)
}

/* Filters permissions from database
qfs is list of utils.QueryFuncs - that are functions that alter query builder
*/
func (p *PermissionManager) FilterPaged(target interface{}, paging *paginator.Paginator, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(p.context, AUTH_PERMISSION_DB_TABLE, paging, qfs...); err != nil {
		return
	}
	// add paging query filter
	qfs = append(qfs, p.QueryFilterPaging(paging))
	_, safe := target.([]*Permission)
	return DBFilter(p.context, AUTH_PERMISSION_DB_TABLE+".*", AUTH_PERMISSION_DB_TABLE, !safe, target, qfs...)
}

/* Returns project by Query filter funcs
 */
func (p *PermissionManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*Permission)
	return DBGet(p.context, "*", AUTH_PERMISSION_DB_TABLE, !safe, target, qfs...)
}

// Returns Permission by codename
func (p *PermissionManager) GetByCodename(target interface{}, codename string) error {
	return p.Get(target, p.QueryFilterWhere("codename = ?", codename))
}
