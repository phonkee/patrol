package utils

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/phonkee/patrol/types"
)

var (
	ErrColumnNotFound = errors.New("column_not_found")
)

/*
DBInfo
*/
func NewDBInfo(db *sqlx.DB) (info *DBInfo, err error) {
	info = &DBInfo{db: db}
	err = info.ReadData()
	return
}

type DBInfo struct {
	db          *sqlx.DB
	Constraints map[string]string
	Columns     []*ColumnInfo
}

// Reads database info
func (d *DBInfo) ReadData() (err error) {
	d.Constraints = make(map[string]string)
	d.Columns = []*ColumnInfo{}

	glog.V(2).Info("DBInfo: reading data.")

	query := `
		SELECT
		    tc.constraint_name,
		    kcu.column_name 
		FROM 
		    information_schema.table_constraints AS tc 
		    JOIN information_schema.key_column_usage AS kcu
		      ON tc.constraint_name = kcu.constraint_name
		    JOIN information_schema.constraint_column_usage AS ccu
		      ON ccu.constraint_name = tc.constraint_name`

	rows, errQuery := d.db.Queryx(query)
	if errQuery != nil {
		return errQuery
	}

	constraint := ""
	column := ""
	for rows.Next() {
		rows.Scan(&constraint, &column)
		d.Constraints[constraint] = column
	}

	d.Columns = []*ColumnInfo{}
	db := d.db.Unsafe()
	err = db.Select(&d.Columns, "select * from information_schema.columns")
	if err != nil {
		return
	}

	return
}

func (d *DBInfo) ColumnInfo(table, column string) (ci *ColumnInfo, err error) {
	for _, c := range d.Columns {
		if c.TableName == table && c.ColumnName == column {
			ci = c
			return
		}
	}
	err = ErrColumnNotFound
	return
}

// @TODO: all columns for table
func (d *DBInfo) TableColumnInfo(table, column string) (ci []*ColumnInfo, err error) {
	return
}

func (d *DBInfo) printColumns() {
	for i, column := range d.Columns {
		r, _ := json.MarshalIndent(column, "", "  ")
		glog.V(2).Infof("column %d is %s", i, string(r))
	}
}

/*
ColumnInfo

can be extended since it uses "unsafe" sqlx db
add additional fields such as character_set_<x>, collation_<x>
*/
type ColumnInfo struct {
	CharacterMaximumLength sql.NullInt64  `db:"character_maximum_length"`
	CharacterOctetLength   sql.NullInt64  `db:"character_octet_length"`
	ColumnDefault          sql.NullString `db:"column_default"`
	ColumnName             string         `db:"column_name"`
	DataType               string         `db:"data_type"`
	DatetimePrecision      sql.NullInt64  `db:"datetime_precision"`
	IsNullable             types.IsField  `db:"is_nullable"`
	IsUpdatable            types.IsField  `db:"is_updatable"`
	NumericPrecision       sql.NullInt64  `db:"numeric_precision"`
	NumericPrecisionRadix  sql.NullInt64  `db:"numeric_precision_radix"`
	NumericScale           sql.NullInt64  `db:"numeric_scale"`
	OrdinalPosition        int            `db:"ordinal_position"`
	TableCatalog           string         `db:"table_catalog"`
	TableName              string         `db:"table_name"`
	TableSchema            string         `db:"table_schema"`
	UDTName                string         `db:"udt_name"`
	UDTCatalog             string         `db:"udt_catalog"`
}
