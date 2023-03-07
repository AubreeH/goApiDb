package entities

import (
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type CreatedAt struct {
	CreatedAt dataTypes.Time `json:"created_at"`
}

func (CreatedAt) Describe() {

}

func (_ CreatedAt) GetPtrFunc(value *CreatedAt) *time.Time {
	return value.CreatedAt.GetPtrFunc(&value.CreatedAt)
}
