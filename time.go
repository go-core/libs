// Package libs ============================================================
// 描述:
// 作者: Simon
// 日期: 2019-08-03 20:39
// 版权: 山东高招信息科技有限公司 @Since 2019
//
//============================================================
package libs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Time is a nullable time.Time.
// JSON marshals to the zero value for time.Time if null.
// Considered to be null to SQL if zero.

var TimeLayOut = "2006-01-02 15:04:05"

type Time struct {
	Time  time.Time
	Valid bool
}

// Scan implements Scanner interface.
func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTime creates a new Time.
func NewTime(t time.Time) Time {
	return Time{
		Time:  t,
		Valid: true,
	}
}

// MarshalJSON implements json.Marshaller.
// It will encode the zero value of time.Time
// if this time is invalid.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.Format(TimeLayOut))
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *Time) UnmarshalJSON(data []byte) error {

	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch value.(type) {
	case string:
		return t.UnmarshalText([]byte(value.(string)))
	case nil:
		t.Time = time.Time{}
		t.Valid = false
		return nil
	default:
		return errors.New("不支持的序列化类型")
	}

}

//MarshalText json 序列化
func (t Time) MarshalText() ([]byte, error) {
	ti := t.Time
	if !t.Valid {
		ti = time.Time{}
	}
	return ti.MarshalText()
}

//UnmarshalText json 反序列化
func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	tt, err := time.ParseInLocation(TimeLayOut, str, time.Local)

	if err != nil {
		return err
	}
	t.Time = tt
	t.Valid = true
	return nil
}

// SetValid changes this Time's value and
// sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value,
// or a nil pointer if this Time is zero.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsNil returns true for null or zero Times, for potential future omitempty support.
func (t Time) IsNil() bool {
	return !t.Valid || t.Time.IsZero()
}

// Get 获取原始日期时间
func (t Time) Get() time.Time {
	return t.Time
}
