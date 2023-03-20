package entities

import (
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type CreatedAt struct {
	CreatedAt dataTypes.Time `json:"created_at" db_type:"DATETIME" db_disallow_external_modification:"true" parse_struct:"false"`
}

func (CreatedAt) Describe() {

}

func (_ CreatedAt) GetPtrFunc(value *CreatedAt) *time.Time {
	return value.CreatedAt.GetPtrFunc(&value.CreatedAt)
}
