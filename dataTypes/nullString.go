package dataTypes

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (value NullString) MarshalJSON() ([]byte, error) {
	if value.Valid {
		return json.Marshal(value.String)
	} else {
		return json.Marshal(nil)
	}
}

func (value NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		value.Valid = true
		value.String = *x
	} else {
		value.Valid = false
	}
	return nil
}

func (value NullString) ExtractDataFunc() any {
	if value.Valid {
		return value.String
	}

	return nil
}

func (_ NullString) GetPtrFunc(value *NullString) *sql.NullString {
	return &value.NullString
}
