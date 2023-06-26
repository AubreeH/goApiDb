package query

import (
	"database/sql"
	"reflect"

	"github.com/AubreeH/goApiDb/database"
)

func (query *Query) exec(db *database.Database) error {
	query.Error = nil

	_, paginationDetailsQuery, _, paginationDetailsQueryParams, err := query.Build()
	if err != nil {
		query.Error = err
		return err
	}
	if query.Error != nil {
		return query.Error
	}

	if paginationDetailsQuery != "" && paginationDetailsQueryParams != nil {
		pdqResults, err := db.Db.Query(paginationDetailsQuery, paginationDetailsQueryParams...)
		if err != nil {
			query.Error = err
			return err
		}
		query.paginationDetailsQueryResult = pdqResults
	}

	result, err := db.Db.Query(query.query, query.args...)
	if err != nil {
		query.Error = err
		return err
	}
	query.result = result

	return nil
}

func ExecuteQuery[T any](db *database.Database, query *Query, _ T) ([]T, error) {
	query.exec(db)
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
