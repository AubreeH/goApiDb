package access

import (
	"errors"
	"fmt"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
	"reflect"
)

func GetById[T any](db *database.Database, entity T, id any) (T, error) {
	var output T

	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return output, err
	}

	var query string
	if tableInfo.SoftDeletes != "" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s IS NULL AND id = ? LIMIT 1", tableInfo.Name, tableInfo.SoftDeletes)
	} else {
		query = fmt.Sprintf("SELECT *  FROM %s WHERE id = ? LIMIT 1", tableInfo.Name)
	}

	result, err := db.Db.Query(query, id)
	if err != nil {
		return entity, err
	}

	args, entityOutput, err := database.BuildRow(db, entity, result)
	if !result.Next() {
		return entity, errors.New("unable to find value")
	}

	err = result.Scan(args...)
	if err != nil {
		return entity, err
	}

	entity = reflect.ValueOf(entityOutput).Elem().Interface().(T)

	return entity, nil
}
