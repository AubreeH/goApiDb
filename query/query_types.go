package query

import (
	"database/sql"

	"github.com/AubreeH/goApiDb/structParsing"
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
	groupBy   string
	result    *sql.Rows
	Error     error
	args      []any
}

type join struct {
	Type string
	To   string
	On   string
}

func (j join) Format(query *Query) (string, error) {
	t := query.tables[j.To]
	tbl, err := t.Format()
	if err != nil {
		return "", err
	}
	return j.Type + " JOIN " + tbl + " ON " + j.On, nil
}

type table struct {
	Entity any
	Alias  string
}

func (t table) Format() (string, error) {
	tblInfo, err := structParsing.GetTableInfo(t.Entity)
	if err != nil {
		return "", err
	}

	return tblInfo.Name + " " + t.Alias, nil
}

type parameter struct {
	Name  string
	Value any
}
