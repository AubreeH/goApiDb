package newQuery

type limit struct {
	limit  uint
	offset uint
}

func (q *query[T]) Limit(limit uint) *query[T] {
	q.limit.limit = limit
	return q
}

func (q *query[T]) Offset(offset uint) *query[T] {
	q.limit.offset = offset
	return q
}

func (q *query[T]) Paginated(pageSize, pageNumber uint) *query[T] {
	q.limit.limit = pageSize
	q.limit.offset = pageSize * (pageNumber - 1)
	return q
}

func (l limit) format(pretty bool) string {
	if l.limit == 0 {
		return ""
	}
	if l.offset > 0 {
		return "LIMIT " + string(l.offset) + ", " + string(l.limit)
	}
	return "LIMIT " + string(l.limit)
}
