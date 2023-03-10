package entities

import (
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type UpdatedAt struct {
	UpdatedAt dataTypes.Time `json:"updated_at" sql_type:"DATETIME" sql_nullable:"false" sql_disallow_external_modification:"true" parse_struct:"false"`
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
