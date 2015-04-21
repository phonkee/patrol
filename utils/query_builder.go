/*
Query builder package

Helper functions based on awesome squirrel package.
*/
package utils

import (
	"errors"

	"github.com/lann/squirrel"
)

// query filter func used in Selects
type QueryFunc func(builder squirrel.SelectBuilder) squirrel.SelectBuilder
type QuerySetFunc func(builder squirrel.UpdateBuilder, field string) squirrel.UpdateBuilder

// group of multiple QueryFuncs
func QueryFuncGroup(funcs ...QueryFunc) QueryFunc {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		for _, f := range funcs {
			builder = f(builder)
		}
		return builder
	}
}

var (
	ErrUpdateQueryNoFieldsGiven = errors.New("update_query_no_fields_given")
)

/* Simple shorthand method
returns select query builder with applied query QueryFuncs
*/
func QueryBuilderTable(dbtable, sel string, qffuncs ...QueryFunc) squirrel.SelectBuilder {
	sq := QueryBuilder().Select(sel).From(dbtable)
	sq = ApplyQueryFuncs(sq, qffuncs...)
	return sq
}

// Returns query builder with postgres dialect
func QueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func ApplyQueryFuncs(builder squirrel.SelectBuilder, funcs ...QueryFunc) squirrel.SelectBuilder {
	for _, qff := range funcs {
		builder = qff(builder)
	}
	return builder
}

/*
Simple Where filter
this can be used e.g. like this:
users := []User{}
type := true
err := usermanager.FilterAll(&users, QueryFilterWhere("is_active:", type))
*/
func QueryFilterWhere(pred interface{}, args ...interface{}) QueryFunc {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(pred, args...)
	}
}

func QueryFilterID(id interface{}) QueryFunc {
	return QueryFilterWhere("id = ?", id)
}
