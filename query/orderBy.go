package query

import "strings"

type orderBy struct {
	columns []string
}

func (q *query[T]) OrderBy(columns ...string) *query[T] {
	q.orderBy.columns = columns
	return q
}

func (o orderBy) format(pretty bool) string {
	if o.columns == nil || len(o.columns) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(o.columns, ", ")
}
