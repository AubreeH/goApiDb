package database

import (
	"errors"
	"reflect"
)

func GetTableSqlDescription[TModel any]() error {
	var model TModel

	refValue := reflect.ValueOf(model)
	refType := reflect.TypeOf(model)

	if refType.Kind() != reflect.Struct {
		return errors.New("provided type is not a struct")
	}

	var columns []string

	for i := 0; i < refValue.NumField(); i++ {
		columns = append(columns, parseColumn(refValue.Field(i)))
	}

	return errors.New("not yet implemented")
}

func parseColumn(col reflect.Value) string {

	return ""
}
