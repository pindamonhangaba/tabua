package generate

import (
	t "github.com/pindamonhangaba/tabua/generate/types"
	"github.com/pindamonhangaba/tabua/reverse"
	"github.com/serenize/snaker"
	"reflect"
	"strings"
)

func colName(table, column string) string {
	table = snaker.SnakeToCamel(table)
	column = snaker.SnakeToCamel(column)
	if table == column {
		return column + "Col"
	}
	return column
}
func columnName(tab, col string) string { return colName(tab, col) }
func cstname(s string) string {
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, "-", "", -1)
	return "CS" + snaker.SnakeToCamel(s)
}
func camel(s string) string { return snaker.SnakeToCamel(s) }
func camelLower(s string) string {
	n := snaker.SnakeToCamel(s)
	return strings.ToLower(string(n[0])) + n[1:]
}
func packageFilename(s string) string { return strings.Replace(s, "_", "", -1) }
func reType(col reverse.Column, nnull bool) reflect.Type {

	udt, ok := t.SQLType2Type(t.SQLType{Name: col.UDTName, Dimension: int(col.Dimension)}, nnull)
	if !ok {
		daty, _ := t.SQLType2Type(t.SQLType{Name: col.DataType}, nnull)
		return daty
	}
	return udt
}
func isJSON(s string) bool {
	tps := strings.Split(s, ", ")
	return tps[0] == t.JSON || tps[0] == t.JSONB
}
func isString(s string) bool {
	tps := strings.Split(s, ", ")
	return !(tps[0] == t.JSON || tps[0] == t.JSONB) && t.SQLTypes[tps[0]] == t.TextSQLType
}
func isTime(s string) bool {
	tps := strings.Split(s, ", ")
	return t.SQLTypes[tps[0]] == t.TimeSQLType
}
func isBlob(s string) bool {
	tps := strings.Split(s, ", ")
	return t.SQLTypes[tps[0]] == t.BlobSQLType
}
