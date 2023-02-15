package goApiDb

import (
	"errors"
	"fmt"
	"goApiDb/access"
	"goApiDb/database"
	"goApiDb/helpers"
	"reflect"
)

func GetById[T any](entity T, id any) (T, error) {
	db := database.GetDb()
	tableName := helpers.GetTableName(entity)

	var query string
	if access.DoesEntitySoftDelete(entity) {
		query = "SELECT *" + " FROM " + tableName + " WHERE deleted = false AND id = ? LIMIT 1"
	} else {
		query = "SELECT *" + " FROM " + tableName + " WHERE id = ? LIMIT 1"
	}

	result, err := db.Query(query, id)
	if err != nil {
		return entity, err
	}

	args, entityOutput, err := database.BuildRow(entity, result)
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

func GetAll[T any](entity T, limit int) ([]T, error) {
	db := database.GetDb()
	tableName := helpers.GetTableName(entity)

	var args []any

	query := "SELECT * FROM " + tableName

	if access.DoesEntitySoftDelete(entity) {
		query += " WHERE deleted = false"
	}

	if limit != 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	result, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var retValue []T
	args, entityOutput, err := database.BuildRow(entity, result)
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

func Create[T any](values []T) (T, error) {
	db := database.GetDb()
	var entity T
	tableName := helpers.GetTableName(entity)

	queryColumns := ""
	queryValues := ""
	var args []any

	var output T
	var id any
	for i := range values {
		rowData, err := access.GetData(values[i], access.CreateOperationHandler)
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

	res, err := db.Exec(query, args...)
	if err != nil {
		return output, err
	}

	if !(id == nil || id == 0 || id == "") {
		return GetById(output, id)
	}

	intId, err := res.LastInsertId()
	if err != nil {
		return output, err
	}

	newEntity, err := GetById(output, intId)

	return newEntity, err
}

func Update[T any](value T, id any) error {
	return update(value, id, access.UpdateOperationHandler)
}

func Delete[T any](entity T, id any) error {
	_, err := GetById(entity, id)
	if err != nil {
		return err
	}

	tableName := helpers.GetTableName(entity)
	db := database.GetDb()

	if access.DoesEntitySoftDelete(entity) {
		return softDelete(entity, id)
	} else {
		_, err = access.DeleteOperationHandler(reflect.ValueOf(entity))
		if err != nil {
			return err
		}

		q := "DELETE FROM " + tableName + " WHERE ID = ?"

		_, err = db.Exec(q, id)
		return err
	}
}

func update[T any](value T, id any, operationHandler access.OperationHandler) error {
	var output T

	existingValue, err := GetById(output, id)
	if err != nil {
		return err
	}

	mergedValue := access.MergeObjects(existingValue, value)
	idColumn, mergedData, err := access.GetDataAndId(mergedValue, operationHandler)
	if err != nil {
		return err
	}
	_, updateData, err := access.GetDataAndId(value, access.NonOperationHandler)
	if err != nil {
		return err
	}

	tableName := helpers.GetTableName(value)
	db := database.GetDb()

	q := "UPDATE " + tableName + " t SET "
	qBase := q
	var where string

	if access.DoesEntitySoftDelete(value) {
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
	_, err = db.Exec(q+where, args...)
	if err != nil {
		return err
	}

	return nil
}

func softDelete[T any](entity T, id any) error {
	existingEntity, err := GetById(entity, id)
	if err != nil {
		return err
	}
	return update(existingEntity, id, access.DeleteOperationHandler)
}
