package access

import (
	"errors"
	"goApiDb/helpers"
	"reflect"
	"strings"
)

func GetDataAndId(value any, operationHandler OperationHandler) (ColumnData, []ColumnData, error) {
	refValue := reflect.ValueOf(value)
	data, err := getColumnData(refValue, operationHandler)
	if err != nil {
		return ColumnData{}, nil, err
	}
	return separateIdAndData(data)
}

func GetData(value any, operationHandler OperationHandler) ([]ColumnData, error) {
	refValue := reflect.ValueOf(value)
	return getColumnData(refValue, operationHandler)
}

func getColumnData(refValue reflect.Value, operationHandler OperationHandler) ([]ColumnData, error) {
	var data []ColumnData

	if refValue.Kind() == reflect.Struct {
		var err error
		refValue, err = operationHandler(refValue)
		if err != nil {
			return nil, err
		}

		structHasParser, parser := checkForParser(refValue)
		if structHasParser {
			result := parser.Call([]reflect.Value{})[0]
			if result.Type() == reflect.TypeOf(data) {
				resultColumnData := result.Interface().([]ColumnData)
				data = append(data, resultColumnData...)
			}
		} else {
			columnData, err := extractData(refValue, operationHandler)
			if err != nil {
				return nil, err
			}

			data = append(data, columnData...)
		}

	} else {
		return nil, errors.New("unable to parse provided value as it is not a struct")
	}

	return data, nil
}

func extractData(refValue reflect.Value, operationHandler OperationHandler) ([]ColumnData, error) {
	var data []ColumnData
	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Field(i)
		fieldType := refValue.Type().Field(i)

		if helpers.ParseBool(fieldType.Tag.Get("sql_ignore")) {
			continue
		}

		fieldExtractor := field.MethodByName("ExtractDataFunc")

		if field.Kind() == reflect.Struct && !fieldExtractor.IsValid() {
			result, err := getColumnData(field, operationHandler)
			if err != nil {
				return nil, err
			}
			data = append(data, result...)
			continue
		}

		sqlName := fieldType.Tag.Get("sql_name")
		nullable := helpers.ParseBool(fieldType.Tag.Get("sql_nullable"))
		primaryKey := strings.ToLower(fieldType.Tag.Get("sql_key")) == "primary"
		hasDefault := fieldType.Tag.Get("sql_default") != ""
		var fieldData any
		if fieldExtractor.IsValid() {
			fieldData = fieldExtractor.Call([]reflect.Value{})[0].Interface()
		} else {
			fieldData = field.Interface()
		}

		if primaryKey && hasDefault {
			data = append(data, ColumnData{Column: sqlName, PrimaryKey: primaryKey, Data: nil, PrioritiseExisting: !nullable})
		} else {
			data = append(data, ColumnData{Column: sqlName, PrimaryKey: primaryKey, Data: fieldData, PrioritiseExisting: !nullable})
		}
	}
	return data, nil
}

func separateIdAndData(data []ColumnData) (ColumnData, []ColumnData, error) {
	var colData []ColumnData

	for i := range data {
		value := data[i]
		if value.PrimaryKey {
			return value, append(colData, data[(i+1):]...), nil
		} else {
			colData = append(colData, value)
		}
	}

	return ColumnData{}, colData, errors.New("unable to find primary key")
}
