package time

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	gotime "time"

	"go-arsenal.kanzhun.tech/arsenal/go-kit/codec"
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

type Time struct {
	gotime.Time
}

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

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	buf := make([]byte, 0, len(layout)+2)
	buf = append(buf, '"')
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

func Today() Time {
	return Time{gotime.Now()}
}

func Yesterday() Time {
	return Time{gotime.Now().AddDate(0, 0, -1)}
}

func Tomorrow() Time {
	return Time{gotime.Now().AddDate(0, 0, 1)}
}

func AnchorTime(years int, months int, days int) Time {
	return Time{gotime.Now().In(Location()).AddDate(years, months, days)}
}

// ********************* The days of the week. ********************* //
func Monday() Time {
	offset := int(gotime.Monday - gotime.Now().Weekday())
	if offset > 0 {
		offset = -6
	}
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Tuesday() Time {
	offset := int(gotime.Tuesday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Wednesday() Time {
	offset := int(gotime.Wednesday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Thursday() Time {
	offset := int(gotime.Thursday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Friday() Time {
	offset := int(gotime.Friday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Saturday() Time {
	offset := int(gotime.Saturday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Sunday() Time {
	offset := int(gotime.Sunday - gotime.Now().Weekday())
	if offset < 0 {
		offset = 7 - int(gotime.Now().Weekday())
	}
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func LastMonday() Time {
	return Time{Monday().AddDate(0, 0, -7)}
}

func LastSunday() Time {
	offset := int(gotime.Sunday - gotime.Now().Weekday())
	if offset == 0 {
		offset = 7 - int(gotime.Now().Weekday())
	}
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func NowInLocation(location *time.Location) Time {
	return Time{gotime.Now().In(location)}
}

func Now() Time {
	return Time{gotime.Now().In(CST())}
}

func (t Time) In(location *gotime.Location) Time {
	return Time{t.Time.In(location)}
}

func (t Time) CST() Time {
	return Time{t.Time.In(CST())}
}

func (t Time) Location() *gotime.Location {
	return t.Time.Location()
}

func (t Time) Start() Time {
	start := fmt.Sprintf("%04d-%02d-%02d 00:00:00", t.Year(), t.Month(), t.Day())
	startTime, _ := gotime.ParseInLocation(Layout(), start, t.Time.Location())
	return Time{startTime}
}

func (t Time) End() Time {
	end := fmt.Sprintf("%04d-%02d-%02d 23:59:59", t.Year(), t.Month(), t.Day())
	endTime, _ := gotime.ParseInLocation(Layout(), end, t.Time.Location())
	return Time{endTime}
}

func (t Time) Date() string {
	return t.Format("2006-01-02")
}

func (t Time) String() string {
	return t.Format(Layout())
}
