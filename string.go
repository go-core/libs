package libs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// String NullString is null friendly type for string.
type String struct {
	s sql.NullString
}

func NewString(str string) String {
	if str == "" {
		return String{s: sql.NullString{
			String: "",
			Valid:  false,
		}}
	}
	return String{s: sql.NullString{
		String: TrimFormatStr(str),
		Valid:  true,
	}}
}

// NullStringOf return NullString that he value is set.
func NullStringOf(value string) String {
	var s String
	s.Set(value)
	return s
}

// Valid return the value is valid. If true, it is not null value.
func (s *String) Valid() bool {
	return s.s.Valid
}

// StringValue return the value.
func (s *String) StringValue() string {
	return s.s.String
}

// Reset set nil to the value.
func (s *String) Reset() {
	s.s.String = ""
	s.s.Valid = false
}

// Set set the value.
func (s *String) Set(value string) {
	s.s.Valid = true
	s.s.String = value
}

// Scan is a method for database/sql.
func (s *String) Scan(value interface{}) error {
	return s.s.Scan(value)
}

// String return string indicated the value.
func (s String) String() string {
	if !s.s.Valid {
		return ""
	}
	return s.s.String
}

// MarshalJSON encode the value to JSON.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.s.String)
}

// UnmarshalJSON decode data to the value.
func (s *String) UnmarshalJSON(data []byte) error {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	s.s.Valid = value != nil
	if value == nil {
		s.s.String = ""
	} else {
		s.s.String = TrimFormatStr(*value)
	}
	return nil
}

// Value implement driver.Valuer.
func (s String) Value() (driver.Value, error) {
	if !s.Valid() {
		return nil, nil
	}
	return s.s.String, nil
}

func (s String) IsNil() bool {
	return !s.s.Valid || s.s.String == ""
}

func (s String) Get() string {
	return s.s.String
}

// TrimFormatStr 去掉字符串中的空格、换行符、回车符
func TrimFormatStr(str string) string {
	// 去除首尾空格
	str = strings.Trim(str, " ")
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	// 去掉回车符
	str = strings.Replace(str, "\r", "", -1)
	return str
}
