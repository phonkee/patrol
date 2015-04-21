package models

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/types"
)

/*
Modeler interface

removes clutter of CRUD methods.
*/
type Modeler interface {
	// returns list of all columns
	Columns() []string

	// returns primary key
	PrimaryKey() types.PrimaryKey

	// sets primary key
	SetPrimaryKey(types.Keyer)

	// returns string representation (identifier for e.g. caching)
	String() string

	// returns database table
	Table() string

	// returns list of all values for given
	Values() []interface{}

	// CRUD
	// still not used but will be in the future
	Insert(ctx *context.Context) (err error)
	Update(ctx *context.Context, fields ...string) (changed bool, err error)
	Delete(ctx *context.Context) (err error)
}

/*
Basic model with primary key
*/
type Model struct {
	ID types.PrimaryKey `db:"id" json:"id"`
}

// Returns primary key value
func (m *Model) PrimaryKey() types.PrimaryKey {
	return m.ID
}

// Sets primary key
func (m *Model) SetPrimaryKey(k types.Keyer) {
	m.ID = types.PrimaryKey(k.Int64())
}
