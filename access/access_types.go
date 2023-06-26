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

type TimedResult struct {
	// Duration of build function for access query in microseconds
	BuildQueryDuration int64

	// Duration of query execution for access query in microseconds
	QueryExecDuration int64

	// Duration of result formatting for access query in microseconds
	FormatResultDuration int64

	// Total duration of function for access query in microseconds
	OverallDuration int64
}
