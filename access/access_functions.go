package access

import (
	"goApiDb/helpers"
	"reflect"
)

func MergeObjects[T any](base T, merger T) T {
	baseValue := reflect.ValueOf(&base).Elem()
	mergerValue := reflect.ValueOf(merger)

	mergeDataFuncMethod := baseValue.MethodByName("MergeDataFunc")

	if mergeDataFuncMethod.IsValid() {
		baseValue = mergeDataFuncMethod.Call([]reflect.Value{mergerValue})[0]
		if baseValue.Kind() == reflect.Pointer {
			return baseValue.Elem().Interface().(T)
		}
		return baseValue.Interface().(T)
	}

	return mergeRefValues(baseValue, mergerValue).Interface().(T)
}

func mergeRefValues(base reflect.Value, merger reflect.Value) reflect.Value {
	for i := 0; i < base.NumField(); i++ {
		baseField := base.Field(i)
		baseFieldType := base.Type().Field(i)
		mergerField := merger.Field(i)

		blockExternalMod := baseFieldType.Tag.Get("sql_disallow_external_modification")
		if helpers.ParseBool(blockExternalMod) {
			continue
		} else {
			parseStruct := baseFieldType.Tag.Get("parse_struct")
			if baseField.Kind() == reflect.Struct && (parseStruct == "" || helpers.ParseBool(parseStruct)) {
				baseField.Set(mergeRefValues(baseField, mergerField))
			} else {
				baseField.Set(mergerField)
			}
		}
	}
	return base
}

func DoesEntitySoftDelete(entity any) bool {
	refValue := reflect.ValueOf(entity)

	if refValue.Kind() == reflect.Struct {
		for i := 0; i < refValue.NumField(); i++ {
			fieldValue := refValue.Field(i)
			fieldType := refValue.Type().Field(i)

			if helpers.ParseBool(fieldType.Tag.Get("soft_deletes")) {
				return true
			}

			if fieldValue.Type().Kind() == reflect.Struct && fieldType.Tag.Get("parse_struct") != "false" {
				if DoesEntitySoftDelete(fieldValue.Interface()) {
					return true
				}
			}
		}
	}

	return false
}

func callMethodIfExists(methodName string, refValue reflect.Value) (reflect.Value, error) {
	method := refValue.MethodByName(methodName)
	if method.IsValid() && method.Type().NumOut() == 2 && method.Type().Out(1).Name() == "error" {
		results := method.Call([]reflect.Value{})
		if len(results) == 2 {
			err := results[1]
			if !err.IsNil() {
				return refValue, err.Interface().(error)
			}

			result := results[0]
			if result.Type() == refValue.Type() {
				return result, nil
			}
		}

	}
	return refValue, nil
}

func checkForParser(refValue reflect.Value) (bool, reflect.Value) {
	method := refValue.MethodByName("ParseDataFunc")
	var f ExtractDataFunc
	funcType := reflect.TypeOf(f)

	if method.IsValid() && method.Type() == funcType {
		return true, method
	}

	return false, method
}
