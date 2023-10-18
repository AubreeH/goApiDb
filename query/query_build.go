package query

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func (query *Query) Build() (baseQuery, paginationDetailsQuery string, baseQueryParams, paginationDetailsQueryParams []any, err error) {
	switch query.operation {
	case SelectKeyword:
		return query.buildSelect()
	case UpdateKeyword:
		return query.buildUpdate(), "", nil, nil, nil
	}

	return "", "", nil, nil, errors.New("operation not supported")
}

func (query *Query) buildSelect() (string, string, []any, []any, error) {
	err := query.validateQuery()
	if err != nil {
		query.Error = err
		return "", "", nil, nil, err
	}

	fromTable, err := query.tables[query.from].Format()
	if err != nil {
		query.Error = err
		return "", "", nil, nil, err
	}

	q := ""

	for i := range query.joins {
		j := query.joins[i]
		join, err := j.Format(query)
		if err != nil {
			query.Error = err
			return "", "", nil, nil, err
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

	limitStatement := ""
	if query.limit != 0 {
		if query.offset == 0 {
			limitStatement = fmt.Sprintf(" LIMIT %d ", query.limit)
		} else {
			limitStatement = fmt.Sprintf(" LIMIT %d OFFSET %d ", query.limit, query.limit*query.offset)
		}
	}

	q1 := fmt.Sprintf("SELECT %s FROM %s %s%s", query.selectStr, fromTable, q, limitStatement)
	q2 := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", fromTable, q)

	q1, q1Args, err := replaceParams(query.params, q1)
	if err != nil {
		query.Error = err
		return "", "", nil, nil, err
	}

	q2, q2Args, err := replaceParams(query.params, q2)
	if err != nil {
		query.Error = err
		return "", "", nil, nil, err
	}

	q1 = strings.Trim(q1, " ")

	query.query = q1
	query.args = q1Args

	query.paginationDetailsQuery = q2
	query.paginationDetailsQueryArgs = q2Args
	return q1, q2, q1Args, q2Args, nil
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
