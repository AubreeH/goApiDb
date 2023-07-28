package query

import "errors"

func (query *Query) validateAlias(alias string) *Query {
	if alias == "" {
		panic("empty alias provided")
	}

	if _, exists := query.tables[alias]; exists {
		panic("alias already in use")
	}
	return query
}

func (query *Query) validateQuery() error {
	if query.selectStr == "" {
		return errors.New("select statement not provided")
	}

	if query.from == "" {
		return errors.New("from table not provided")
	}
	return nil
}
