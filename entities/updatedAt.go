package baseEntities

import (
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type UpdatedAt struct {
	UpdatedAt dataTypes.Time `json:"updated_at" sql_name:"updated_at" sql_type:"DATETIME" sql_nullable:"false" parse_struct:"false" sql_disallow_external_modification:"true"`
}

func (val UpdatedAt) OnCreate() (UpdatedAt, error) {
	val.UpdatedAt = dataTypes.Time{Time: time.Now()}
	return val, nil
}

func (val UpdatedAt) OnUpdate() (UpdatedAt, error) {
	val.UpdatedAt = dataTypes.Time{Time: time.Now()}
	return val, nil
}
