package helpers

func ArrAdd[T any](arr *[]T, val ...T) {
	*arr = append(*arr, val...)
}
