package entities

import (
	"time"

	"github.com/AubreeH/goApiDb/dataTypes"
)

type UpdatedAt struct {
	UpdatedAt dataTypes.Time `json:"updated_at" db_type:"DATETIME" db_nullable:"false" db_disallow_external_modification:"true" parse_struct:"false" db_default:"CURRENT_TIMESTAMP"`
}

func (val UpdatedAt) OnCreate() (UpdatedAt, error) {
	val.UpdatedAt = dataTypes.Time{Time: time.Now()}
	return val, nil
}

func (val UpdatedAt) OnUpdate() (UpdatedAt, error) {
	val.UpdatedAt = dataTypes.Time{Time: time.Now()}
	return val, nil
}

func (_ UpdatedAt) GetPtrFunc(value *UpdatedAt) *time.Time {
	return value.UpdatedAt.GetPtrFunc(&value.UpdatedAt)
}
