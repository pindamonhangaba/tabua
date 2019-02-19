package column

import (
	"database/sql/driver"
	tbu "github.com/pindamonhangaba/tabua"
)

// NullCol implements a Column with a nil value
func NullCol(c tbu.Column) tbu.Column {
	return nullCol{c}
}

type nullCol struct {
	origenCol tbu.Column
}

func (n nullCol) Name() string {
	return n.origenCol.Name()
}

func (n nullCol) SQLType() string {
	return n.origenCol.SQLType()
}

func (n nullCol) NonNull() bool {
	return n.origenCol.NonNull()
}

func (n nullCol) Table() tbu.Table {
	return n.origenCol.Table()
}

func (n nullCol) Value() (driver.Value, error) {
	return driver.Value(nil), nil
}

func (n *nullCol) Scan(src interface{}) error {
	return nil
}
