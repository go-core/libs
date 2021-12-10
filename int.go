// Package libs ============================================================
// 描述:
// 作者: Simon
// 日期: 2019-08-03 20:38
// 版权: 山东高招信息科技有限公司 @Since 2019
//
//============================================================
package libs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Int is a nullable int64.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int struct {
	sql.NullInt64
}

// NewInt creates a new Int
func NewInt(i int64) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: true,
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case string:
		if len(x) == 0 {
			i.Valid = false
			return nil
		}

		i.Int64, err = strconv.ParseInt(x, 10, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = (err == nil) && (i.Int64 != 0)
	return err
}

// MarshalJSON implements json.Marshaller.
// It will encode 0 if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// MarshalText implements encoding.TextMarshaller.
// It will encode a zero if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// SetValid changes this Int value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsNil returns true for null or zero Int, for future omitempty support (Go 1.4?)
func (i Int) IsNil() bool {
	return !i.Valid
}

// Get 获取值
func (i *Int) Get() int {
	return int(i.Int64)
}

func (i *Int) String() string {
	return strconv.Itoa(i.Get())
}
