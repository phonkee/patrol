package types

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

type Keyer interface {
	Int64() int64

	String() string
}

// FOreignKey used in all models
type ForeignKey int64

/*
provide interface for all foreignKeys and don't need to change all code everywhere
in case of move from ForeignKey to NullForeihnKey We provide method Set()
this method the value is defined as interface{} so we can accept nil's.
*/
func (f ForeignKey) Int64() int64 {
	return int64(f)
}

func (f ForeignKey) String() string {
	return strconv.FormatInt(f.Int64(), 10)
}

func (f ForeignKey) ToPrimaryKey() PrimaryKey {
	return PrimaryKey(f.Int64())
}

func (f *ForeignKey) Scan(value interface{}) error {
	*f = ForeignKey(value.(int64))
	return nil
}
func (f ForeignKey) Value() (driver.Value, error) {
	return f.Int64(), nil
}

/*
NullFOreignKey - nullable foreign key
*/
type NullForeignKey struct {
	sql.NullInt64
}

func (n NullForeignKey) String() string {
	return strconv.FormatInt(n.Int64, 10)
}
