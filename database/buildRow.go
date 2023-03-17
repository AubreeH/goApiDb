package database

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/structParsing"
	"reflect"
)

func getEntityConstruction[T any](entity *T) (map[string]any, T, string, error) {
	val := reflect.ValueOf(entity).Elem()

	tmp := reflect.New(val.Elem().Type()).Elem()
	tmp.Set(val.Elem())

	columnVariables := make(map[string]any)
	getColumnsFromStruct(tmp, columnVariables)

	currentValue := reflect.ValueOf(entity).Elem().Interface()
	tableInfo, err := structParsing.GetTableInfo(currentValue)
	if err != nil {
		var output T
		return nil, output, "", err
	}

	return columnVariables, tmp.Addr().Interface().(T), tableInfo.Name, nil
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
			if result.Elem().Kind() == reflect.Map {
				val := result.Elem().Interface().(map[string]any)
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

func BuildRow(db *Database, entity interface{}, result *sql.Rows) ([]interface{}, interface{}, error) {
	columnVariables, ptr, tableName, err := getEntityConstruction(&entity)
	if err != nil {
		return nil, ptr, err
	}

	var columns []string
	if db.tableColumns[tableName] != nil {
		columns = db.tableColumns[tableName]
	} else {
		resultColumns, err := result.Columns()
		if err != nil {
			return nil, ptr, err
		}
		db.tableColumns[tableName] = resultColumns
		columns = resultColumns
	}

	retArgs := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		retArgs[i] = columnVariables[columns[i]]
	}

	return retArgs, ptr, nil
}
