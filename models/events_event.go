package models

import (
	"time"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/parser"

	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

/*
Migrations
*/
var (
	MIGRATION_EVENTS_EVENT_INITIAL_ID = "events-event-initial"
	MIGRATION_EVENTS_EVENT_INITIAL    = `CREATE TABLE ` + EVENTS_EVENT_DB_TABLE + `
	(
		id bigserial NOT NULL PRIMARY KEY,
		event_id character varying (32),
		event_group_id bigint REFERENCES ` + EVENTS_EVENTGROUP_DB_TABLE + `,
		project_id bigint REFERENCES ` + PROJECTS_PROJECT_DB_TABLE + `,
		message character varying (1000) NOT NULL,
		platform character varying (64) NOT NULL,
		datetime timestamp with time zone NOT NULL,
		time_spent bigint NOT NULL,
		data bytea NOT NULL
	)`
)

/*
Event model
*/
type Event struct {
	Model
	EventID      string           `db:"event_id" json:"event_id"`
	EventGroupID types.ForeignKey `db:"event_group_id" json:"event_group_id"`
	ProjectID    types.ForeignKey `db:"project_id" json:"project_id"`
	Message      string           `db:"message" json:"message"`
	Platform     string           `db:"platform" json:"platform"`
	Datetime     time.Time        `db:"date_time" json:"datetime"`
	TimeSpent    int64            `db:"time_spent" json:"time_spent"`
	Data         types.GzippedMap `db:"data" json:"data"`
}

// returns all columns except of primary key
func (e *Event) Columns() []string {
	return []string{
		"event_id", "event_group_id", "project_id", "message",
		"platform", "datetime", "time_spent", "data",
	}
}
func (e *Event) Values() []interface{} {
	return []interface{}{
		e.EventID, e.EventGroupID, e.ProjectID, e.Message,
		e.Platform, e.Datetime, e.TimeSpent, e.Data,
	}
}
func (e *Event) String() string { return "events:event:" + e.PrimaryKey().String() }
func (e *Event) Table() string  { return EVENTS_EVENT_DB_TABLE }

/*
CRUD
*/
func (e *Event) Insert(ctx *context.Context) (err error) {
	if err = DBInsert(ctx, e); err != nil {
		return
	}
	return
}

func (e *Event) Update(ctx *context.Context, fields ...string) (changed bool, err error) {
	changed, err = DBUpdate(ctx, e, fields...)
	return
}

func (e *Event) Delete(ctx *context.Context) (err error) {
	return DBDelete(ctx, e)
}

func (e *Event) Manager(ctx *context.Context) *EventManager {
	return NewEventManager(ctx)
}

/*
Event manager
All sql queries are performed on managers
*/
type EventManager struct {
	Manager
	context *context.Context
}

// constructor function to create new manager
func NewEventManager(context *context.Context) *EventManager {
	em := &EventManager{context: context}
	return em
}

/*
	New... methods to create blank model instances
*/
func NewEvent(funcs ...func(*Event)) (event *Event) {
	event = &Event{
		Datetime: utils.NowTruncated(),
	}
	for _, f := range funcs {
		f(event)
	}
	return
}

/*
	New... methods to create blank model instances
*/
func (e *EventManager) NewEvent(funcs ...func(*Event)) (event *Event) {
	return NewEvent(funcs...)
}
func (e *EventManager) NewEventList() []*Event { return []*Event{} }

// Filter results without paging
func (e *EventManager) Filter(target interface{}, qfs ...utils.QueryFunc) error {
	_, safe := target.([]*Event)
	return DBFilter(e.context, EVENTS_EVENT_DB_TABLE+".*", EVENTS_EVENT_DB_TABLE, !safe, target, qfs...)
}

// Filter results with paging
func (e *EventManager) FilterPaged(target interface{}, paging *utils.Paging, qfs ...utils.QueryFunc) (err error) {
	if err = DBFilterCount(e.context, EVENTS_EVENT_DB_TABLE, paging, qfs...); err != nil {
		return
	}

	// add paging query filter
	qfs = append(qfs, e.QueryFilterPaging(paging))

	_, safe := target.([]*Event)

	return DBFilter(e.context, EVENTS_EVENT_DB_TABLE+".*", EVENTS_EVENT_DB_TABLE, !safe, target, qfs...)
}

// get from database
func (e *EventManager) Get(target interface{}, qfs ...utils.QueryFunc) (err error) {
	_, safe := target.(*Event)
	return DBGet(e.context, "*", EVENTS_EVENT_DB_TABLE, !safe, target, qfs...)
}

// returns by id )and possibly other queryFuncs
func (e *EventManager) GetByID(target interface{}, id types.Keyer, qfs ...utils.QueryFunc) (err error) {
	// add by id
	qfs = append(qfs, e.QueryFilterWhere("id = ?", id.Int64()))
	return e.Get(target, qfs...)
}

// NewEventFromRaw creates new event from raw event
func (e *EventManager) NewEventFromRaw(raw *parser.RawEvent) (event *Event, eventgroup *EventGroup, err error) {
	egm := NewEventGroupManager(e.context)

	// some serious error occured
	if eventgroup, err = egm.GetByRaw(raw); err != nil {
		return
	}

	event = e.NewEvent(func(ev *Event) {
		ev.EventID = raw.EventID
		ev.EventGroupID = eventgroup.ID.ToForeignKey()
		ev.ProjectID = eventgroup.ProjectID
		ev.Message = raw.Message
		ev.Platform = raw.Platform
		ev.Datetime = utils.NowTruncated()
		ev.Data = raw.Data
	})

	// some serious error occured
	if err = event.Insert(e.context); err != nil {
		return
	}

	return
}
