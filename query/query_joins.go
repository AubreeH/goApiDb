package query

func (query *Query) Join(joinType string, entity any, alias string, on string) *Query {
	query.validateAlias(alias)

	query.tables[alias] = table{Entity: entity, Alias: alias}
	query.joins = append(query.joins, join{
		Type: joinType,
		To:   alias,
		On:   on,
	})

	return query
}

func (query *Query) LeftJoin(entity any, alias string, on string) *Query {
	return query.Join("LEFT", entity, alias, on)
}

func (query *Query) InnerJoin(entity any, alias string, on string) *Query {
	return query.Join("INNER", entity, alias, on)
}
