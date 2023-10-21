package query

import "strings"

type groupBy struct {
	columns []string
}

func (q *query[T]) GroupBy(columns ...string) *query[T] {
	q.groupBy.columns = columns
	return q
}

func (g groupBy) format(pretty bool) string {
	if g.columns == nil || len(g.columns) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(g.columns, ", ")
}
