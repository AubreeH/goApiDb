package access

import "reflect"

type ColumnData struct {
	Column             string
	Data               any
	PrimaryKey         bool
	PrioritiseExisting bool
}

type OperationHandler = func(value reflect.Value) (reflect.Value, error)

type ExtractDataFunc = func() any

type ParseDataFunc = func() []ColumnData

type MergeDataFunc = func() any
