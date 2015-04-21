package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

type NullInt struct {
	sql.NullInt64
}

func (ns *NullInt) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		ns.Int64 = 0
		ns.Valid = false
		return
	}
	err = json.Unmarshal(b, &ns.Int64)
	if err != nil {
		return
	}
	ns.Valid = true
	return
}

func (ns NullInt) MarshalJSON() (b []byte, err error) {
	if ns.Int64 == 0 && !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Int64)
}
