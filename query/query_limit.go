package query

import (
	"fmt"
)

func (query *Query) Limit(value uint) *Query {
	query.limit = fmt.Sprintf("%d", value)
	return query
}

func (query *Query) Offset(value uint) *Query {
	query.offset = fmt.Sprintf("%d", value)
	return query
}

func (query *Query) Paginated(itemsPerPage, offset uint) *Query {
	query.limit = fmt.Sprintf("%d", itemsPerPage)
	query.offset = fmt.Sprintf("%d", offset)
	return query
}
