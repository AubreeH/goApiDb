package query

import "strings"

func (query *Query) Where(clauses ...string) *Query {
	query.clauses = strings.Join(clauses, " AND ")

	return query
}

func (query *Query) AndWhere(clause string) *Query {
	query.clauses += " AND " + clause

	return query
}
