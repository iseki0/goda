package goda

import (
	"database/sql/driver"
	"errors"
	"time"
)

// String returns the ISO 8601 string representation (yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]).
func (dt LocalDateTime) String() string {
	return stringImpl(dt)
}

// AppendText implements encoding.TextAppender.
func (dt LocalDateTime) AppendText(b []byte) ([]byte, error) {
	if dt.IsZero() {
		return b, nil
	}
	b, _ = dt.date.AppendText(b)
	b = append(b, 'T')
	b, _ = dt.time.AppendText(b)
	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (dt LocalDateTime) MarshalText() ([]byte, error) {
	return marshalTextImpl(dt)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Accepts ISO 8601 format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]
func (dt *LocalDateTime) UnmarshalText(text []byte) (e error) {
	defer deferOpInParse(text, &e)
	if len(text) == 0 {
		*dt = LocalDateTime{}
		return nil
	}

	// Find the 'T' separator
	sepIdx := -1
	for i, ch := range text {
		if ch == 'T' || ch == 't' || ch == ' ' {
			sepIdx = i
			break
		}
	}

	if sepIdx < 0 {
		return errors.New("invalid date-time format: missing 'T' separator")
	}

	// Parse date part
	var date LocalDate
	if err := date.UnmarshalText(text[:sepIdx]); err != nil {
		return err
	}

	// Parse time part
	var timePart LocalTime
	if err := timePart.UnmarshalText(text[sepIdx+1:]); err != nil {
		return err
	}

	*dt = date.AtTime(timePart)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (dt LocalDateTime) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(dt)
}

// UnmarshalJSON implements json.Unmarshaler.
func (dt *LocalDateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		*dt = LocalDateTime{}
		return nil
	}
	return unmarshalJsonImpl(dt, data)
}

// Scan implements sql.Scanner.
func (dt *LocalDateTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dt = LocalDateTime{}
		return nil
	case string:
		return dt.UnmarshalText([]byte(v))
	case []byte:
		return dt.UnmarshalText(v)
	case time.Time:
		*dt = LocalDateTimeOfGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements driver.Valuer.
func (dt LocalDateTime) Value() (driver.Value, error) {
	if dt.IsZero() {
		return nil, nil
	}
	return dt.String(), nil
}
