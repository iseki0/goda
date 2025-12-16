package goda

import (
	"database/sql/driver"
	"errors"
	"time"
)

// UnmarshalJSON implements the json.Unmarshaler interface.
// It accepts JSON strings in yyyy-MM-dd format or JSON null.
func (d *LocalDate) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*d = LocalDate{}
		return nil
	}
	return unmarshalJsonImpl(d, bytes)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses dates in yyyy-MM-dd format.
// Empty input is treated as zero value.
func (d *LocalDate) UnmarshalText(text []byte) (e error) {
	defer deferOpInParse(text, &e)
	if len(text) == 0 {
		*d = LocalDate{}
		return nil
	}
	if text[len(text)-3] != '-' || text[len(text)-6] != '-' {
		return errors.New("'-' required")
	}
	var y int64
	var m, dom int
	dom, e = parseInt(text[len(text)-2:])
	if e != nil {
		return
	}
	m, e = parseInt(text[len(text)-5 : len(text)-3])
	if e != nil {
		return
	}
	y, e = parseInt64(text[:len(text)-6])
	if e != nil {
		return
	}
	dd, e := LocalDateOf(Year(y), Month(m), dom)
	if e != nil {
		return
	}
	*d = dd
	return
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the date in yyyy-MM-dd format, or empty for zero value.
func (d LocalDate) MarshalText() (text []byte, err error) {
	return marshalTextImpl(d)
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the date as a JSON string in yyyy-MM-dd format, or empty string for zero value.
func (d LocalDate) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(d)
}

// String returns the date in yyyy-MM-dd format, or empty string for zero value.
func (d LocalDate) String() string {
	return stringImpl(d)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the date in yyyy-MM-dd format to b and returns the extended buffer.
func (d LocalDate) AppendText(b []byte) ([]byte, error) {
	if d.IsZero() {
		return b, nil
	}
	b, _ = d.Year().AppendText(b)
	b = append(b, '-', byte('0'+d.Month()/10), byte('0'+d.Month()%10), '-', byte('0'+d.DayOfMonth()/10), byte('0'+d.DayOfMonth()%10))
	return b, nil
}

// Value implements the driver.Valuer interface.
// It returns nil for zero values, otherwise returns the date as a string in yyyy-MM-dd format.
func (d LocalDate) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.String(), nil
}

// Scan implements the sql.Scanner interface.
// It supports scanning from nil, string, []byte, and time.Time.
// Nil values are converted to the zero value of LocalDate.
func (d *LocalDate) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*d = LocalDate{}
		return nil
	case string:
		return d.UnmarshalText([]byte(v))
	case []byte:
		return d.UnmarshalText(v)
	case time.Time:
		*d = LocalDateOfGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}
