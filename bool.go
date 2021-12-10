// Package libs ============================================================
// 描述:
// 作者: Simon
// 日期: 2019-08-03 20:33
// 版权: 山东高招信息科技有限公司 @Since 2019
//
//============================================================
package libs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Bool is a nullable bool. False input is considered null.
// JSON marshals to false if null.
// Considered null to SQL unmarshalled from a false value.
type Bool struct {
	sql.NullBool
}

// NewBool creates a new Bool
func NewBool(b bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
// "false" will be considered a null Bool.
// It also supports unmarshalling a sql.NullBool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
		b.Valid = true
	case map[string]interface{}:
		err = json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Valid = false
		return nil
	case string:
		if len(x) == 0 {
			b.Valid = false
			return nil
		}
		b.Bool, err = strconv.ParseBool(x)
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is a false or not a bool.
// It will return an error if the input is not a float, blank, or "null".
func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Bool = false
		b.Valid = false
		return nil
	case "true":
		b.Valid = true
		b.Bool = true
	case "false":
		b.Bool = false
		b.Valid = true
	case "1":
		b.Valid = true
		b.Bool = true
	case "0":
		b.Bool = false
		b.Valid = true

	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	}

	return nil
}

// MarshalJSON implements json.Marshaller.
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaller.
// It will encode a zero if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this Bool value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a poBooler to this Bool value, or a nil poBooler if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsNil returns true for null or zero Bools, for future omitempty support (Go 1.4?)
func (b Bool) IsNil() bool {
	return !b.Valid
}

//Get returns base value for Bools
func (b Bool) Get() bool {
	return b.Bool
}

func (b Bool) String() string {
	if !b.Valid {
		return ""
	}
	return fmt.Sprint(b.Get())
}
