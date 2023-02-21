package query

func (query *Query) OrderBy(column string, direction string) *Query {
	query.orderBy = column + " " + direction

	return query
}

func (query *Query) OrderByAscending(column string) *Query {
	return query.OrderBy(column, Ascending)
}

func (query *Query) OrderByDescending(column string) *Query {
	return query.OrderBy(column, Descending)
}
