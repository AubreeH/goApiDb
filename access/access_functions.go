package access

import (
	"reflect"

	"github.com/AubreeH/goApiDb/helpers"
	"github.com/AubreeH/goApiDb/structParsing"
)

func MergeObjects[T any](base T, merger T) T {
	baseValue := helpers.GetRootValue(reflect.ValueOf(&base))
	mergerValue := reflect.ValueOf(merger)

	mergeDataFuncMethod := baseValue.MethodByName("MergeDataFunc")

	if mergeDataFuncMethod.IsValid() {
		baseValue = mergeDataFuncMethod.Call([]reflect.Value{mergerValue})[0]
		if baseValue.Kind() == reflect.Pointer {
			return helpers.GetRootValue(baseValue).Interface().(T)
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

		blockExternalMod := structParsing.FormatBoolean(structParsing.DbDisallowExternalModification.Get(baseFieldType))
		if blockExternalMod == 1 {
			continue
		} else {
			parseStruct := structParsing.FormatBoolean(structParsing.DbParseStruct.Get(baseFieldType))
			if baseField.Kind() == reflect.Struct && (parseStruct != 0) {
				baseField.Set(mergeRefValues(baseField, mergerField))
			} else {
				baseField.Set(mergerField)
			}
		}
	}
	return base
}

func doesEntitySoftDelete(entity any) bool {
	refValue := reflect.ValueOf(entity)

	if refValue.Kind() == reflect.Struct {
		for i := 0; i < refValue.NumField(); i++ {
			fieldValue := refValue.Field(i)
			fieldType := refValue.Type().Field(i)

			if structParsing.FormatBoolean(structParsing.DbSoftDeletes.Get(fieldType)) == 1 {
				return true
			}
			if fieldValue.Type().Kind() == reflect.Struct && structParsing.FormatParseStruct(fieldType) {
				if doesEntitySoftDelete(fieldValue.Interface()) {
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
