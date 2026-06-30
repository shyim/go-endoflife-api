package endoflife

import (
	"strings"
	"time"
)

// dateLayout is the format used by the endoflife.date API for date-only fields.
const dateLayout = "2006-01-02"

// Date represents a calendar date (without time) as returned by the
// endoflife.date API. It marshals to and from the "YYYY-MM-DD" format.
//
// API fields documented as date-or-null are represented as *Date, where a nil
// pointer denotes a JSON null (i.e. the information is not known).
type Date struct {
	time.Time
}

// NewDate returns a Date for the given year, month and day (UTC).
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON implements json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format(dateLayout) + `"`), nil
}

// String returns the date formatted as "YYYY-MM-DD", or an empty string if the
// date is the zero value.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format(dateLayout)
}
