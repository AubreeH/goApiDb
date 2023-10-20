package newQuery

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
	out = append(out, q.formatSelect(pretty))
	out = append(out, q.from.format(pretty))
	out = append(out, q.joins.format(pretty))
	out = append(out, "WHERE "+q.where.format(pretty))
	out = append(out, q.groupBy.format(pretty))
	out = append(out, q.orderBy.format(pretty))
	out = append(out, q.limit.format(pretty))

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
