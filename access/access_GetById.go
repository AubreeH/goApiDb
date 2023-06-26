package access

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func GetById[T any](db *database.Database, entity T, id any) (T, error) {
	out, _, err := getById[T](db, entity, id, false)
	return out, err
}

func GetByIdTimed[T any](db *database.Database, entity T, id any) (T, *TimedResult, error) {
	out, timedResult, err := getById[T](db, entity, id, true)
	return out, timedResult, err
}

func getById[T any](db *database.Database, entity T, id any, timed bool) (T, *TimedResult, error) {
	var overallDurationStart time.Time
	var overallDurationEnd time.Time
	var buildQueryDurationStart time.Time
	var buildQueryDurationEnd time.Time
	var queryExecDurationStart time.Time
	var queryExecDurationEnd time.Time
	var formatResultDurationStart time.Time
	var formatResultDurationEnd time.Time

	if timed {
		overallDurationStart = time.Now()
		buildQueryDurationStart = time.Now()
	}

	var output T

	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return output, nil, err
	}

	var query string
	if tableInfo.SoftDeletes != "" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s IS NULL AND id = ? LIMIT 1", tableInfo.Name, tableInfo.SoftDeletes)
	} else {
		query = fmt.Sprintf("SELECT *  FROM %s WHERE id = ? LIMIT 1", tableInfo.Name)
	}

	if timed {
		buildQueryDurationEnd = time.Now()
		queryExecDurationStart = time.Now()
	}

	result, err := db.Db.Query(query, id)
	if err != nil {
		return output, nil, err
	}

	if timed {
		queryExecDurationEnd = time.Now()
		formatResultDurationStart = time.Now()
	}

	args, entityOutput, err := database.BuildRow(db, entity, result)
	if err != nil {
		return output, nil, err
	}

	if !result.Next() {
		return output, nil, errors.New("unable to find value")
	}

	err = result.Scan(args...)
	if err != nil {
		return output, nil, err
	}

	output = reflect.ValueOf(entityOutput).Elem().Interface().(T)

	if timed {
		formatResultDurationEnd = time.Now()
		overallDurationEnd = time.Now()

		return output, &TimedResult{
			BuildQueryDuration:   buildQueryDurationEnd.UnixMicro() - buildQueryDurationStart.UnixMicro(),
			QueryExecDuration:    queryExecDurationEnd.UnixMicro() - queryExecDurationStart.UnixMicro(),
			OverallDuration:      overallDurationEnd.UnixMicro() - overallDurationStart.UnixMicro(),
			FormatResultDuration: formatResultDurationEnd.UnixMicro() - formatResultDurationStart.UnixMicro(),
		}, nil
	}

	return output, nil, nil
}
