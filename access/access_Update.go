package access

import (
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func Update[T any](db *database.Database, value T, id any) error {
	return update(db, value, id, updateOperationHandler)
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
	_, updateData, err := GetDataAndId(existingValue, nilOperationHandler)
	if err != nil {
		return err
	}

	tableInfo, err := structParsing.GetTableInfo(value)
	if err != nil {
		return err
	}

	q := "UPDATE " + tableInfo.Name + " t SET "
	qBase := q
	var where string

	if doesEntitySoftDelete(value) {
		where = " WHERE deleted = false AND t."
	} else {
		where = " WHERE t."
	}

	var args []any
	for j, column := range mergedData {
		if updateData[j].Data != column.Data && !column.PrioritiseExisting {
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
