package query

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/helpers"
)

const (
	Ascending     = "ASC"
	Descending    = "DESC"
	SelectKeyword = "SELECT"
	UpdateKeyword = "UPDATE"
)

type Query struct {
	operation string
	query     string
	selectStr string
	from      string
	joins     []join
	clauses   string
	params    map[string]parameter
	tables    map[string]table
	orderBy   string
	result    *sql.Rows
	Error     error
	args      []any
}

type join struct {
	Type string
	To   string
	On   string
}

func (j join) Format(query *Query) string {
	t := query.tables[j.To]
	return j.Type + " JOIN " + t.Format() + " ON " + j.On
}

type table struct {
	Entity any
	Alias  string
}

func (t table) Format() string {
	return helpers.GetTableName(t.Entity) + " " + t.Alias
}

type parameter struct {
	Name  string
	Value any
}
