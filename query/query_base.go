package query

import "strings"

func (query *Query) Select(columns ...string) *Query {
	if query.selectStr != "" {
		query.selectStr += ", " + strings.Join(columns, ", ")
	} else {
		query.selectStr = strings.Join(columns, ", ")
	}

	return query
}

func (query *Query) From(entity any, alias string) *Query {
	query.validateAlias(alias)

	query.tables[alias] = table{Entity: entity, Alias: alias}
	query.from = alias

	return query
}

func (query *Query) FromTable(t string, alias string) *Query {
	query.tables[alias] = table{Entity: t, Alias: alias}
	query.from = alias

	return query
}
