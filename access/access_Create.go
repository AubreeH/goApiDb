package access

import (
	"fmt"
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func Create[T any](db *database.Database, values []T) error {
	_, err := create[T](db, values, false)
	return err
}

func CreateTimed[T any](db *database.Database, values []T) (*TimedResult, error) {
	timedResult, err := create[T](db, values, true)
	return timedResult, err
}

func create[T any](db *database.Database, values []T, timed bool) (*TimedResult, error) {
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
		return nil, err
	}

	queryColumns := ""
	queryValues := ""
	var args []any

	for i := range values {
		var rowData []ColumnData
		rowData, err = GetData(values[i], createOperationHandler)
		if err != nil {
			return nil, err
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

			args = append(args, columnData.Data)
		}
	}

	query := fmt.Sprintf("INSERT"+" INTO %s (%s) values (%s)", tableInfo.Name, queryColumns, queryValues)

	if timed {
		buildQueryDurationEnd = time.Now()
		queryExecDurationStart = time.Now()
	}

	_, err = db.Db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	if timed {
		queryExecDurationEnd = time.Now()
		overallDurationEnd = time.Now()

		return &TimedResult{
			BuildQueryDuration: buildQueryDurationEnd.UnixMicro() - buildQueryDurationStart.UnixMicro(),
			QueryExecDuration:  queryExecDurationStart.UnixMicro() - queryExecDurationEnd.UnixMicro(),
			OverallDuration:    overallDurationStart.UnixMicro() - overallDurationEnd.UnixMicro(),
		}, nil
	}

	return nil, err
}
