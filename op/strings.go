package op

import (
	tab "github.com/pindamonhangaba/tabua"
	"strings"
)

//Q encases a Named object in quotes
func Q(n tab.Named) string {
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
	flatCols := []string{}
	for _, c := range cols {
		tname := Q(c.Table())
		cname := Q(c)
		flatCols = append(flatCols, tname+"."+cname)
	}
	return flatCols
}
