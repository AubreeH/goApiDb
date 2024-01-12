package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/AubreeH/goApiDb/structParsing"
)

type Query[T any] struct {
	operation operation
	distinct  bool
	where     Where
	from      from
	joins     joins
	groupBy   groupBy
	orderBy   orderBy
	limit     limit
	params    params

	QuerySettings
}

type QueryResult[T any] struct {
	Results    []T
	Total      uint
	Paginated  bool
	Page       uint
	TotalPages uint
	Query      *Query[T]
}

func Select[T any](s T) *Query[T] {
	return &Query[T]{
		operation: selectOperation,
	}
}

func SelectDistinct[T any](s T) *Query[T] {
	return &Query[T]{
		operation: selectOperation,
		distinct:  true,
	}
}

func (q *Query[T]) format(pretty bool) (string, error) {
	var out []string
	if selectStatement, err := q.formatSelect(pretty); err != nil {
		return "", err
	} else if selectStatement == "" {
		return "", fmt.Errorf("no select statement found")
	} else {
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
		return strings.Join(out, "\n"), nil
	} else {
		return strings.Join(out, " "), nil
	}
}

func (q *Query[T]) formatSelect(pretty bool) (string, error) {
	var t T
	refType := reflect.TypeOf(t)

	columns, err := q.parseSelectColumns(refType)
	if err != nil {
		return "", err
	}

	var formattedColumns string
	if pretty {
		formattedColumns = strings.Join(columns, ",\n")
	} else {
		formattedColumns = strings.Join(columns, ", ")
	}

	if q.distinct {
		return fmt.Sprintf("SELECT DISTINCT %s", formattedColumns), nil
	} else {
		return fmt.Sprintf("SELECT %s", formattedColumns), nil
	}
}

func (q *Query[T]) parseSelectColumns(refType reflect.Type) ([]string, error) {
	if refType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("provided type in query is not a struct")
	}

	var columns []string

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)

		if structParsing.FormatSqlIgnore(field) {
			continue
		}

		if val, ok := structParsing.QueryIgnore.Lookup(field); ok && structParsing.FormatBoolean(val) == 1 {
			continue
		}

		if field.Type.Kind() == reflect.Struct && structParsing.FormatParseStruct(field) {
			newColumns, err := q.parseSelectColumns(field.Type)
			if err != nil {
				return nil, err
			}

			columns = append(columns, newColumns...)
		} else if val, ok := structParsing.QuerySelect.Lookup(field); !ok || val == "" {
			if q.preventFieldNameAutoMapping {
				continue
			}
			columns = append(columns, q.from.alias+"."+structParsing.FormatSqlName(field))
		} else {
			columns = append(columns, val+" as "+field.Name)
		}
	}

	return columns, nil
}

func (q *Query[T]) formatWhere(pretty bool) string {
	val := q.where.format(pretty)
	if val == "" {
		return ""
	}
	return "WHERE " + val
}
