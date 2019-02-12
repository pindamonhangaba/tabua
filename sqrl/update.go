package sqrl

import (
	sq "github.com/elgris/sqrl"
	t "github.com/pindamonhangaba/tabua"
)

// UpdateBuilder builds SQL UPDATE statements.
type UpdateBuilder struct {
	*sq.UpdateBuilder
}

// Table sets the table to be updateb.
func (b *UpdateBuilder) Table(table t.Table) *UpdateBuilder {
	b.UpdateBuilder = b.UpdateBuilder.Table(t.Q(table))
	return b
}

// Set adds SET clauses to the query.
func (b *UpdateBuilder) Set(cols ...t.Column) *UpdateBuilder {
	for _, c := range cols {
		b.UpdateBuilder = b.UpdateBuilder.Set(t.Q(c), c)
	}
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b *UpdateBuilder) Where(cols ...t.Column) *UpdateBuilder {
	where := sq.Eq{}
	for _, c := range cols {
		where[t.Q(c)] = c
	}
	b.UpdateBuilder = b.UpdateBuilder.Where(where)
	return b
}
