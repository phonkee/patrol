package models

import (
	"time"

	"github.com/phonkee/patrol/context"
)

/* BufferedCounter

BufferedCounter is counter value which is written to database. Values are buffered
in cache first and after some time they will be written to database

*/

type BufferedCounter struct {
	Model
	Identifier string    `db:"identifier" json:"identifier"`
	DateAdded  time.Time `db:"date_added" json:"date_added"`
}

// Buffered counter manager
func NewBufferedCounterManager(context *context.Context) *BufferedCounterManager {
	return &BufferedCounterManager{
		context: context,
	}
}

type BufferedCounterManager struct {
	Manager
	context *context.Context
}
