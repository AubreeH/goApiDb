package helpers

import "reflect"

func GetRootValue(refValue reflect.Value) reflect.Value {
	for refValue.Kind() == reflect.Ptr || refValue.Kind() == reflect.Interface {
		refValue = refValue.Elem()
	}
	return refValue
}
