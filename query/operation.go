package query

type operation string

const (
	selectOperation operation = "SELECT"
	insertOperation operation = "INSERT"
	updateOperation operation = "UPDATE"
	deleteOperation operation = "DELETE"
)

func (o operation) format(s interface{}, pretty bool) string {
	return string(o)
}