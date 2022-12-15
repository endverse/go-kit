package time

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	gotime "time"

	"github.com/endverse/go-kit/codec"
)

const layout = "2006-01-02 15:04:05"

var cst = gotime.FixedZone("CST", 8*3600)

func Location() *gotime.Location {
	return cst
}

func CST() *gotime.Location {
	return cst
}

func Layout() string {
	return layout
}

func ParseInLocation(t string) *Time {
	tt, _ := time.ParseInLocation(layout, t, cst)

	return &Time{tt}
}

type Time struct {
	gotime.Time
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (t *Time) UnmarshalJSON(b []byte) error {
	if (len(b) == 4 && string(b) == "null") || (len(b) == 2 && string(b) == "\"\"") {
		t.Time = gotime.Time{}
		return nil
	}

	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	pt, err := gotime.ParseInLocation(layout, str, cst)
	if err != nil {
		return err
	}

	t.Time = pt.Local()
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		// Encode unset/nil objects as JSON's "null".
		return []byte("null"), nil
	}

	buf := make([]byte, 0, len(layout)+2)
	buf = append(buf, '"')
	// time cannot contain non escapable JSON characters
	buf = t.In(cst).AppendFormat(buf, layout)
	buf = append(buf, '"')
	return buf, nil
}

func (t *Time) Scan(value interface{}) error {
	return codec.Scan(t, value)
}

func (t Time) Value() (driver.Value, error) {
	return codec.Value(&t)
}

func (t *Time) IsZero() bool {
	if t == nil {
		return true
	}
	return t.Time.IsZero()
}

func (t *Time) In(location *gotime.Location) *Time {
	return &Time{t.Time.In(location)}
}

func (t *Time) CST() *Time {
	return &Time{t.Time.In(CST())}
}

func (t *Time) Location() *gotime.Location {
	if t.Time.IsZero() {
		return cst
	}

	return t.Time.Location()
}

func (t *Time) Start() *Time {
	if t == nil || t.Time.IsZero() {
		return &Time{}
	}

	start := fmt.Sprintf("%04d-%02d-%02d 00:00:00", t.Year(), t.Month(), t.Day())
	startTime, _ := gotime.ParseInLocation(Layout(), start, t.Time.Location())
	return &Time{startTime}
}

func (t *Time) End() *Time {
	if t.Time.IsZero() || t.Time.IsZero() {
		return &Time{}
	}

	end := fmt.Sprintf("%04d-%02d-%02d 23:59:59", t.Year(), t.Month(), t.Day())
	endTime, _ := gotime.ParseInLocation(Layout(), end, t.Time.Location())
	return &Time{endTime}
}

func (t *Time) Date() string {
	if t == nil || t.Time.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func (t *Time) String() string {
	if t == nil || t.Time.IsZero() {
		return ""
	}

	return t.Format(Layout())
}
