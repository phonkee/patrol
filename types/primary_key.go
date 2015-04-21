package types

import (
	"database/sql/driver"
	"strconv"
)

// Primary key used in all models
type PrimaryKey int64

func (p PrimaryKey) Int64() int64 {
	return int64(p)
}

func (p PrimaryKey) String() string {
	return strconv.FormatInt(p.Int64(), 10)
}

func (p PrimaryKey) ToForeignKey() ForeignKey {
	return ForeignKey(p.Int64())
}

func (p *PrimaryKey) Parse(value string) error {
	if id, err := strconv.ParseInt(value, 10, 0); err != nil {
		return err
	} else {
		*p = PrimaryKey(id)
	}
	return nil
}

func (p *PrimaryKey) Scan(value interface{}) error {
	*p = PrimaryKey(value.(int64))
	return nil
}
func (p PrimaryKey) Value() (driver.Value, error) {
	return p.Int64(), nil
}
