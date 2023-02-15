package query

import (
	"database/sql"
	"goApiDb/database"
	"goApiDb/helpers"
	"reflect"
	"regexp"
	"sort"
	"strings"
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

// -------------------------------
// NEW QUERY FUNCTIONS
// -------------------------------

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

// -------------------------------
// VALIDATION FUNCTIONS
// -------------------------------

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

// -------------------------------
// SELECT FUNCTIONS
// -------------------------------

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

// -------------------------------
// JOIN FUNCTIONS
// -------------------------------

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

// -------------------------------
// WHERE FUNCTIONS
// -------------------------------

func (query *Query) Where(clauses ...string) *Query {
	query.clauses = strings.Join(clauses, " AND ")

	return query
}

func (query *Query) AndWhere(clause string) *Query {
	query.clauses += " AND " + clause

	return query
}

// -------------------------------
// ORDER BY FUNCTIONS
// -------------------------------

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

// -------------------------------
// PARAMETER FUNCTIONS
// -------------------------------

func (query *Query) SetParam(name string, param any) *Query {
	query.params[name] = parameter{Name: name, Value: param}
	return query
}

// -------------------------------
// BUILD FUNCTIONS
// -------------------------------

func (query *Query) Build() {
	switch query.operation {
	case SelectKeyword:
		query.buildSelect()
	case UpdateKeyword:
		query.buildUpdate()
	}
}

func (query *Query) buildSelect() {
	query.validateQuery()

	q := "SELECT " + query.selectStr + " "
	q += "FROM " + query.tables[query.from].Format() + " "

	for i := range query.joins {
		j := query.joins[i]
		q += j.Format(query) + " "
	}

	if query.clauses != "" {
		q += "WHERE " + query.clauses + " "
	}

	if query.orderBy != "" {
		q += "ORDER BY " + query.orderBy
	}

	q, args, err := replaceParams(query.params, q)
	if err != nil {
		query.Error = err
		return
	}

	query.query = q
	query.args = args
}

func replaceParams(parameters map[string]parameter, q string) (string, []any, error) {
	sqlParams := make(map[int]any)
	for s := range parameters {
		param := parameters[s]
		re, err := regexp.Compile(":" + param.Name)
		if err != nil {
			return "", nil, err
		}
		result := re.FindAllIndex([]byte(q), -1)
		for i := range result {
			sqlParams[result[i][0]] = param.Value
		}
		q = strings.Replace(q, ":"+s, "?", -1)
	}

	keys := make([]int, 0, len(sqlParams))
	for i := range sqlParams {
		keys = append(keys, i)
	}
	sort.Ints(keys)

	args := make([]any, 0, len(sqlParams))
	for _, i := range keys {
		args = append(args, sqlParams[i])
	}

	return q, args, nil
}

func (query *Query) buildUpdate() {
	panic("not implemented")
}

func (query *Query) Exec() {
	query.Error = nil

	query.Build()
	if query.Error != nil {
		return
	}

	db := database.GetDb()
	result, err := db.Query(query.query, query.args...)

	query.result = result
	query.Error = err
}

func ExecuteQuery[T any](query *Query, _ T) ([]T, error) {
	query.Exec()
	if query.Error != nil {
		return nil, query.Error
	}

	var output []T

	result := query.result
	for result.Next() {
		var row T
		args := GetRowArgs(&row)
		err := result.Scan(args...)
		if err != nil {
			return nil, err
		}
		output = append(output, row)
	}

	return output, nil
}

func GetRowArgs[T any](row *T) []any {
	var args []any

	refValue := reflect.ValueOf(row).Elem()
	for i := 0; i < refValue.NumField(); i++ {
		args = append(args, refValue.Field(i).Addr().Interface())
	}

	return args
}

func GetResult(query *Query) *sql.Rows {
	return query.result
}
