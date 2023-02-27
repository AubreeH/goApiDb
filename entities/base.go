package entities

import (
	"errors"
	"reflect"
)

type EntityBase struct {
}

type TableInfo struct {
	Name string
}

func GetTableInfo(entity interface{}) (TableInfo, error) {
	entityVal := reflect.ValueOf(entity)
	entityBaseType := reflect.TypeOf(EntityBase{})

	if entityVal.Kind() == reflect.Ptr || entityVal.Kind() == reflect.Interface {
		entityVal = entityVal.Elem()

		if entityVal.Kind() == reflect.Interface {
			entityVal = entityVal.Elem()
		}
	}

	entityType := entityVal.Type()

	if entityType.Kind() != reflect.Struct {
		return TableInfo{}, errors.New("provided entity is not a struct")
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if field.Type == entityBaseType {
			tableInfo := TableInfo{Name: field.Tag.Get("table_name")}
			return tableInfo, nil
		}
	}

	return TableInfo{}, errors.New("no entity base in struct")
}
