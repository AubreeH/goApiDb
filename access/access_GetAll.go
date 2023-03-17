package access

import (
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
	"reflect"
)

func GetAll[T any](db *database.Database, entity T, limit int) ([]T, error) {
	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return nil, err
	}

	var args []any

	query := "SELECT * FROM " + tableInfo.Name

	if tableInfo.SoftDeletes != "" {
		query += " WHERE " + tableInfo.SoftDeletes + " IS NULL"
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
