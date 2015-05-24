package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/golang/glog"
	"github.com/mgutz/ansi"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/paginator"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/utils"
)

const (
	// runtime.caller skip value - do it points to correct code
	SQL_CALLER_SKIP         = 3
	NIL_POINTER_CALLER_SKIP = 3
)

var (
	// V2 logging colorize
	alert       = ansi.ColorFunc("red+h:black")
	phosphorize = ansi.ColorFunc("green+h:black")
)

/*
Handles situations if nil pointer is passed
*/
func handleNilPointer(value interface{}) (err error) {
	if value == nil {
		if _, file, line, ok := runtime.Caller(NIL_POINTER_CALLER_SKIP); ok {
			err = fmt.Errorf("nil pointer passed %s:%s", file, line)
		} else {
			err = fmt.Errorf("nil pointer")
		}
		if settings.DEBUG {
			panic(err)
		} else {
			glog.Error(err)
		}
	}
	return
}

/*
DBFilter method implementation

All Managers can use this "general" method for filtering results.
In case of something more complicated please provide your own implementation
*/
func DBFilter(ctx *context.Context, sel, dbtable string, unsafe bool, target interface{}, qfs ...utils.QueryFunc) (err error) {
	if target == nil {
		return ErrNilPointer
	}

	selfunc := ctx.DB.Select
	if ctx.Tx != nil {
		if unsafe {
			selfunc = ctx.Tx.Unsafe().Select
		} else {
			selfunc = ctx.Tx.Select
		}
	} else {
		if unsafe {
			selfunc = ctx.DB.Unsafe().Select
		}
	}

	var (
		query string
		args  []interface{}
	)

	qb := utils.QueryBuilderTable(dbtable, sel, qfs...)

	if query, args, err = qb.ToSql(); err != nil {
		// log query
		LogSQL(query, args, err, SQL_CALLER_SKIP)
		return
	} else {
		// log query
		LogSQL(query, args, err, SQL_CALLER_SKIP)
	}

	err = selfunc(target, query, args...)
	if err != nil {
		// log query
		LogSQL(query, args, err, SQL_CALLER_SKIP)
	}
	return
}

func DBFilterCount(ctx *context.Context, dbtable string, paging *paginator.Paginator, qfs ...utils.QueryFunc) (err error) {
	qb := utils.QueryBuilderTable(dbtable, "COUNT(*)", qfs...)
	var (
		query string
		args  []interface{}
	)

	qrxfunc := ctx.DB.QueryRowx
	if ctx.Tx != nil {
		qrxfunc = ctx.Tx.QueryRowx
	}

	query, args, err = qb.ToSql()
	if err != nil {
		// log query
		LogSQL(query, args, err, SQL_CALLER_SKIP)
		return
	}

	// Get count of returned objects
	var count int
	err = qrxfunc(query, args...).Scan(&count)
	// log query
	LogSQL(query, args, err, SQL_CALLER_SKIP)

	if err != nil {
		return
	}
	paging.SetCount(count)
	return
}

/*
Scans single object into target interface
This method should not be used directly from views, rather shiould be used
in manager Filter, FilterAll, Get, GetByID
*/
func DBGet(ctx *context.Context, sel, table string, unsafe bool, target interface{}, qfs ...utils.QueryFunc) (err error) {

	if target == nil {
		return ErrNilPointer
	}

	getfunc := ctx.DB.Get
	if ctx.Tx != nil {
		if unsafe {
			getfunc = ctx.Tx.Unsafe().Get
		} else {
			getfunc = ctx.Tx.Get
		}
	} else {
		if unsafe {
			getfunc = ctx.DB.Unsafe().Get
		}
	}

	getQuery := utils.QueryBuilderTable(table, sel, qfs...)

	var args []interface{}
	var query string

	query, args, err = getQuery.ToSql()
	if err != nil {
		LogSQL(query, args, err, SQL_CALLER_SKIP)
		return
	}

	err = getfunc(target, query, args...)
	LogSQL(query, args, err, SQL_CALLER_SKIP)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return ErrObjectDoesNotExists
		default:
			return err
		}
	}
	return
}

/*
DBInsert - generic insert used in case of modeler
*/
func DBInsert(ctx *context.Context, model Modeler) (err error) {
	handleNilPointer(model)

	// object cannot be inserted, it's already stored in db
	if model.PrimaryKey().Int64() != 0 {
		return ErrObjectAlreadyExists
	}

	qb := utils.QueryBuilder().Insert(model.Table()).
		Columns(model.Columns()...).
		Values(model.Values()...).
		Suffix("RETURNING id")

	var (
		query string
		args  []interface{}
	)
	if query, args, err = qb.ToSql(); err != nil {
		return err
	}

	// change this when #62 will be actual
	var id int64

	qrfunc := ctx.DB.QueryRow
	if ctx.Tx != nil {
		qrfunc = ctx.Tx.QueryRow
	}

	if err = qrfunc(query, args...).Scan(&id); err != nil {
		return
	}

	// set id
	model.SetPrimaryKey(types.PrimaryKey(id))

	return
}

