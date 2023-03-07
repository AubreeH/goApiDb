package helpers

import "reflect"

func TagLookup(refType reflect.StructField, key string, output *string) (ok bool) {
	val, success := refType.Tag.Lookup(key)
	if success {
		*output = val
	}

	return ok
}
