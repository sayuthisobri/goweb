package web

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Date time.Time

func (d *Date) Scan(v interface{}) error {
	if v != nil {
		fmt.Print(v)
	}
	return nil
}

func (d *Date) Value() (driver.Value, error) {
	return time.Time(*d), nil
}

//
// UnmarshalJSON -
//
func (d *Date) UnmarshalJSON(bytes []byte) error {
	layout := "2006-01-02"
	str := strings.Trim(string(bytes), "\"")
	if t, err := time.Parse(layout, str); err != nil {
		return err
	} else {
		*d = Date(t)
		return nil
	}
}
