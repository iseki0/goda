package goda

import (
	"database/sql/driver"
	"errors"
	"time"
)

// String returns the ISO 8601 string representation (yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm).
func (odt OffsetDateTime) String() string {
	return stringImpl(odt)
}

// AppendText implements encoding.TextAppender.
func (odt OffsetDateTime) AppendText(b []byte) ([]byte, error) {
	if odt.IsZero() {
		return b, nil
	}
	b, _ = odt.datetime.AppendText(b)
	b, _ = odt.offset.AppendText(b)
	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (odt OffsetDateTime) MarshalText() ([]byte, error) {
	return marshalTextImpl(odt)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Accepts ISO 8601 format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss] or Z for UTC.
func (odt *OffsetDateTime) UnmarshalText(text []byte) (e error) {
	defer deferOpInParse(text, &e)
	if len(text) == 0 {
		*odt = OffsetDateTime{}
		return nil
	}

	// Find the offset part (starts with +, -, or Z)
	offsetIdx := -1
	for i := len(text) - 1; i >= 0; i-- {
		ch := text[i]
		if ch == '+' || ch == '-' {
			offsetIdx = i
			break
		}
		if ch == 'Z' || ch == 'z' {
			offsetIdx = i
			break
		}
	}

	if offsetIdx < 0 {
		return errors.New("invalid offset date-time format: missing offset")
	}

	// Parse date-time part
	var dt LocalDateTime
	if err := dt.UnmarshalText(text[:offsetIdx]); err != nil {
		return err
	}

	// Parse offset part
	var offset ZoneOffset
	if err := offset.UnmarshalText(text[offsetIdx:]); err != nil {
		return err
	}

	*odt = OffsetDateTime{
		datetime: dt,
		offset:   offset,
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (odt OffsetDateTime) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(odt)
}

// UnmarshalJSON implements json.Unmarshaler.
func (odt *OffsetDateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		*odt = OffsetDateTime{}
		return nil
	}
	return unmarshalJsonImpl(odt, data)
}

// Scan implements sql.Scanner.
func (odt *OffsetDateTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*odt = OffsetDateTime{}
		return nil
	case string:
		return odt.UnmarshalText([]byte(v))
	case []byte:
		return odt.UnmarshalText(v)
	case time.Time:
		*odt = OffsetDateTimeOfGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements driver.Valuer.
func (odt OffsetDateTime) Value() (driver.Value, error) {
	if odt.IsZero() {
		return nil, nil
	}
	return odt.String(), nil
}
