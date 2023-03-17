package access

import (
	"fmt"
	"reflect"
)

// createOperationHandler runs the OnDelete function for the provided value. This is used prior to Create operations.
func createOperationHandler(refValue reflect.Value) (value reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an error occurred whilst running the onCreate operation handler:%v", r.(error))
		}
	}()

	return callMethodIfExists("OnCreate", refValue)
}

// updateOperationHandler runs the OnDelete function for the provided value. This is used prior to Update operations.
func updateOperationHandler(refValue reflect.Value) (value reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an error occurred whilst running the onCreate operation handler:%v", r.(error))
		}
	}()

	return callMethodIfExists("OnUpdate", refValue)
}

// deleteOperationHandler runs the OnDelete function for the provided value. This is used prior to Delete operations.
func deleteOperationHandler(refValue reflect.Value) (value reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an error occurred whilst running the onCreate operation handler:%v", r.(error))
		}
	}()

	return callMethodIfExists("OnDelete", refValue)
}

// nilOperationHandler is used to extract data without updating the original values.
func nilOperationHandler(refValue reflect.Value) (reflect.Value, error) {
	return refValue, nil
}