/*
DBUpdate - generic update for modeler
*/

func DBUpdate(ctx *context.Context, model Modeler, fields ...string) (changed bool, err error) {

	handleNilPointer(model)

	if model.PrimaryKey().Int64() == 0 {
		return false, ErrObjectDoesNotExists
	}

	processed := []string{}

	if len(fields) == 0 {
		fields = model.Columns()
	}

	values := model.Values()

	qb := utils.QueryBuilder().Update(model.Table())

	for i, column := range model.Columns() {
		if utils.StringIndex(fields, column) != -1 {
			processed = append(processed, column)
			qb = qb.Set(column, values[i])
		}
	}

	for _, field := range fields {
		if utils.StringIndex(processed, field) == -1 {
			err = fmt.Errorf("field %s not found.", field)
			return
		}
	}

	// add where
	qb = qb.Where("id = ?", model.PrimaryKey())

	var (
		query string
		args  []interface{}
	)

	if query, args, err = qb.ToSql(); err != nil {
		return
	}

	execfunc := ctx.DB.Exec
	if ctx.Tx != nil {
		execfunc = ctx.Tx.Exec
	}

	var result sql.Result
	if result, err = execfunc(query, args...); err != nil {
		return
	}
	affected, _ := result.RowsAffected()
	changed = affected > 0

	return
}

/*
DBDelete - generic delete for Modeler type
*/
func DBDelete(ctx *context.Context, model Modeler) (err error) {
	handleNilPointer(model)
	if model.PrimaryKey().Int64() == 0 {
		return ErrObjectDoesNotExists
	}

	qb := utils.QueryBuilder().Delete(model.Table()).Where("id = ?", model.PrimaryKey())

	var (
		args  []interface{}
		query string
	)

	query, args, err = qb.ToSql()
	if err != nil {
		LogSQL(query, args, err, SQL_CALLER_SKIP)
		return
	}

	execfunc := ctx.DB.Exec
	if ctx.Tx != nil {
		execfunc = ctx.Tx.Exec
	}

	_, err = execfunc(query, args...)
	LogSQL(query, args, err, SQL_CALLER_SKIP)

	model.SetPrimaryKey(types.PrimaryKey(0))

	return
}

/*
Returns list of changed fields between two model instances
*/
func ChangedModelFields(model, another Modeler) (result []string, err error) {
	handleNilPointer(model)
	handleNilPointer(another)

	result = []string{}

	if model.Table() != another.Table() {
		return result, ErrIncorrectModel
	}

	columns := model.Columns()
	modelValues := model.Values()
	anotherValues := another.Values()

	// incorrect something
	if len(modelValues) != len(anotherValues) || len(modelValues) != len(columns) {
		return result, ErrIncorrectModel
	}

	// hacky hack
	if model.PrimaryKey().Int64() != another.PrimaryKey().Int64() {
		result = append(result, "id")
	}

	for i := range modelValues {
		if !reflect.DeepEqual(modelValues[i], anotherValues[i]) {
			result = append(result, columns[i])
		}
	}

	return
}

/*
	Caching
*/
func Cache(context *context.Context, cacheKey string, target interface{}) (err error) {
	var body []byte
	if body, err = json.Marshal(target); err != nil {
		return
	}

	return context.Cache.Set(cacheKey, body)
}

func GetCached(context *context.Context, cacheKey string, target interface{}) (err error) {
	var result []byte
	if result, err = context.Cache.Get(cacheKey); err != nil {
		return
	}
	return json.Unmarshal(result, target)
}

func RemoveCached(context *context.Context, cacheKey string) (err error) {
	return context.Cache.Delete(cacheKey)
}

/*
LogSQL method
*/

func getfilename(path string) (result string) {
	_, result = filepath.Split(path)
	return
}

/*
LogSQL - logs sql query depending on glog.V and arguments given
*/
func LogSQL(query string, args []interface{}, err error, skip int) {

	if err != nil {
		if glog.V(2) {
			if _, file, line, ok := runtime.Caller(skip); ok {
				glog.Infof("%s %s, query: %s, args: %v, file:[%s:%d]",
					alert("sql:"), err.Error(), query, args,
					getfilename(file), line)
			} else {
				glog.Infof("%s %s, query: %s, args: %v",
					alert("sql:"), err.Error(), query, args)
			}
		} else {
			glog.Infof("sql: %s.", err.Error())
		}
		return
	}

	if glog.V(2) {
		if _, file, line, ok := runtime.Caller(skip); ok {
			glog.V(2).Infof("%s %s, args: %v, file: [%s:%d]", phosphorize("sql:"), query, args, getfilename(file), line)
		} else {
			glog.V(2).Infof("%s %s, args: %v", phosphorize("sql:"), query, args)
		}
	}
}
