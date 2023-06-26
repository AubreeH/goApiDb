package access

import (
	"reflect"
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func Delete[T any](db *database.Database, entity T, id any) error {
	_, err := delete[T](db, entity, id, false)
	return err
}

func DeleteTimed[T any](db *database.Database, entity T, id any) (*TimedResult, error) {
	return delete[T](db, entity, id, true)
}

func delete[T any](db *database.Database, entity T, id any, timed bool) (*TimedResult, error) {
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

	var err error
	entity, err = GetById(db, entity, id)
	if err != nil {
		return nil, err
	}

	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return nil, err
	}

	if doesEntitySoftDelete(entity) {
		if timed {
			buildQueryDurationEnd = time.Now()
		}

		timedResult, err := softDelete(db, entity, id, timed)
		if err != nil {
			return nil, err
		}

		if timed {
			overallDurationEnd = time.Now()

			return &TimedResult{
				QueryExecDuration:  timedResult.QueryExecDuration,
				BuildQueryDuration: buildQueryDurationEnd.UnixMicro() - buildQueryDurationStart.UnixMicro() + timedResult.BuildQueryDuration,
				OverallDuration:    overallDurationEnd.UnixMicro() - overallDurationStart.UnixMicro(),
			}, nil
		}

		return nil, nil
	}

	_, err = deleteOperationHandler(reflect.ValueOf(entity))
	if err != nil {
		return nil, err
	}

	q := "DELETE FROM " + tableInfo.Name + " WHERE ID = ?"

	if timed {
		buildQueryDurationEnd = time.Now()
		queryExecDurationStart = time.Now()
	}

	_, err = db.Db.Exec(q, id)
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

	return nil, err
}

func softDelete[T any](db *database.Database, entity T, id any, timed bool) (*TimedResult, error) {
	existingEntity, err := GetById(db, entity, id)
	if err != nil {
		return nil, err
	}
	return update(db, existingEntity, id, deleteOperationHandler, timed)
}
