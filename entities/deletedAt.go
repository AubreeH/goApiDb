package entities

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type DeletedAt struct {
	DeletedAt dataTypes.NullTime `json:"deleted_at" sql_name:"deleted_at" sql_type:"DATETIME" sql_nullable:"true" parse_struct:"false" sql_disallow_external_modification:"true"`
	SoftDeletes
}

func (val DeletedAt) OnDelete() (DeletedAt, error) {
	val.DeletedAt = dataTypes.NullTime{NullTime: sql.NullTime{Time: time.Now(), Valid: true}}
	return val, nil
}

func (_ DeletedAt) GetPtrFunc(value *DeletedAt) []any {
	return []any{
		value.DeletedAt.GetPtrFunc(&value.DeletedAt),
		value.SoftDeletes.GetPtrFunc(&value.SoftDeletes),
	}
}
