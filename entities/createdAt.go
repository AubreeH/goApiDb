package entities

import (
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type CreatedAt struct {
	CreatedAt dataTypes.Time `json:"created_at" db_type:"DATETIME" db_disallow_external_modification:"true" parse_struct:"false"`
}

func (val CreatedAt) OnCreate() (CreatedAt, error) {
	val.CreatedAt = dataTypes.Time{Time: time.Now()}
	return val, nil
}

func (CreatedAt) Describe() {

}

func (_ CreatedAt) GetPtrFunc(value *CreatedAt) *time.Time {
	return value.CreatedAt.GetPtrFunc(&value.CreatedAt)
}
