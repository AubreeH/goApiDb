package query

type QuerySettings struct {
	preventFieldNameAutoMapping bool
}

// PreventFieldNameAutoMapping prevents the query from automatically mapping the struct field names to the database column names. Use this if the auto mapping is causing issues.
func (q *QuerySettings) PreventFieldNameAutoMapping() {
	q.preventFieldNameAutoMapping = true
}
