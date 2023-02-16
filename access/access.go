package access

import (
	"errors"
	"fmt"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/helpers"
	"reflect"
)

func GetById[T any](db *database.Database, entity T, id any) (T, error) {
	tableName := helpers.GetTableName(entity)

	var query string
	if DoesEntitySoftDelete(entity) {
		query = "SELECT *" + " FROM " + tableName + " WHERE deleted = false AND id = ? LIMIT 1"
	} else {
		query = "SELECT *" + " FROM " + tableName + " WHERE id = ? LIMIT 1"
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

func Create[T any](db *database.Database, values []T) (T, error) {
	var entity T
	tableName := helpers.GetTableName(entity)

	queryColumns := ""
	queryValues := ""
	var args []any

	var output T
	var id any
	for i := range values {
		rowData, err := GetData(values[i], CreateOperationHandler)
		if err != nil {
			return output, err
		}

		for j := range rowData {
			columnData := rowData[j]

			if queryColumns == "" {
				queryColumns += columnData.Column
				queryValues += "?"
			} else {
				queryColumns += ", " + columnData.Column
				queryValues += ", ?"
			}

			if columnData.PrimaryKey {
				id = columnData.Data
			}

			args = append(args, columnData.Data)
		}
	}

	query := fmt.Sprintf("INSERT"+" INTO %s (%s) values (%s)", tableName, queryColumns, queryValues)

	res, err := db.Db.Exec(query, args...)
	if err != nil {
		return output, err
	}

	if !(id == nil || id == 0 || id == "") {
		return GetById(db, output, id)
	}

	intId, err := res.LastInsertId()
	if err != nil {
		return output, err
	}

	newEntity, err := GetById(db, output, intId)

	return newEntity, err
}

func Update[T any](db *database.Database, value T, id any) error {
	return update(db, value, id, UpdateOperationHandler)
}

func Delete[T any](db *database.Database, entity T, id any) error {
	var err error
	entity, err = GetById(db, entity, id)
	if err != nil {
		return err
	}

	tableName := helpers.GetTableName(entity)

	if DoesEntitySoftDelete(entity) {
		return softDelete(db, entity, id)
	} else {
		_, err = DeleteOperationHandler(reflect.ValueOf(entity))
		if err != nil {
			return err
		}

		q := "DELETE FROM " + tableName + " WHERE ID = ?"

		_, err = db.Db.Exec(q, id)
		return err
	}
}

func update[T any](db *database.Database, value T, id any, operationHandler OperationHandler) error {
	var output T

	existingValue, err := GetById(db, output, id)
	if err != nil {
		return err
	}

	mergedValue := MergeObjects(existingValue, value)
	idColumn, mergedData, err := GetDataAndId(mergedValue, operationHandler)
	if err != nil {
		return err
	}
	_, updateData, err := GetDataAndId(value, NonOperationHandler)
	if err != nil {
		return err
	}

	tableName := helpers.GetTableName(value)

	q := "UPDATE " + tableName + " t SET "
	qBase := q
	var where string

	if DoesEntitySoftDelete(value) {
		where = " WHERE deleted = false AND t."
	} else {
		where = " WHERE t."
	}

	var args []any
	for j := range mergedData {
		column := mergedData[j]
		if updateData[j].Data != column.Data {
			if q != qBase {
				q += ", "
			}
			q += "t." + column.Column + " = ?"
			args = append(args, column.Data)
		}
	}
	args = append(args, id)
	where += idColumn.Column + " = ?"
	_, err = db.Db.Exec(q+where, args...)
	if err != nil {
		return err
	}

	return nil
}

func softDelete[T any](db *database.Database, entity T, id any) error {
	existingEntity, err := GetById(db, entity, id)
	if err != nil {
		return err
	}
	return update(db, existingEntity, id, DeleteOperationHandler)
}
