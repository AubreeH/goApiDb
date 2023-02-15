package helpers

func ArrFind[T comparable](array []T, item T) (T, bool) {
	for i := range array {
		if item == array[i] {
			return array[i], true
		}
	}

	var result T
	return result, false
}

func ArrFindIndex[T comparable](array []T, item T) int {
	for i := range array {
		if item == array[i] {
			return i
		}
	}

	return -1
}

func ArrFindFunc[T any](array []T, handler func(item T) bool) (T, bool) {
	for i := range array {
		if handler(array[i]) {
			return array[i], true
		}
	}

	var result T
	return result, false
}

func ArrFindIndexFunc[T any](array []T, handler func(item T) bool) int {
	for i := range array {
		if handler(array[i]) {
			return i
		}
	}

	return -1
}
