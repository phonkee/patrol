package models

import (
	"strconv"
	"time"

	"github.com/lann/squirrel"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/parser"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

type EventGroup struct {
	Model
	ProjectID      types.ForeignKey `db:"project_id" json:"project_id"`
	Logger         string           `db:"logger" json:"logger"`
	Level          int              `db:"level" json:"level"`
	Message        string           `db:"message" json:"message"`
	Culprit        string           `db:"culprit" json:"culprit"`
	Checksum       string           `db:"checksum" json:"checksum"`
	Platform       string           `db:"platform" json:"platform"`
	Status         EventGroupStatus `db:"status" json:"status"`
	TimesSeen      int64            `db:"times_seen" json:"times_seen"`
	FirstSeen      time.Time        `db:"first_seen" json:"first_seen"`
	LastSeen       time.Time        `db:"last_seen" json:"last_seen"`
	ResolvedAt     time.Time        `db:"resolved_at" json:"resolved_at"`
	ActiveAt       time.Time        `db:"active_at" json:"active_at"`
	TimeSpentTotal int              `db:"time_spent_total" json:"time_spent_total"`
	TimeSpentCount int              `db:"time_spent_count" json:"time_spent_count"`
	Score          int              `db:"score" json:"score"`
	Data           types.GzippedMap `db:"data" json:"data"`
}

// returns all columns except of primary key
func (e *EventGroup) Columns() []string {
	return []string{
		"project_id", "logger", "level", "message", "culprit",
		"checksum", "platform", "status", "times_seen", "first_seen",
		"last_seen", "resolved_at", "active_at", "time_spent_total",
		"time_spent_count", "score", "data",
	}
}
func (e *EventGroup) Values() []interface{} {
	return []interface{}{
		e.ProjectID, e.Logger, e.Level, e.Message, e.Culprit,
		e.Checksum, e.Platform, e.Status, e.TimesSeen, e.FirstSeen,
		e.LastSeen, e.ResolvedAt, e.ActiveAt, e.TimeSpentTotal,
		e.TimeSpentCount, e.Score, e.Data,
	}
}
func (e *EventGroup) String() string { return "events:eventgroup:" + e.PrimaryKey().String() }
func (e *EventGroup) Table() string  { return EVENTS_EVENTGROUP_DB_TABLE }

/*
CRUD
*/
func (e *EventGroup) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, e); err != nil {
		return
	}
	if err = Cache(ctx, e.String(), e); err != nil {
		return
	}
	return
}

func (e *EventGroup) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	changed, err = DBUpdate(ctx, e, fields...)

	if err = Cache(ctx, e.String(), e); err != nil {
		return
	}

	return
}

func (e *EventGroup) Delete(ctx *context.Context) (err error) {
	cacheKey := e.String()
	if err = DBDelete(ctx, e); err != nil {
		return
	}

	if err = RemoveCached(ctx, cacheKey); err != nil {
		return
	}
	return
}

// returns manager
func (e *EventGroup) Manager(ctx *context.Context) *EventGroupManager {
	return NewEventGroupManager(ctx)
}

/*
Validates eventgroup
*/
func (e *EventGroup) Validate(context *context.Context) (result *validator.Result, err error) {
	result = validator.NewResult()
	if !e.Status.IsValid() {
		result.AddFieldError("status", ErrInvalidChoice)
	}

	return result, nil
}

var (
	MIGRATION_EVENTS_EVENTGROUP_INITIAL_ID = "events-eventgroup-initial"
	MIGRATION_EVENTS_EVENTGROUP_INITIAL    = `CREATE TABLE ` + EVENTS_EVENTGROUP_DB_TABLE + `
    (
        id bigserial NOT NULL PRIMARY KEY,
        project_id bigint REFERENCES ` + PROJECTS_PROJECT_DB_TABLE + `,
        logger character varying (64) NOT NULL,
        level integer NOT NULL,
        message text NOT NULL,
        culprit character varying(` + strconv.Itoa(MAX_CULPRIT_LENGTH) + `),
        checksum character varying(32) NOT NULL,
        platform character varying (64) NOT NULL,
        status integer CHECK (status > 0),
        times_seen integer NOT NULL,
        first_seen timestamp with time zone NOT NULL,
        last_seen timestamp with time zone NOT NULL,
        resolved_at timestamp with time zone NOT NULL,
        active_at timestamp with time zone NOT NULL,
        time_spent_total integer NOT NULL,
        time_spent_count integer NOT NULL,
        score integer NOT NULL,
        data bytea NOT NULL,
        UNIQUE (project_id, checksum)
    )`
	MIGRATION_EVENTS_EVENTGROUP_INITIAL_DEPENDENCIES = []string{settings.PROJECTS_PLUGIN_ID + ":" + MIGRATION_PROJECT_INITIAL_ID}
)

