package entities

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type DeletedAt struct {
	DeletedAt dataTypes.NullTime `json:"deleted_at" db_type:"DATETIME" db_nullable:"true" db_disallow_external_modification:"true" parse_struct:"false" soft_deletes:"true"`
}

func (val DeletedAt) OnDelete() (DeletedAt, error) {
	val.DeletedAt = dataTypes.NullTime{NullTime: sql.NullTime{Time: time.Now(), Valid: true}}
	return val, nil
}

func (_ DeletedAt) GetPtrFunc(value *DeletedAt) any {
	return value.DeletedAt.GetPtrFunc(&value.DeletedAt)
}
