package query

func (query *Query) Limit(value uint) *Query {
	query.limit = value
	return query
}

func (query *Query) Offset(value uint) *Query {
	query.offset = value
	return query
}

func (query *Query) Paginated(itemsPerPage, offset uint) *Query {
	query.limit = itemsPerPage
	query.offset = offset
	return query
}
