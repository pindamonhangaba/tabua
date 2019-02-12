package reverse

import (
	"strconv"
	"strings"
)

// SQLFromPsql returns a query to reverse tables to structs
func SQLFromPsql(f Filter) (string, []interface{}, error) {
	args := []interface{}{f.Schema}
	filterArgs := []string{}
	for _, n := range f.Tables {
		count := len(args) + 1
		filterArgs = append(filterArgs, "$"+strconv.Itoa(count))
		args = append(args, n)
	}

	tableFilter := ""
	if len(filterArgs) > 0 {
		tableFilter = "WHERE table_name IN (" + strings.Join(filterArgs, ",") + ")"
	}

	query := `
	with
	cols as (
	select attrelid, attnum, json_build_object('table',relname, 'column',attname) as col, attndims as dimension from pg_attribute
	join pg_class on attrelid = oid
	),
	loc as (
		select conrelid as attrelid, unnest(conkey) as attnum, conname, contype from pg_constraint
	),
	local_cols as (
		select conname, contype, json_agg(col) as columns_local from loc
		join cols using(attrelid, attnum)
		group by conname, contype
	),
	fog as (
		select confrelid as attrelid, unnest(confkey) as attnum, conname, contype from pg_constraint
	),
	foreign_cols as  (
		select conname, contype, json_agg(col) as columns_foreign from fog
		join cols using(attrelid, attnum)
		group by conname, contype
	),
	constrs as (
		select relname, conname, columns_local, columns_foreign, pc.contype, pg_get_constraintdef(pc.oid) as constraintdef	from pg_constraint pc 
		join pg_class pt on pc.conrelid = pt.oid
		left join local_cols using(conname, contype)
		left join foreign_cols using(conname, contype)
		join pg_namespace n ON n.oid = pc.connamespace
		where n.nspname = $1
	),
	table_constraints as (
		SELECT relname as table_name,
			json_agg(json_build_object(
				'name', conname,
				'definition', constraintdef,
				'type', CASE
					WHEN contype = 'f' THEN 'FOREIGN KEY'
					WHEN contype = 'p' THEN 'PRIMARY KEY'
					WHEN contype = 'c' THEN 'CHECK'
					WHEN contype = 'u' THEN 'UNIQUE'
				END ,
				'columns_local', columns_local,
				'columns_foreign', columns_foreign
			)) as table_constraints
		from constrs
		GROUP BY table_name
	),
	column_comments as (
		SELECT n.nspname  as table_schema, c.relname As table_name, a.attname As column_name,  d.description as comment
		FROM pg_class As c
		INNER JOIN pg_attribute As a ON c.oid = a.attrelid
		LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_tablespace t ON t.oid = c.reltablespace
		LEFT JOIN pg_description As d ON (d.objoid = c.oid AND d.objsubid = a.attnum)
		WHERE  c.relkind IN('r', 'v') AND d.description is not null AND n.nspname = $1
		ORDER BY n.nspname, c.relname, a.attname
	),
	table_comments as (
		select * from (select table_name, (select obj_description(TABLE_NAME::regclass)) as comment from table_constraints) ok
		where comment is not null
	),
	columns_list AS (
		SELECT
			table_name, json_agg(json_build_object('name', COLUMN_NAME, 'udt_name', udt_name, 'non_null', NOT CAST (is_nullable AS BOOLEAN), 'data_type', data_type,'comment',cc.comment, 'dimension', dimension )) as columns
		FROM
		INFORMATION_SCHEMA. COLUMNS incol
		LEFT JOIN column_comments cc using(TABLE_SCHEMA, table_name, COLUMN_NAME)
		LEFT JOIN cols ON cols.col->>'column' = incol.COLUMN_NAME and cols.col->>'table' = incol.TABLE_NAME
		WHERE  TABLE_SCHEMA = 'public'
		GROUP BY table_name
	),
	all_tables as (
		select json_build_object('name',table_name, 'columns',columns, 'constraints',table_constraints, 'comment',comment) as table from table_constraints
		left join columns_list using(table_name)
		left join table_comments using(table_name)
		` + tableFilter + `
	)

	select json_agg(all_tables.table) as tables from all_tables
		`

	return query, args, nil
}
