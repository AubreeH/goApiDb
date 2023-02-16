package baseEntities

import (
	"github.com/AubreeH/goApiDb/dataTypes"
	"time"
)

type CreatedAt struct {
	CreatedAt dataTypes.Time `json:"created_at" sql_name:"created_at" sql_type:"DATETIME" sql_nullable:"false" parse_struct:"false" sql_disallow_external_modification:"true"`
}

func (val CreatedAt) OnCreate() (CreatedAt, error) {
	val.CreatedAt = dataTypes.Time{Time: time.Now()}
	return val, nil
}