/*
EventGroupManager
*/
type EventGroupManager struct {
	Manager
	context *context.Context
}

func NewEventGroupManager(context *context.Context) *EventGroupManager {
	return &EventGroupManager{context: context}
}

// returns new model instance
func NewEventGroup(funcs ...func(eg *EventGroup)) (eg *EventGroup) {
	eg = &EventGroup{}
	for _, f := range funcs {
		f(eg)
	}
	return
}

// returns new model instance
func (e *EventGroupManager) NewEventGroup(funcs ...func(eg *EventGroup)) (eg *EventGroup) {
	return NewEventGroup(funcs...)
}
func (e *EventGroupManager) NewEventGroupList() []*EventGroup { return []*EventGroup{} }

// Filter results without paging
func (e *EventGroupManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*EventGroup)
	return DBFilter(e.context, EVENTS_EVENTGROUP_DB_TABLE+".*", EVENTS_EVENTGROUP_DB_TABLE, !safe, target, qfs...)
}

// Filter results with paging
func (e *EventGroupManager) FilterPaged(target interface{}, paging *utils.Paging, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(e.context, EVENTS_EVENTGROUP_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, e.QueryFilterPaging(paging))

	_, safe := target.([]*EventGroup)

	return DBFilter(e.context, EVENTS_EVENTGROUP_DB_TABLE+".*", EVENTS_EVENTGROUP_DB_TABLE, !safe, target, qfs...)
}

// get from database
func (e *EventGroupManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*EventGroup)
	return DBGet(e.context, "*", EVENTS_EVENTGROUP_DB_TABLE, !safe, target, qfs...)
}

// returns by id )and possibly other queryFuncs
func (e *EventGroupManager) GetByID(target interface{}, id int64, qfs ...utils.QueryFunc) (err error) {
	handleNilPointer(target)

	cacheKey := e.NewEventGroup(func(eg *EventGroup) { eg.SetPrimaryKey(types.PrimaryKey(id)) }).String()

	if err = GetCached(e.context, cacheKey, target); err == nil {
		return
	}

	qfs = append(qfs, e.QueryFilterWhere("id = ?", id))
	if err = e.Get(target, qfs...); err != nil {
		return
	}

	// cache it's safe
	if _, ok := target.(*EventGroup); ok {
		if err = Cache(e.context, cacheKey, target); err != nil {
			return
		}
	}

	return
}

// returns eventgroup by raw event (returns from db or craetes one)
func (e *EventGroupManager) GetByRaw(raw *parser.RawEvent) (eventgroup *EventGroup, err error) {

	eventgroup = e.NewEventGroup()

	if err = e.Get(eventgroup, utils.QueryFilterWhere("checksum = ? AND project_id = ?", raw.Checksum, raw.ProjectID)); err != nil {
		if err != ErrObjectDoesNotExists {
			return
		}

		// create new eventgroup
		eventgroup = e.NewEventGroup(func(eg *EventGroup) {
			eg.ProjectID = raw.ProjectID
			eg.Logger = raw.Logger
			// @TODO: map from string to int value
			// eg.Level = raw.Level
			eg.Message = raw.Message
			eg.Culprit = raw.Culprit
			eg.Checksum = raw.Checksum
			eg.Platform = raw.Platform
			eg.Status = EVENT_GROUP_STATUS_UNRESOLVED
			//eg.TimesSeen = 0
			eg.FirstSeen = utils.NowTruncated()
			eg.LastSeen = utils.NowTruncated()
			eg.Data = raw.Data
		})

		// something bad happened
		if err = eventgroup.Insert(e.context); err != nil {
			return
		}
	}

	return
}

// increments counter safe way
func (e *EventGroupManager) IncrementCounters(eventgroup *EventGroup) (err error) {
	builder := utils.QueryBuilder().
		Update(eventgroup.Table()).
		Set("times_seen", squirrel.Expr("times_seen + 1")).
		Suffix("RETURNING times_seen")

	query, args, err := builder.ToSql()

	qrfunc := e.context.DB.QueryRow
	if e.context.Tx != nil {
		qrfunc = e.context.Tx.QueryRow
	}

	if err = qrfunc(query, args...).Scan(&eventgroup.TimesSeen); err != nil {
		return
	}

	return
}
