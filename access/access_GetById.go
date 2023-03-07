package access

import (
	"errors"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/entities"
	"reflect"
)

func GetById[T any](db *database.Database, entity T, id any) (T, error) {
	var output T

	tableInfo, err := entities.GetTableInfo(entity)
	if err != nil {
		return output, err
	}

	var query string
	if DoesEntitySoftDelete(entity) {
		query = "SELECT *" + " FROM " + tableInfo.Name + " WHERE deleted = false AND id = ? LIMIT 1"
	} else {
		query = "SELECT *" + " FROM " + tableInfo.Name + " WHERE id = ? LIMIT 1"
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
