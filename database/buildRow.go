package database

import (
	"database/sql"
	"reflect"

	"github.com/AubreeH/goApiDb/helpers"
	"github.com/AubreeH/goApiDb/structParsing"
)

func BuildRow[T interface{}](entity T, result *sql.Rows) ([]interface{}, *T, error) {
	columnVariables, ptr, _, err := getEntityConstruction(&entity)
	if err != nil {
		return nil, ptr, err
	}
	resultColumns, err := result.Columns()
	if err != nil {
		return nil, ptr, err
	}
	columns := resultColumns

	retArgs := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		retArgs[i] = columnVariables[columns[i]]
	}

	return retArgs, ptr, nil
}

func getEntityConstruction[T any](entity *T) (map[string]any, *T, string, error) {
	val := helpers.GetRootValue(reflect.ValueOf(entity))
	tmp := helpers.GetRootValue(reflect.New(val.Type()))
	tmp.Set(val)

	columnVariables := make(map[string]any)
	getColumnsFromStruct(tmp, columnVariables)

	tableName := structParsing.GetTableName(entity)

	return columnVariables, tmp.Addr().Interface().(*T), tableName, nil
}

func getColumnsFromStruct(refValue reflect.Value, columnVariables map[string]any) map[string]any {
	numFields := refValue.NumField()
	for i := 0; i < numFields; i++ {
		field := refValue.Type().Field(i)

		if structParsing.FormatSqlIgnore(field) {
			continue
		}

		valueField := refValue.Field(i)
		getPtrFunc := valueField.MethodByName("GetPtrFunc")

		sqlName := structParsing.FormatSqlName(field)
		if getPtrFunc.IsValid() {
			result := getPtrFunc.Call([]reflect.Value{valueField.Addr()})[0]
			if helpers.GetRootValue(result).Kind() == reflect.Map {
				val := helpers.GetRootValue(result).Interface().(map[string]any)
				for s, p := range val {
					columnVariables[s] = p
				}
			} else if result.Kind() == reflect.Pointer {
				columnVariables[sqlName] = result.Interface()
			}
		} else if valueField.Kind() == reflect.Struct && structParsing.FormatParseStruct(field) {
			getColumnsFromStruct(valueField, columnVariables)
		} else {
			columnVariables[sqlName] = valueField.Addr().Interface()
		}
	}
	return columnVariables
}
