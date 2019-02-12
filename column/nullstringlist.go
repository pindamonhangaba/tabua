package column

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullStringList marshal to empty list if null.
type NullStringList struct {
	Valid bool
	List  []string
}

// ListFrom creates a new string list.
func ListFrom(s []string) NullStringList {
	return NewStringList(s, true)
}

// NewStringList creates a new StringList
func NewStringList(s []string, valid bool) NullStringList {
	return NullStringList{
		List:  s,
		Valid: valid,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *NullStringList) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.List)
	if err == nil {
		s.Valid = true
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (s NullStringList) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("[]"), nil
	}
	return json.Marshal(s.List)
}

// MarshalText implements encoding.TextMarshaler.
func (s NullStringList) MarshalText() ([]byte, error) {
	return s.MarshalJSON()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *NullStringList) UnmarshalText(text []byte) error {
	return s.UnmarshalJSON(text)
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *NullStringList) SetValid(v []string) {
	s.List = v
	s.Valid = true
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s NullStringList) IsZero() bool {
	return !s.Valid
}

// Value implements the driver Valuer interface.
func (i NullStringList) Value() (driver.Value, error) {
	b, err := json.Marshal(i.List)
	return driver.Value(b), err
}

// Scan implements the Scanner interface.
func (i *NullStringList) Scan(src interface{}) error {
	var source []byte
	// let's support string and []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	case nil:
		i.Valid = false
		return nil
	default:
		return errors.New("Incompatible type for NullStringList")
	}
	err := json.Unmarshal(source, &i.List)
	if err == nil {
		i.Valid = true
	}
	return err
}
