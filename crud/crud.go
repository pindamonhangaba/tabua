package crud

import (
	sq "github.com/elgris/sqrl"
	tab "github.com/pindamonhangaba/tabua"
	"github.com/pindamonhangaba/tabua/op"
	"strings"
)

type qSB struct {
	sq.StatementBuilderType
}

// Q return an sq.StatementBuilderType with sq.Question (?) as placeholder
var Q qSB

// Default placeholder is sq.Dollar ($#)
var builder sq.StatementBuilderType

func init() {
	Q = qSB{sq.StatementBuilder.PlaceholderFormat(sq.Question)}
	builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (qs qSB) Select(cols []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	return selectOnly(qs.StatementBuilderType, cols, conditions...)
}

func (qs qSB) Insert(cols []tab.Column) (q tab.Query, err error) {
	return insertOnly(qs.StatementBuilderType, cols)
}

func (qs qSB) Upsert(cols []tab.Column, conflict tab.Column, update []tab.Column) (q tab.Query, err error) {
	return upsertOnly(qs.StatementBuilderType, cols, conflict, update)
}

func (qs qSB) InsertR(cols []tab.Column, returning ...tab.Column) (q tab.Query, err error) {
	return insertOnlyR(qs.StatementBuilderType, cols, returning...)
}

func (qs qSB) Update(cols []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	return update(qs.StatementBuilderType, cols, conditions...)
}

func (qs qSB) Delete(t tab.Table, conditions ...tab.Column) (q tab.Query, err error) {
	return delete(qs.StatementBuilderType, t, conditions...)
}

// Select generates query to select only cols columns, given conditions
// The table is defined by the first column
func Select(cols []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	return selectOnly(builder, cols, conditions...)
}

// Insert generates query to insert table
func Insert(cols []tab.Column) (q tab.Query, err error) {
	return insertOnly(builder, cols)
}

// Upsert generates query to upsert table
func Upsert(cols []tab.Column, onConflict tab.Column, update []tab.Column) (q tab.Query, err error) {
	return upsertOnly(builder, cols, onConflict, update)
}

// InsertR generates query to insert columns while returning selected columns
// The table is defined by the first column
func InsertR(cols []tab.Column, returning ...tab.Column) (q tab.Query, err error) {
	return insertOnlyR(builder, cols, returning...)
}

// Update generates query to update a table given conditions
func Update(cols []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	return update(builder, cols, conditions...)
}

// UpdateR generates query to update a table given conditions
func UpdateR(cols []tab.Column, returning []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	return updateR(builder, cols, returning, conditions...)
}

// Delete generates query to remove entry given conditions
func Delete(t tab.Table, conditions ...tab.Column) (q tab.Query, err error) {
	return delete(builder, t, conditions...)
}

// DeleteR generates query to remove entry given conditions and return selected columns
func DeleteR(t tab.Table, returning []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	return deleteR(builder, t, returning, conditions...)
}

func selectOnly(builder sq.StatementBuilderType, cols []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	if len(cols) < 1 {
		return nil, tab.QueryGenerationError{"No columns to select"}
	}
	t := cols[0].Table()

	var columns []string
	for _, c := range cols {
		columns = append(columns, op.Q(c))
	}
	stmt := builder.Select(columns...).
		From(op.Q(t))

	if len(conditions) > 0 {
		where := sq.Eq{}
		for _, cond := range conditions {
			where[op.Q(cond)] = cond
		}
		stmt = stmt.Where(where)
	}

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func insertOnly(builder sq.StatementBuilderType, cols []tab.Column, returning ...tab.Column) (q tab.Query, err error) {
	if len(cols) < 1 {
		return nil, tab.QueryGenerationError{"No columns to insert"}
	}
	t := cols[0].Table()

	var columns []string
	var values []interface{}
	for _, v := range cols {
		columns = append(columns, op.Q(v))
		values = append(values, v)
	}
	stmt := builder.Insert(op.Q(t)).
		Columns(columns...).
		Values(values...)

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func upsertOnly(builder sq.StatementBuilderType, cols []tab.Column, onConflict tab.Column, update []tab.Column, returning ...tab.Column) (q tab.Query, err error) {
	if len(cols) < 1 {
		return nil, tab.QueryGenerationError{"No columns to insert"}
	}
	t := cols[0].Table()

	var columns []string
	var values []interface{}
	for _, v := range cols {
		columns = append(columns, op.Q(v))
		values = append(values, v)
	}
	stmt := builder.Insert(op.Q(t)).
		Columns(columns...).
		Values(values...)

	cft := onConflict.Name()
	upd := []string{}
	for _, col := range update {
		upd = append(upd, col.Name()+"= EXCLUDED."+col.Name())
	}

	stmtConf := "ON CONFLICT (" + cft + ") DO NOTHING"
	if len(update) > 0 {
		stmtConf = " ON CONFLICT (" + cft + ") DO UPDATE SET " +
			strings.Join(upd, ", ")
	}
	stmt = stmt.Suffix(stmtConf)

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func insertOnlyR(builder sq.StatementBuilderType, cols []tab.Column, returning ...tab.Column) (q tab.Query, err error) {
	if len(cols) < 1 {
		return nil, tab.QueryGenerationError{"No columns to insert"}
	}
	t := cols[0].Table()

	var columns []string
	var values []interface{}
	for _, v := range cols {
		columns = append(columns, op.Q(v))
		values = append(values, v)
	}
	stmt := builder.Insert(op.Q(t)).
		Columns(columns...).
		Values(values...)

	if len(returning) > 0 {
		rtn := ""
		for _, r := range returning {
			rtn += op.Q(r) + ","
		}
		rtnR := []rune(rtn)
		rtnR = rtnR[0 : len(rtnR)-1]
		stmt = stmt.Suffix("RETURNING " + string(rtnR))
	}

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func update(builder sq.StatementBuilderType, columns []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	if len(columns) < 1 {
		return nil, tab.QueryGenerationError{"No columns to select"}
	}
	t := columns[0].Table()

	stmt := builder.Update(op.Q(t))

	for _, v := range columns {
		stmt = stmt.Set(op.Q(v), v)
	}

	if len(conditions) > 0 {
		where := sq.Eq{}
		for _, cond := range conditions {
			where[op.Q(cond)] = cond
		}
		stmt = stmt.Where(where)
	}

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func updateR(builder sq.StatementBuilderType, columns []tab.Column, returning []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {
	if len(columns) < 1 {
		return nil, tab.QueryGenerationError{"No columns to select"}
	}
	t := columns[0].Table()

	stmt := builder.Update(op.Q(t))

	for _, v := range columns {
		stmt = stmt.Set(op.Q(v), v)
	}

	if len(conditions) > 0 {
		where := sq.Eq{}
		for _, cond := range conditions {
			where[op.Q(cond)] = cond
		}
		stmt = stmt.Where(where)
	}

	if len(returning) > 0 {
		rtn := ""
		for _, r := range returning {
			rtn += op.Q(r) + ","
		}
		rtnR := []rune(rtn)
		rtnR = rtnR[0 : len(rtnR)-1]
		stmt = stmt.Suffix("RETURNING " + string(rtnR))
	}

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func delete(builder sq.StatementBuilderType, t tab.Table, conditions ...tab.Column) (q tab.Query, err error) {

	stmt := builder.Delete(op.Q(t))

	if len(conditions) > 0 {
		where := sq.Eq{}
		for _, cond := range conditions {
			where[op.Q(cond)] = cond
		}
		stmt = stmt.Where(where)
	}

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}

func deleteR(builder sq.StatementBuilderType, t tab.Table, returning []tab.Column, conditions ...tab.Column) (q tab.Query, err error) {

	stmt := builder.Delete(op.Q(t))

	if len(conditions) > 0 {
		where := sq.Eq{}
		for _, cond := range conditions {
			where[op.Q(cond)] = cond
		}
		stmt = stmt.Where(where)
	}

	if len(returning) > 0 {
		rtn := ""
		for _, r := range returning {
			rtn += op.Q(r) + ","
		}
		rtnR := []rune(rtn)
		rtnR = rtnR[0 : len(rtnR)-1]
		stmt = stmt.Suffix("RETURNING " + string(rtnR))
	}

	sql, args, err := stmt.ToSql()

	return tab.Sq{sql, args}, err
}
