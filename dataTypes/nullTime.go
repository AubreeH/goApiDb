package dataTypes

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (value NullTime) MarshalJSON() ([]byte, error) {
	if value.Valid {
		return json.Marshal(value.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (value NullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		value.Valid = true
		value.Time = *x
	} else {
		value.Valid = false
	}
	return nil
}

func (value NullTime) ExtractDataFunc() any {
	if value.Valid {
		return value.Time
	}

	return nil
}

func (_ NullTime) GetPtrFunc(value *NullTime) *sql.NullTime {
	return &value.NullTime
}
