package access

import (
	"reflect"
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func GetAll[T any](db database.DbInstance, entity T, limit int) ([]T, error) {
	out, _, err := getAll[T](db, entity, limit, false)
	return out, err
}

func GetAllTimed[T any](db database.DbInstance, entity T, limit int) ([]T, *TimedResult, error) {
	out, timedResult, err := getAll[T](db, entity, limit, true)
	return out, timedResult, err
}

func getAll[T any](db database.DbInstance, entity T, limit int, timed bool) ([]T, *TimedResult, error) {
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

	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return nil, nil, err
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

	if timed {
		buildQueryDurationEnd = time.Now()
		queryExecDurationStart = time.Now()
	}

	result, err := db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}

	defer result.Close()

	if timed {
		queryExecDurationEnd = time.Now()
		formatResultDurationStart = time.Now()
	}

	var retValue []T
	args, entityOutput, err := database.BuildRow(entity, result)
	for result.Next() {
		if err != nil {
			return nil, nil, err
		}
		err = result.Scan(args...)
		if entityOutput != nil {
			retValue = append(retValue, reflect.ValueOf(entityOutput).Elem().Interface().(T))
		}
	}

	if timed {
		formatResultDurationEnd = time.Now()
		overallDurationEnd = time.Now()

		return retValue, &TimedResult{
			BuildQueryDuration:   buildQueryDurationEnd.UnixMicro() - buildQueryDurationStart.UnixMicro(),
			QueryExecDuration:    queryExecDurationEnd.UnixMicro() - queryExecDurationStart.UnixMicro(),
			OverallDuration:      overallDurationEnd.UnixMicro() - overallDurationStart.UnixMicro(),
			FormatResultDuration: formatResultDurationEnd.UnixMicro() - formatResultDurationStart.UnixMicro(),
		}, nil
	}

	return retValue, nil, nil
}
