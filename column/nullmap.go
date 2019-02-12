package column

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullStringMap will marshal to null if null.
// Blank string input will be considered null.
type NullStringMap struct {
	Valid bool
	Map   map[string]string
}

// MapFrom creates a new Map that will never be blank.
func MapFrom(m map[string]string) NullStringMap {
	return NewStringMap(m, true)
}

// NewStringMap creates a new StringMap
func NewStringMap(m map[string]string, valid bool) NullStringMap {
	return NullStringMap{
		Map:   m,
		Valid: valid,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *NullStringMap) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &m.Map)
	if err == nil {
		m.Valid = true
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (m NullStringMap) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("{}"), nil
	}
	return json.Marshal(m.Map)
}

// MarshalText implements encoding.TextMarshaler.
func (m NullStringMap) MarshalText() ([]byte, error) {
	return m.MarshalJSON()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (m *NullStringMap) UnmarshalText(text []byte) error {
	return m.UnmarshalJSON(text)
}

// SetValid changes this Map's value and also sets it to be non-null.
func (m *NullStringMap) SetValid(v map[string]string) {
	m.Map = v
	m.Valid = true
}

// IsZero returns true for null maps, for potential future omitempty support.
func (m NullStringMap) IsZero() bool {
	return !m.Valid
}

// Value implements the driver Valuer interface.
func (m NullStringMap) Value() (driver.Value, error) {
	b, err := json.Marshal(m.Map)
	return driver.Value(b), err
}

// Scan implements the Scanner interface.
func (m *NullStringMap) Scan(src interface{}) error {
	var source []byte
	// let's support string and []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	case nil:
		m.Valid = false
		return nil
	default:
		return errors.New("Incompatible type for NullStringMap")
	}
	err := json.Unmarshal(source, &m.Map)
	if err == nil {
		m.Valid = true
	}
	return err
}
