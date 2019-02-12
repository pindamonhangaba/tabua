package types

import (
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
	"reflect"
	"sort"
	"strings"
	"time"
)

// use xorm/core/type.go for type conversion between go and db

// Database types
// (unused)
const (
	POSTGRES = "postgres"
	SQLITE   = "sqlite3"
	MYSQL    = "mysql"
	MSSQL    = "mssql"
	ORACLE   = "oracle"
)

// SQLType is a cpy from xorm SQL types
type SQLType struct {
	Name           string
	DefaultLength  int
	DefaultLength2 int
	Dimension      int
}

// SQL Types
const (
	UnknownSQLType = iota
	TextSQLType
	BlobSQLType
	TimeSQLType
	NumericSQLType
)

// IsType check the type of an SQLType
func (s *SQLType) IsType(st int) bool {
	if t, ok := SQLTypes[s.Name]; ok && t == st {
		return true
	}
	return false
}

// IsText checks if the SQLType is of a TextSQLType
func (s *SQLType) IsText() bool {
	return s.IsType(TextSQLType)
}

// IsBlob checks if the SQLType is of a BlobSQLType
func (s *SQLType) IsBlob() bool {
	return s.IsType(BlobSQLType)
}

// IsTime checks if the SQLType is of a TimeType
func (s *SQLType) IsTime() bool {
	return s.IsType(TimeSQLType)
}

// IsNumeric checks if the SQLType is of a NumericSQLType
func (s *SQLType) IsNumeric() bool {
	return s.IsType(NumericSQLType)
}

// IsJSON checks if the SQLType is JSON(B)
func (s *SQLType) IsJSON() bool {
	return s.Name == JSON || s.Name == JSONB
}

// SQL types
var (
	Bit       = "BIT"
	TinyInt   = "TINYINT"
	SmallInt  = "SMALLINT"
	MediumInt = "MEDIUMINT"
	Int       = "INT"
	Integer   = "INTEGER"
	BigInt    = "BIGINT"

	Enum = "ENUM"
	Set  = "SET"

	Char       = "CHAR"
	Varchar    = "VARCHAR"
	NVarchar   = "NVARCHAR"
	TinyText   = "TINYTEXT"
	Text       = "TEXT"
	Clob       = "CLOB"
	MediumText = "MEDIUMTEXT"
	LongText   = "LONGTEXT"
	UUID       = "UUID"

	Date       = "DATE"
	DateTime   = "DATETIME"
	Time       = "TIME"
	TimeStamp  = "TIMESTAMP"
	TimeStampz = "TIMESTAMPZ"

	Decimal = "DECIMAL"
	Numeric = "NUMERIC"

	Real   = "REAL"
	Float  = "FLOAT"
	Double = "DOUBLE"

	Binary     = "BINARY"
	VarBinary  = "VARBINARY"
	TinyBlob   = "TINYBLOB"
	Blob       = "BLOB"
	MediumBlob = "MEDIUMBLOB"
	LongBlob   = "LONGBLOB"
	Bytea      = "BYTEA"

	Bool = "BOOL"

	Serial    = "SERIAL"
	BigSerial = "BIGSERIAL"

	JSON  = "JSON"
	JSONB = "JSONB"

	SQLTypes = map[string]int{
		Bit:       NumericSQLType,
		TinyInt:   NumericSQLType,
		SmallInt:  NumericSQLType,
		MediumInt: NumericSQLType,
		Int:       NumericSQLType,
		Integer:   NumericSQLType,
		BigInt:    NumericSQLType,

		Enum:  TextSQLType,
		Set:   TextSQLType,
		JSON:  TextSQLType,
		JSONB: TextSQLType,

		Char:       TextSQLType,
		Varchar:    TextSQLType,
		NVarchar:   TextSQLType,
		TinyText:   TextSQLType,
		Text:       TextSQLType,
		MediumText: TextSQLType,
		LongText:   TextSQLType,
		UUID:       TextSQLType,
		Clob:       TextSQLType,

		Date:       TimeSQLType,
		DateTime:   TimeSQLType,
		Time:       TimeSQLType,
		TimeStamp:  TimeSQLType,
		TimeStampz: TimeSQLType,

		Decimal: NumericSQLType,
		Numeric: NumericSQLType,
		Real:    NumericSQLType,
		Float:   NumericSQLType,
		Double:  NumericSQLType,

		Binary:    BlobSQLType,
		VarBinary: BlobSQLType,

		TinyBlob:   BlobSQLType,
		Blob:       BlobSQLType,
		MediumBlob: BlobSQLType,
		LongBlob:   BlobSQLType,
		Bytea:      BlobSQLType,

		Bool: NumericSQLType,

		Serial:    NumericSQLType,
		BigSerial: NumericSQLType,
	}

	intTypes  = sort.StringSlice{"*int", "*int16", "*int32", "*int8"}
	uintTypes = sort.StringSlice{"*uint", "*uint16", "*uint32", "*uint8"}
)

// !nashtsai! treat following var as interal const values, these are used for reflect.TypeOf comparision
var (
	cEmptyString       string
	cBoolDefault       bool
	cByteDefault       byte
	cComplex64Default  complex64
	cComplex128Default complex128
	cFloat32Default    float32
	cFloat64Default    float64
	cInt64Default      int64
	cUInt64Default     uint64
	cInt32Default      int32
	cUInt32Default     uint32
	cInt16Default      int16
	cUInt16Default     uint16
	cInt8Default       int8
	cUInt8Default      uint8
	cIntDefault        int
	cUIntDefault       uint
	cTimeDefault       time.Time
)

