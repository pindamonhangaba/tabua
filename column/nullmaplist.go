package column

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullMapList will marshal to null if null. Blank string input will be considered null.
type NullMapList struct {
	Valid bool
	List  []map[string]string
}

// MapListFrom creates a new String that will never be blank.
func MapListFrom(s []map[string]string) NullMapList {
	return NewMapList(s, true)
}

// NewMapList creates a new StringList
func NewMapList(s []map[string]string, valid bool) NullMapList {
	return NullMapList{
		List:  s,
		Valid: valid,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *NullMapList) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.List)
	if err == nil {
		s.Valid = true
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (s NullMapList) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("[]"), nil
	}
	return json.Marshal(s.List)
}

// MarshalText implements encoding.TextMarshaler.
func (s NullMapList) MarshalText() ([]byte, error) {
	return s.MarshalJSON()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *NullMapList) UnmarshalText(text []byte) error {
	return s.UnmarshalJSON(text)
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *NullMapList) SetValid(v []map[string]string) {
	s.List = v
	s.Valid = true
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s NullMapList) IsZero() bool {
	return !s.Valid
}

// Value implements the driver Valuer interface.
func (s NullMapList) Value() (driver.Value, error) {
	b, err := json.Marshal(s.List)
	return driver.Value(b), err
}

// Scan implements the Scanner interface.
func (s *NullMapList) Scan(src interface{}) error {
	var source []byte
	// let's support string and []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	case nil:
		s.Valid = false
		return nil
	default:
		return errors.New("Incompatible type for NullMapList")
	}
	err := json.Unmarshal(source, &s.List)
	if err == nil {
		s.Valid = true
	}
	return err
}
