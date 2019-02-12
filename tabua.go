package tabua

// ValueQuotes sets the quote char to use when quoting values
var ValueQuotes = `'`

// NameQuotes sets the quote char to use when quoting table and columnnames
var NameQuotes = `"`

type ConstraintType string

// PSQL Constraint types
const (
	ConstraintUnique ConstraintType = "UNIQUE"
	ConstraintCheck  ConstraintType = "CHECK"
	ConstraintFK     ConstraintType = "FOREIGN KEY"
	ConstraintPK     ConstraintType = "PRIMARY KEY"
)

// Constraint returns the table constraint name
type Constraint string

// Constrainer returns the constraint's type
type Constrainer interface {
	Type() ConstraintType
	Definition() string
}

// UniqueConstrainer limits Column to unique values
type UniqueConstrainer interface {
	Constrainer
	Uniques() []Column
}

// CheckConstrainer SQL check attributed to Column
type CheckConstrainer interface {
	Constrainer
	Columns() []Column
}

// PKConstrainer table primary key
type PKConstrainer interface {
	Constrainer
	Keys() []Column
}

// FKConstrainer table foreign keys
type FKConstrainer interface {
	Constrainer
	Key() FK
}

// Namer describes an objects name
type Namer interface {
	Name() string
}

//FK represents basic information for a database foreignkey entry
type FK struct {
	From []Column
	To   []Column
}

// Column descrives an SQL column
type Column interface {
	Namer
	Table() Table
	SQLType() string
	NonNull() bool
}

// Table describes an SQL table
type Table interface {
	Namer
	Constraints() []Constraint
	Columns() []Column
}

// Querier serves to organize an SQL query and it's arguments
type Querier interface {
	SQL() string
	Args() []interface{}
}
