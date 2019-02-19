package op

import (
	tab "github.com/pindamonhangaba/tabua"
	"strings"
)

//Q encases a Namer object in quotes
func Q(n tab.Namer) string {
	return tab.NameQuotes + n.Name() + tab.NameQuotes
}

//D returns a Column as table.column
func D(c tab.Column) string {
	return Q(c.Table()) + "." + Q(c)
}

// Join returns a comma separated list of quoted column names
func Join(cols ...tab.Column) string {
	flatCols := Columns(cols...)
	return strings.Join(flatCols, ",")
}

// Columns returns an array with column names
func Columns(cols ...tab.Column) []string {
	return columns(false, cols...)
}

// QualifiedColumns returns an array with column names, qualified with their table name
func QualifiedColumns(cols ...tab.Column) []string {
	return columns(true, cols...)
}

func columns(withTableName bool, cols ...tab.Column) []string {
	flatCols := []string{}
	for _, c := range cols {
		tname := Q(c.Table())
		cname := Q(c)
		n := cname
		if withTableName {
			n = tname + "." + cname
		}
		flatCols = append(flatCols, n)
	}
	return flatCols
}
