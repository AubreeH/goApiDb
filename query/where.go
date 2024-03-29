package query

import "strings"

type clauseType string

const (
	andClause clauseType = "AND"
	orClause  clauseType = "OR"
)

type Where struct {
	clauseType clauseType
	wheres     []Where
	statements []string
}

type WhereOperation func(w *Where, state func(statements ...string))

func (q *Query[T]) WhereBuilder(function func(w *Where)) *Query[T] {
	function(&q.where)
	return q
}

func (q *Query[T]) Where(statements ...string) *Query[T] {
	q.where.clauseType = andClause
	q.where.statements = append(q.where.statements, statements...)
	return q
}

func (parent *Where) And(function WhereOperation) {
	parent.do(andClause, function)
}

func (parent *Where) Or(function WhereOperation) {
	parent.do(orClause, function)
}

func (parent *Where) do(t clauseType, function WhereOperation) {
	where := Where{
		clauseType: t,
	}
	state := func(statements ...string) {
		where.statements = append(where.statements, statements...)
	}
	function(&where, state)
	parent.wheres = append(parent.wheres, where)
}

func (w Where) format(pretty bool) string {
	var out []string
	out = append(out, w.statements...)
	for _, where := range w.wheres {
		out = append(out, where.format(pretty))
	}

	val := strings.Join(out, ") "+string(w.clauseType)+"(")
	if len(out) > 0 {
		return "(" + val + ")"
	}
	return ""
}
