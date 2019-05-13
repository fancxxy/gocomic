package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time struct custom time formatter
type Time struct {
	time.Time
}

// MarshalJSON function override
func (t Time) MarshalJSON() ([]byte, error) {
	var date string
	if t.IsZero() {
		date = `"-"`
	} else {
		date = fmt.Sprintf(`"%s"`, t.Format("2006-01-02"))
	}
	return []byte(date), nil
}

// Value function override
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan function override
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
