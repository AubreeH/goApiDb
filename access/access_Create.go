package access

import (
	"fmt"
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

// Create new records in the database. Returns the last inserted value.
func Create[T any](db database.DbInstance, values ...T) (T, error) {
	out, _, err := create[T](db, values, false)
	return out, err
}

// Create new records in the database. Returns the last inserted value.
func CreateTimed[T any](db database.DbInstance, values ...T) (T, *TimedResult, error) {
	out, timedResult, err := create[T](db, values, true)
	return out, timedResult, err
}

func create[T any](db database.DbInstance, values []T, timed bool) (T, *TimedResult, error) {
	var overallDurationStart time.Time
	var overallDurationEnd time.Time
	var buildQueryDurationStart time.Time
	var buildQueryDurationEnd time.Time
	var queryExecDurationStart time.Time
	var queryExecDurationEnd time.Time

	if timed {
		overallDurationStart = time.Now()
		buildQueryDurationStart = time.Now()
	}

	var entity T
	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return entity, nil, err
	}

	queryColumns := ""
	queryValues := ""
	var args []any

	var id interface{}

	for i := range values {
		var rowData []ColumnData
		rowData, err = GetData(values[i], createOperationHandler)
		if err != nil {
			return entity, nil, err
		}

		for j := range rowData {
			columnData := rowData[j]

			if columnData.Column == tableInfo.PrimaryKey {
				id = columnData.Data
			}

			if queryColumns == "" {
				queryColumns += columnData.Column
				queryValues += "?"
			} else {
				queryColumns += ", " + columnData.Column
				queryValues += ", ?"
			}

			args = append(args, columnData.Data)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) values (%s)", tableInfo.Name, queryColumns, queryValues)

	if timed {
		buildQueryDurationEnd = time.Now()
		queryExecDurationStart = time.Now()
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return entity, nil, err
	}

	if lastInsertId, err := result.LastInsertId(); err != nil {
		return entity, nil, err
	} else if lastInsertId != 0 {
		id = lastInsertId
	}

	if id != nil {
		entity, _, err = getById(db, entity, id, false)
		if err != nil {
			return entity, nil, err
		}
	}

	if timed {
		queryExecDurationEnd = time.Now()
		overallDurationEnd = time.Now()

		return entity, &TimedResult{
			BuildQueryDuration: buildQueryDurationEnd.UnixMicro() - buildQueryDurationStart.UnixMicro(),
			QueryExecDuration:  queryExecDurationStart.UnixMicro() - queryExecDurationEnd.UnixMicro(),
			OverallDuration:    overallDurationStart.UnixMicro() - overallDurationEnd.UnixMicro(),
		}, nil
	}

	return entity, nil, err
}
