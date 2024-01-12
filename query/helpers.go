package query

import (
	"database/sql"
	"reflect"

	"github.com/AubreeH/goApiDb/helpers"
	"github.com/AubreeH/goApiDb/structParsing"
)

func resetStruct[T any](r *T) {
	refVal := helpers.GetRootValue(reflect.ValueOf(r))
	refVal.Set(reflect.Zero(refVal.Type()))
}

func getScannableRow[TScanType any](rs *sql.Rows) (*TScanType, []interface{}, error) {
	var out TScanType

	columns, err := rs.Columns()
	if err != nil {
		return nil, nil, err
	}

	args := make([]interface{}, len(columns))
	ptrs := getPointers(&out)

	for i, col := range columns {
		args[i] = ptrs[col]
	}

	return &out, args, nil
}

func getPointers[T any](s *T) map[string]interface{} {
	out := make(map[string]interface{})

	refVal := helpers.GetRootValue(reflect.ValueOf(s))
	getPointersFromRefValue(refVal, out)

	return out
}

func getPointersFromRefValue(refVal reflect.Value, out map[string]interface{}) {
	for i := 0; i < refVal.NumField(); i++ {
		field := refVal.Type().Field(i)
		valueField := refVal.Field(i)

		colName := field.Name
		if colName == "" {
			continue
		}

		sqlName := structParsing.FormatSqlName(field)

		getPtrFunc := valueField.MethodByName("GetPtrFunc")
		if getPtrFunc.IsValid() {
			result := getPtrFunc.Call([]reflect.Value{valueField.Addr()})[0]
			if result.Elem().Kind() == reflect.Map {
				val := result.Elem().Interface().(map[string]any)
				for s, p := range val {
					out[s] = p
				}
			} else if result.Kind() == reflect.Pointer {
				out[sqlName] = result.Interface()
			} else if result.Kind() != reflect.Invalid {
				out[sqlName] = result.Addr().Interface()
			}
		} else if valueField.Kind() == reflect.Struct && structParsing.FormatParseStruct(field) {
			getPointersFromRefValue(valueField, out)
		} else {
			out[sqlName] = valueField.Addr().Interface()
		}
	}
}

func scanRows[TScanType any](rs *sql.Rows) ([]TScanType, error) {
	out := []TScanType{}
	row, rowArgs, err := getScannableRow[TScanType](rs)
	if err != nil {
		return []TScanType{}, err
	}

	for rs.Next() {
		err := rs.Scan(rowArgs...)
		if err != nil {
			return []TScanType{}, err
		}

		out = append(out, *row)
		resetStruct(row)
	}
	return out, nil
}

func scanRow[TScanType any](rs *sql.Rows) (TScanType, error) {
	var out TScanType

	row, rowArgs, err := getScannableRow[TScanType](rs)
	if err != nil {
		return out, err
	}

	if !rs.Next() {
		return out, nil
	}

	err = rs.Scan(rowArgs...)
	if err != nil {
		return out, err
	}

	return *row, nil
}

func tempSet[T any](ptr *T, value T) func() {
	refVal := helpers.GetRootValue(reflect.ValueOf(ptr))
	oldVal := refVal.Interface()
	refVal.Set(reflect.ValueOf(value))

	return func() {
		refVal.Set(reflect.ValueOf(oldVal))
	}
}
