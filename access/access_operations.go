package access

import "reflect"

func CreateOperationHandler(refValue reflect.Value) (reflect.Value, error) {
	return callMethodIfExists("OnCreate", refValue)
}

func UpdateOperationHandler(refValue reflect.Value) (reflect.Value, error) {
	return callMethodIfExists("OnUpdate", refValue)
}

func DeleteOperationHandler(refValue reflect.Value) (reflect.Value, error) {
	return callMethodIfExists("OnDelete", refValue)
}

func NonOperationHandler(refValue reflect.Value) (reflect.Value, error) {
	return refValue, nil
}
