package sqrl

import (
	"strings"

	sq "github.com/elgris/sqrl"
	t "github.com/pindamonhangaba/tabua"
)

// SelectBuilder builds SQL SELECT statements.
type SelectBuilder struct {
	*sq.SelectBuilder
}

// Columns adds result columns to the query.
func (b *SelectBuilder) Columns(cols ...t.Column) *SelectBuilder {
	b.SelectBuilder = b.SelectBuilder.Columns(t.Columns(cols...)...)
	return b
}

// From sets the FROM clause of the query.
func (b *SelectBuilder) From(table t.Table) *SelectBuilder {
	b.SelectBuilder = b.SelectBuilder.From(t.Q(table))
	return b
}

// Join adds a JOIN clause to the query.
func (b *SelectBuilder) Join(cols ...t.Column) *SelectBuilder {
	b.SelectBuilder = b.SelectBuilder.JoinClause("JOIN " + columnsToJoin(cols...))
	return b
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b *SelectBuilder) LeftJoin(cols ...t.Column) *SelectBuilder {
	b.SelectBuilder = b.SelectBuilder.JoinClause("LEFT JOIN " + columnsToJoin(cols...))
	return b
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b *SelectBuilder) RightJoin(cols ...t.Column) *SelectBuilder {
	b.SelectBuilder = b.SelectBuilder.JoinClause("RIGHT JOIN " + columnsToJoin(cols...))
	return b
}

// Where adds an expression to the WHERE clause of the query.
func (b *SelectBuilder) Where(cols ...t.Column) *SelectBuilder {
	where := sq.Eq{}
	for _, c := range cols {
		where[t.Q(c)] = c
	}
	b.SelectBuilder = b.SelectBuilder.Where(where)
	return b
}

func columnsToJoin(cols ...t.Column) string {
	flatCols := []string{}
	tableName := ""
	for _, c := range cols {
		if c2, ok := c.FK(); ok {
			t1name := t.Q(c.Table())
			c1name := t.Q(c)
			t2name := t.Q(c2.Table())
			c2name := t.Q(c2)
			flatCols = append(flatCols, t1name+"."+c1name+"="+t2name+"."+c2name)
			tableName = t1name
		}
	}
	return tableName + " ON " + strings.Join(flatCols, " AND ")
}
