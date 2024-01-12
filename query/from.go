package query

import (
	"fmt"

	"github.com/AubreeH/goApiDb/structParsing"
)

type from struct {
	entity interface{}
	alias  string
}

func (q *Query[T]) From(entity interface{}, alias string) *Query[T] {
	q.from = from{
		entity: entity,
		alias:  alias,
	}
	return q
}

func (f from) format(pretty bool) string {
	tableInfo, _ := structParsing.GetTableInfo(f.entity)
	return fmt.Sprintf("FROM `%s` %s", tableInfo.Name, f.alias)
}
