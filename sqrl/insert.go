package sqrl

import (
	sq "github.com/elgris/sqrl"
	t "github.com/pindamonhangaba/tabua"
)

// InsertBuilder builds SQL INSERT statements.
type InsertBuilder struct {
	*sq.InsertBuilder
}

// Into sets the INTO clause of the query.
func (b *InsertBuilder) Into(table t.Table) *InsertBuilder {
	b.InsertBuilder = b.InsertBuilder.Into(t.Q(table))
	return b
}

// Columns adds insert columns to the query.
func (b *InsertBuilder) Columns(cols ...t.Column) *InsertBuilder {
	b.InsertBuilder = b.InsertBuilder.Columns(t.Columns(cols...)...)
	return b
}
