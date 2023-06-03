package query

import (
	"strings"
)

func (query *Query) GroupBy(value ...string) *Query {
	query.groupBy = strings.Join(value, ", ")

	return query
}

func (query *Query) AddGroupBy(value ...string) *Query {
	if query.groupBy == "" {
		return query.GroupBy(value...)
	}

	query.groupBy += ", " + strings.Join(value, ", ")

	return query
}
