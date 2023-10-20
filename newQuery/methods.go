package newQuery

import (
	"fmt"
	"reflect"

	"github.com/AubreeH/goApiDb/database"
)

func (q *query[T]) Sql() string {
	return q.format(false)
}

func (q *query[T]) First(db database.Database) (T, error) {
	var out T
	lim := q.limit.limit
	q.limit.limit = 1

	parsedQuery, queryArgs := q.params.parse(q.format(false))

	rows, err := db.Db.Query(parsedQuery, queryArgs...)
	if err != nil {
		return out, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return out, err
	}

	for rows.Next() {
		row, rowArgs := getScannableRow(out, columns)
		err := rows.Scan(rowArgs...)
		if err != nil {
			return out, err
		}

		out = row
	}

	q.limit.limit = lim

	return out, err
}

func (q *query[T]) FirstN(db *database.Database) ([]T, error) {
	var out []T

	return out, nil
}

func (q *query[T]) All(db *database.Database) ([]T, error) {
	var t T
	var out []T

	parsedQuery, queryArgs := q.params.parse(q.format(false))
	fmt.Println(parsedQuery)

	rows, err := db.Db.Query(parsedQuery, queryArgs...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return out, err
	}

	for rows.Next() {
		row, rowArgs := getScannableRow(t, columns)
		err := rows.Scan(rowArgs...)
		if err != nil {
			return out, err
		}

		out = append(out, row)
	}

	return out, nil
}

func getScannableRow[T any](s T, sqlRowsColumns []string) (T, []interface{}) {
	var out T

	var args = make([]interface{}, len(sqlRowsColumns))

	ptrs := getPointers(&out)

	for i, col := range sqlRowsColumns {
		args[i] = ptrs[col]
	}

	return out, args
}

func getPointers[T any](s *T) map[string]interface{} {
	out := make(map[string]interface{})

	refVal := reflect.ValueOf(s).Elem()

	for i := 0; i < refVal.NumField(); i++ {
		field := refVal.Type().Field(i)
		if field.Anonymous {
			continue
		}

		colName := field.Name
		if colName == "" {
			continue
		}

		out[colName] = refVal.Field(i).Addr().Interface()
	}

	return out
}
