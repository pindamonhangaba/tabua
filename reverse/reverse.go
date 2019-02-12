package reverse

import (
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

// Table represents a database table
type Table struct {
	Name        string       `json:"name"`
	Columns     []Column     `json:"columns"`
	Constraints []Constraint `json:"constraints"`
	Comment     *string      `json:"comment"`
}

// Column represents a database column
type Column struct {
	Name      string  `json:"name"`
	UDTName   string  `json:"udt_name"`
	NonNull   bool    `json:"non_null"`
	DataType  string  `json:"data_type"`
	Comment   *string `json:"comment"`
	Dimension int32   `json:"dimension"`
}

// Constraint represents a database constraint
type Constraint struct {
	Name           string             `json:"name"`
	Definition     string             `json:"definition"`
	Type           string             `json:"type"`
	ColumnsLocal   []ConstraintColumn `json:"columns_local"`
	ColumnsForeign []ConstraintColumn `json:"columns_foreign"`
}

// ConstraintColumn represents a database column constraint definition
type ConstraintColumn struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

// Filter holds schema and tables to filter results by
type Filter struct {
	Schema string
	Tables []string
}

// SQLGenerator generates an SQL query to reverse database -> structs
type SQLGenerator func(Filter) (string, []interface{}, error)

// Reverser reverses a database to its Table(s)
type Reverser struct {
	DB     *sqlx.DB
	GetSQL SQLGenerator
}

// Run starts the reversing process
func (r *Reverser) Run(f Filter) (t []Table, err error) {
	qSQL, args, err := r.GetSQL(f)
	if err != nil {
		return nil, err
	}
	res := struct {
		Tables *types.JSONText `db:"tables"`
	}{}
	err = r.DB.Get(&res, qSQL, args...)
	if err != nil {
		return nil, err
	}
	if res.Tables == nil {
		return nil, NoTablesErr(errors.New("No result from query"))
	}
	err = json.Unmarshal((*res.Tables), &t)
	return t, err
}

// New creates a new Reverser
func New(db *sqlx.DB, r SQLGenerator) (*Reverser, error) {
	return &Reverser{GetSQL: r, DB: db}, nil
}

// NoTablesErr is the error returned when no results return from the query
type NoTablesErr error
