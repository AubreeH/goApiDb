package query

func NewSelectQuery() *Query {
	query := Query{}
	query.operation = SelectKeyword
	query.tables = make(map[string]table)
	query.params = make(map[string]parameter)
	return &query
}

func NewUpdateQuery() *Query {
	query := Query{}
	query.operation = UpdateKeyword
	query.tables = make(map[string]table)
	query.params = make(map[string]parameter)
	return &query
}