// Golang types
var (
	IntType   = reflect.TypeOf(cIntDefault)
	Int8Type  = reflect.TypeOf(cInt8Default)
	Int16Type = reflect.TypeOf(cInt16Default)
	Int32Type = reflect.TypeOf(cInt32Default)
	Int64Type = reflect.TypeOf(cInt64Default)

	UintType   = reflect.TypeOf(cUIntDefault)
	Uint8Type  = reflect.TypeOf(cUInt8Default)
	Uint16Type = reflect.TypeOf(cUInt16Default)
	Uint32Type = reflect.TypeOf(cUInt32Default)
	Uint64Type = reflect.TypeOf(cUInt64Default)

	Float32Type = reflect.TypeOf(cFloat32Default)
	Float64Type = reflect.TypeOf(cFloat64Default)

	Complex64Type  = reflect.TypeOf(cComplex64Default)
	Complex128Type = reflect.TypeOf(cComplex128Default)

	StringType = reflect.TypeOf(cEmptyString)
	BoolType   = reflect.TypeOf(cBoolDefault)
	ByteType   = reflect.TypeOf(cByteDefault)

	TimeType = reflect.TypeOf(cTimeDefault)

	SliceInt8 = reflect.TypeOf([]int8{})
)

// Golang pointer types
var (
	PtrIntType   = reflect.PtrTo(IntType)
	PtrInt8Type  = reflect.PtrTo(Int8Type)
	PtrInt16Type = reflect.PtrTo(Int16Type)
	PtrInt32Type = reflect.PtrTo(Int32Type)
	PtrInt64Type = reflect.PtrTo(Int64Type)

	PtrUintType   = reflect.PtrTo(UintType)
	PtrUint8Type  = reflect.PtrTo(Uint8Type)
	PtrUint16Type = reflect.PtrTo(Uint16Type)
	PtrUint32Type = reflect.PtrTo(Uint32Type)
	PtrUint64Type = reflect.PtrTo(Uint64Type)

	PtrFloat32Type = reflect.PtrTo(Float32Type)
	PtrFloat64Type = reflect.PtrTo(Float64Type)

	PtrComplex64Type  = reflect.PtrTo(Complex64Type)
	PtrComplex128Type = reflect.PtrTo(Complex128Type)

	PtrStringType = reflect.PtrTo(StringType)
	PtrBoolType   = reflect.PtrTo(BoolType)
	PtrByteType   = reflect.PtrTo(ByteType)

	PtrTimeType = reflect.PtrTo(TimeType)
)

// package null types
var (
	NullString = reflect.TypeOf(null.String{})
	NullInt    = reflect.TypeOf(null.Int{})   // 64
	NullFloat  = reflect.TypeOf(null.Float{}) // 64
	NullBool   = reflect.TypeOf(null.Bool{})
	NullTime   = reflect.TypeOf(null.Time{})
)

// package sqlx/types
var (
	JSONText     = reflect.TypeOf(types.JSONText{})
	NullJSONText = reflect.TypeOf(types.NullJSONText{})
)

// SQLType2Type default sql type change to golang types
func SQLType2Type(st SQLType, nnull bool) (t reflect.Type, ok bool) {
	name := strings.ToUpper(st.Name)
	isArray := strings.HasPrefix(name, "_") || st.Dimension > 0
	name = strings.Replace(name, "_", "", -1)
	var res reflect.Type

	switch name {
	case Bit, TinyInt, SmallInt, MediumInt, Int, Integer, Serial:
		if nnull {
			res, ok = IntType, true
		} else {
			res, ok = NullInt, true
		}
	case BigInt, BigSerial:
		if nnull {
			res, ok = Int64Type, true
		} else {
			res, ok = NullInt, true
		}
	case Float, Real:
		if nnull {
			res, ok = Float32Type, true
		} else {
			res, ok = NullFloat, true
		}
	case Double:
		if nnull {
			res, ok = Float64Type, true
		} else {
			res, ok = NullFloat, true
		}
	case Char, Varchar, NVarchar, TinyText, Text, MediumText, LongText, Enum, Set, UUID, Clob:
		if nnull {
			res, ok = StringType, true
		} else {
			res, ok = NullString, true
		}
	case TinyBlob, Blob, LongBlob, Bytea, Binary, MediumBlob, VarBinary:
		r := []byte{}
		if nnull {
			res, ok = reflect.TypeOf(r), true
		} else {
			res, ok = reflect.TypeOf(r), true
		}
	case Bool:
		if nnull {
			res, ok = BoolType, true
		} else {
			res, ok = NullBool, true
		}
	case DateTime, Date, Time, TimeStamp, TimeStampz:
		if nnull {
			res, ok = TimeType, true
		} else {
			res, ok = NullTime, true
		}
	case Decimal, Numeric:
		if nnull {
			res, ok = StringType, true
		} else {
			res, ok = NullString, true
		}
	case JSON, JSONB:
		if nnull {
			res, ok = JSONText, true
		} else {
			res, ok = NullJSONText, true
		}
	default:
		res, ok = SliceInt8, false
	}

	if isArray {
		res = makeDimensions(res, st.Dimension)
	}

	return res, true
}

func makeDimensions(t reflect.Type, d int) reflect.Type {
	if d == 0 {
		return t
	}
	return makeDimensions(reflect.SliceOf(t), d-1)
}
