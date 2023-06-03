package query

import (
	"regexp"
	"sort"
	"strings"
)

func (query *Query) Build() string {
	switch query.operation {
	case SelectKeyword:
		return query.buildSelect()
	case UpdateKeyword:
		return query.buildUpdate()
	}

	return ""
}

func (query *Query) buildSelect() string {
	query.validateQuery()

	fromTable, err := query.tables[query.from].Format()
	if err != nil {
		query.Error = err
		return ""
	}
	q := "SELECT " + query.selectStr + " "
	q += "FROM " + fromTable + " "

	for i := range query.joins {
		j := query.joins[i]
		join, err := j.Format(query)
		if err != nil {
			query.Error = err
			return ""
		}
		q += join + " "
	}

	if query.clauses != "" {
		q += "WHERE " + query.clauses + " "
	}

	if query.groupBy != "" {
		q += "GROUP BY " + query.groupBy + " "
	}

	if query.orderBy != "" {
		q += "ORDER BY " + query.orderBy + " "
	}

	q, args, err := replaceParams(query.params, q)
	if err != nil {
		query.Error = err
		return ""
	}

	q = strings.Trim(q, " ")

	query.query = q
	query.args = args
	return q
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

func (query *Query) buildUpdate() string {
	panic("not implemented")
}
