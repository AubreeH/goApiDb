package access

import (
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/helpers"
	"reflect"
)

func GetAll[T any](db *database.Database, entity T, limit int) ([]T, error) {
	tableName := helpers.GetTableName(entity)

	var args []any

	query := "SELECT * FROM " + tableName

	if DoesEntitySoftDelete(entity) {
		query += " WHERE deleted = false"
	}

	if limit != 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	result, err := db.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var retValue []T
	args, entityOutput, err := database.BuildRow(db, entity, result)
	for result.Next() {
		if err != nil {
			return nil, err
		}
		err = result.Scan(args...)
		if entityOutput != nil {
			retValue = append(retValue, reflect.ValueOf(entityOutput).Elem().Interface().(T))
		}
	}

	return retValue, nil
}
