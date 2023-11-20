package access

import (
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func Update[T any](db database.DbInstance, value T, id any) error {
	_, err := update(db, value, id, updateOperationHandler, false)
	return err
}

func UpdateTimed[T any](db database.DbInstance, value T, id any) (*TimedResult, error) {
	return update[T](db, value, id, updateOperationHandler, true)
}

func update[T any](db database.DbInstance, value T, id any, operationHandler OperationHandler, timed bool) (*TimedResult, error) {
	var overallDurationStart time.Time
	var overallDurationEnd time.Time
	var buildQueryDurationStart time.Time
	var buildQueryDurationEnd time.Time
	var queryExecDurationStart time.Time
	var queryExecDurationEnd time.Time

	var output T

	if timed {
		overallDurationStart = time.Now()
		buildQueryDurationStart = time.Now()
	}

	existingValue, err := GetById(db, output, id)
	if err != nil {
		return nil, err
	}

	mergedValue := MergeObjects(existingValue, value)
	idColumn, mergedData, err := GetDataAndId(mergedValue, operationHandler)
	if err != nil {
		return nil, err
	}
	_, updateData, err := GetDataAndId(existingValue, nilOperationHandler)
	if err != nil {
		return nil, err
	}

	tableInfo, err := structParsing.GetTableInfo(value)
	if err != nil {
		return nil, err
	}

	q := "UPDATE " + tableInfo.Name + " t SET "
	qBase := q
	var where string

	// if doesEntitySoftDelete(value) {
	// 	where = " WHERE deleted = false AND t."
	// } else {
	// 	where = " WHERE t."
	// }

	where = " WHERE t."

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

	if timed {
		buildQueryDurationEnd = time.Now()
		queryExecDurationStart = time.Now()
	}

	_, err = db.Exec(q+where, args...)
	if err != nil {
		return nil, err
	}

	if timed {
		queryExecDurationEnd = time.Now()
		overallDurationEnd = time.Now()

		return &TimedResult{
			BuildQueryDuration: buildQueryDurationEnd.UnixMicro() - buildQueryDurationStart.UnixMicro(),
			QueryExecDuration:  queryExecDurationEnd.UnixMicro() - queryExecDurationStart.UnixMicro(),
			OverallDuration:    overallDurationEnd.UnixMicro() - overallDurationStart.UnixMicro(),
		}, nil
	}

	return nil, nil
}
