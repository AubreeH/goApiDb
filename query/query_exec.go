package query

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/database"
	"reflect"
)

func (query *Query) Exec(db *database.Database) {
	query.Error = nil

	query.Build()
	if query.Error != nil {
		return
	}

	result, err := db.Db.Query(query.query, query.args...)

	query.result = result
	query.Error = err
}

func ExecuteQuery[T any](db *database.Database, query *Query, _ T) ([]T, error) {
	query.Exec(db)
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
