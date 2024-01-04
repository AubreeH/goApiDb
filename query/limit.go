package query

import "fmt"

type limit struct {
	limit  uint
	offset uint
}

func (q *Query[T]) Limit(limit uint) *Query[T] {
	q.limit.limit = limit
	return q
}

func (q *Query[T]) Offset(offset uint) *Query[T] {
	q.limit.offset = offset
	return q
}

func (l limit) format(pretty bool) string {
	if l.limit == 0 {
		return ""
	}
	if l.offset > 0 {
		return fmt.Sprintf("LIMIT %d, %d", l.offset, l.limit)
	}
	return fmt.Sprintf("LIMIT %d", l.limit)
}
