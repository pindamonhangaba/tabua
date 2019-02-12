package op

import (
	tab "github.com/pindamonhangaba/tabua"
)

// IList returns Columns as an interface{} list
func IList(cols ...tab.Column) (i []interface{}) {
	for _, col := range cols {
		i = append(i, col)
	}
	return i
}

// HasColumn checks and returns a matching Column in cols
func HasColumn(cols []tab.Column, col tab.Column) (*tab.Column, bool) {
	for i, c := range cols {
		if c.Name() == col.Name() {
			return &cols[i], true
		}
	}
	return nil, false
}

// Exclude returns Columns except the excluded
func Exclude(cols []tab.Column, exclude ...tab.Column) (i []tab.Column) {
	for _, col := range cols {
		_, has := HasColumn(exclude, col)
		if !has {
			i = append(i, col)
		}
	}
	return i
}

// InOrder returns Columns except the excluded
func InOrder(cols []tab.Column, order []tab.Column) (i []tab.Column) {
	oMap := map[string]tab.Column{}

	for _, col := range cols {
		oMap[col.Name()] = col
	}

	for _, col := range order {
		oc, has := oMap[col.Name()]
		if has {
			i = append(i, oc)
		}
	}
	return i
}

// NullFill returns subset's missing Columns as NullCol's
func NullFill(set []tab.Column, subset []tab.Column) (i []tab.Column) {
	oMap := map[string]tab.Column{}

	for _, col := range subset {
		oMap[col.Name()] = col
	}

	for _, col := range set {
		oc, has := oMap[col.Name()]
		if has {
			i = append(i, oc)
		} else {
			i = append(i, tab.NullCol(col))
		}
	}
	return i
}
