package query

func (query *Query) SetParam(name string, param any) *Query {
	query.params[name] = parameter{Name: name, Value: param}
	return query
}
