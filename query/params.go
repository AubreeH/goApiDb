package query

import "regexp"

type params map[string]any

var paramRegex = regexp.MustCompile(":([a-zA-Z0-9]+)")

func (p params) parse(q string) (string, []any) {
	matches := paramRegex.FindAllStringSubmatch(q, -1)
	if matches == nil {
		return q, nil
	}

	out := make([]any, len(matches))
	for i, match := range matches {
		out[i] = p[match[1]]
	}

	return paramRegex.ReplaceAllString(q, "?"), out
}

func (q *Query[T]) SetParam(key string, value any) *Query[T] {
	if q.params == nil {
		q.params = make(params)
	}

	q.params[key] = value
	return q
}
