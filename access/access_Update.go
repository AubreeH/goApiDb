package access

import (
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/entities"
)

func Update[T any](db *database.Database, value T, id any) error {
	return update(db, value, id, UpdateOperationHandler)
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

	tableInfo, err := entities.GetTableInfo(value)
	if err != nil {
		return err
	}

	q := "UPDATE " + tableInfo.Name + " t SET "
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
