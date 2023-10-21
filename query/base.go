package query

import (
	"fmt"
	"reflect"
	"strings"
)

type query[T any] struct {
	operation operation
	distinct  bool
	where     Where
	from      from
	joins     joins
	groupBy   groupBy
	orderBy   orderBy
	limit     limit
	params    params
}

type QueryResult[T any] struct {
	Results    []T
	Total      uint
	Paginated  bool
	Page       uint
	TotalPages uint
	Query      *query[T]
}

func Select[T any](s T) *query[T] {
	return &query[T]{
		operation: selectOperation,
	}
}

func SelectDistinct[T any](s T) *query[T] {
	return &query[T]{
		operation: selectOperation,
		distinct:  true,
	}
}

func (q *query[T]) format(pretty bool) string {

	var out []string
	if selectStatement := q.formatSelect(pretty); selectStatement != "" {
		out = append(out, selectStatement)
	}
	if fromStatement := q.from.format(pretty); fromStatement != "" {
		out = append(out, fromStatement)
	}
	if joinsStatement := q.joins.format(pretty); joinsStatement != "" {
		out = append(out, joinsStatement)
	}
	if whereStatement := q.formatWhere(pretty); whereStatement != "" {
		out = append(out, whereStatement)
	}
	if groupByStatement := q.groupBy.format(pretty); groupByStatement != "" {
		out = append(out, groupByStatement)
	}
	if orderByStatement := q.orderBy.format(pretty); orderByStatement != "" {
		out = append(out, orderByStatement)
	}
	if limitStatement := q.limit.format(pretty); limitStatement != "" {
		out = append(out, limitStatement)
	}

	if pretty {
		return strings.Join(out, "\n")
	} else {
		return strings.Join(out, " ")
	}
}

func (q *query[T]) formatSelect(pretty bool) string {
	var t T
	refType := reflect.TypeOf(t)

	var columns []string
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		columns = append(columns, field.Tag.Get("s")+" as "+field.Name)
	}

	var formattedColumns string
	if pretty {
		formattedColumns = strings.Join(columns, ",\n")
	} else {
		formattedColumns = strings.Join(columns, ", ")
	}

	if q.distinct {
		return fmt.Sprintf("SELECT DISTINCT %s", formattedColumns)
	} else {
		return fmt.Sprintf("SELECT %s", formattedColumns)
	}
}

func (q *query[T]) formatWhere(pretty bool) string {
	val := q.where.format(pretty)
	if val == "" {
		return ""
	}
	return "WHERE " + val
}
