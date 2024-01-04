package query

import (
	"fmt"
	"strings"

	"github.com/AubreeH/goApiDb/structParsing"
)

type joinType string

const (
	innerJoin joinType = "INNER"
	leftJoin  joinType = "LEFT"
	rightJoin joinType = "RIGHT"
)

type join struct {
	joinType joinType
	entity   interface{}
	alias    string
	on       string
}

type joins []join

func (j joins) join(joinType joinType, entity interface{}, alias, on string) joins {
	return append(j, join{
		joinType: joinType,
		entity:   entity,
		alias:    alias,
		on:       on,
	})
}

func (q Query[T]) LeftJoin(entity interface{}, alias, on string) *Query[T] {
	q.joins = q.joins.join(leftJoin, entity, alias, on)
	return &q
}

func (q Query[T]) RightJoin(entity interface{}, alias, on string) *Query[T] {
	q.joins = q.joins.join(rightJoin, entity, alias, on)
	return &q
}

func (q Query[T]) InnerJoin(entity interface{}, alias, on string) *Query[T] {
	q.joins = q.joins.join(innerJoin, entity, alias, on)
	return &q
}

func (j *join) format(pretty bool) string {
	tableInfo, _ := structParsing.GetTableInfo(j.entity)
	return fmt.Sprintf("%s JOIN `%s` %s ON %s", string(j.joinType), tableInfo.Name, j.alias, j.on)
}

func (j joins) format(pretty bool) string {
	var out []string
	for _, join := range j {
		out = append(out, join.format(pretty))
	}
	if pretty {
		return strings.Join(out, "\n")
	}
	return strings.Join(out, " ")
}
