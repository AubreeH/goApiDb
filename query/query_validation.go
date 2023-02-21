package query

func (query *Query) validateAlias(alias string) *Query {
	if alias == "" {
		panic("empty alias provided")
	}

	if _, exists := query.tables[alias]; exists {
		panic("alias already in use")
	}
	return query
}

func (query *Query) validateQuery() *Query {
	if query.selectStr == "" {
		panic("select statement not provided")
	}

	if query.from == "" {
		panic("from table not provided")
	}
	return query
}
