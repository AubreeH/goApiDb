package access

import (
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/entities"
	"reflect"
)

func Delete[T any](db *database.Database, entity T, id any) error {
	var err error
	entity, err = GetById(db, entity, id)
	if err != nil {
		return err
	}

	tableInfo, err := entities.GetTableInfo(entity)
	if err != nil {
		return err
	}

	if DoesEntitySoftDelete(entity) {
		return softDelete(db, entity, id)
	} else {
		_, err = DeleteOperationHandler(reflect.ValueOf(entity))
		if err != nil {
			return err
		}

		q := "DELETE FROM " + tableInfo.Name + " WHERE ID = ?"

		_, err = db.Db.Exec(q, id)
		return err
	}
}

func softDelete[T any](db *database.Database, entity T, id any) error {
	existingEntity, err := GetById(db, entity, id)
	if err != nil {
		return err
	}
	return update(db, existingEntity, id, DeleteOperationHandler)
}
