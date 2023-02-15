package goApiDb

import "goApiDb/queryBuilder"

func ExecuteQuery[T any](q *query.Query, _ T) ([]T, error) {
	q.Exec()
	if q.Error != nil {
		return nil, q.Error
	}

	var output []T

	result := query.GetResult(q)
	for result.Next() {
		var row T
		args := query.GetRowArgs(&row)
		err := result.Scan(args...)
		if err != nil {
			return nil, err
		}
		output = append(output, row)
	}

	return output, nil
}

// -------------------------------
// NEW QUERY FUNCTIONS
// -------------------------------

func NewSelectQuery() *query.Query {
	return query.NewSelectQuery()
}

func NewUpdateQuery() *query.Query {
	return query.NewUpdateQuery()
}
