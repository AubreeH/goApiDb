package query

import (
	"database/sql"
	"fmt"

	"github.com/AubreeH/goApiDb/database"
)

func (q *Query[T]) Sql(pretty bool) (string, error) {
	return q.format(pretty)
}

func (q *Query[T]) First(db *database.Database) (T, error) {
	var out T
	reset := tempSet(&q.limit.limit, 1)
	defer reset()

	rs, err := q.execQuery(db)
	if err != nil {
		return out, err
	}
	defer rs.Close()

	return scanRow[T](rs)
}

func (q *Query[T]) FirstN(db *database.Database, n uint) ([]T, error) {
	var out []T

	if n == 0 {
		return out, nil
	}

	reset := tempSet(&q.limit.limit, n)
	defer reset()

	return q.All(db)
}

func (q *Query[T]) All(db *database.Database) ([]T, error) {
	rs, err := q.execQuery(db)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	return scanRows[T](rs)
}

func (q *Query[T]) Paginated(db *database.Database, itemsPerPage, page uint) ([]T, error) {
	var out []T

	if itemsPerPage == 0 {
		return out, nil
	}

	resetLimit := tempSet(&q.limit.limit, itemsPerPage)
	defer resetLimit()

	resetOffset := tempSet(&q.limit.offset, itemsPerPage*page)
	defer resetOffset()

	return q.All(db)
}

func (q *Query[T]) execQuery(db *database.Database) (*sql.Rows, error) {
	queryStr, err := q.format(false)
	if err != nil {
		return nil, err
	}

	fmt.Println(queryStr)

	parsedQuery, queryArgs := q.params.parse(queryStr)
	return db.Db.Query(parsedQuery, queryArgs...)
}
