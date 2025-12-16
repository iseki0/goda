package goda

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Value implements driver.Valuer for database serialization.
func (t LocalTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.String(), nil
}

// Scan implements sql.Scanner for database deserialization.
func (t *LocalTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = LocalTime{}
		return nil
	case []byte:
		return t.UnmarshalText(v)
	case string:
		return t.UnmarshalText([]byte(v))
	case time.Time:
		*t = LocalTimeOfGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(src)
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *LocalTime) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*t = LocalTime{}
		return nil
	}
	return unmarshalJsonImpl(t, bytes)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *LocalTime) UnmarshalText(text []byte) (e error) {
	defer deferOpInParse(text, &e)
	if len(text) == 0 {
		*t = LocalTime{}
		return nil
	}

	if len(text) < 8 {
		return errors.New("invalid format")
	}

	hour, err := parseInt(text[0:2])
	if err != nil {
		return err
	}
	if text[2] != ':' {
		return errors.New("expect ':'")
	}
	minute, err := parseInt(text[3:5])
	if err != nil {
		return err
	}
	if text[5] != ':' {
		return errors.New("expect ':'")
	}
	second, err := parseInt(text[6:8])
	if err != nil {
		return err
	}

	var nano int64 = 0
	if len(text) > 8 {
		if text[8] != '.' {
			return errors.New("expect '.'")
		}
	}
	if len(text) > 9 {
		var nanoBuf [9]byte
		copy(nanoBuf[:], text[9:min(len(text), 18)])
		for i := 8; i >= 0 && nanoBuf[i] == 0; i-- {
			nanoBuf[i] = '0'
		}
		nano, e = parseInt64(nanoBuf[:])
		if e != nil {
			return e
		}
	}
	*t, e = LocalTimeOf(hour, minute, second, int(nano))
	return e
}

// AppendText in ISO-8601 format
func (t LocalTime) AppendText(b []byte) ([]byte, error) {
	if t.IsZero() {
		return b, nil
	}

	var buf [18]byte
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	nano := t.Nano()

	// Write hours, minutes, seconds
	buf[0] = byte('0' + hour/10)
	buf[1] = byte('0' + hour%10)
	buf[2] = ':'
	buf[3] = byte('0' + minute/10)
	buf[4] = byte('0' + minute%10)
	buf[5] = ':'
	buf[6] = byte('0' + second/10)
	buf[7] = byte('0' + second%10)

	if nano == 0 {
		return append(b, buf[:8]...), nil
	}

	buf[8] = '.'
	// Format nanoseconds
	ns := nano
	for i := 17; i >= 9; i-- {
		buf[i] = byte('0' + ns%10)
		ns /= 10
	}

	// Trim trailing zeros in groups of 3 (align to milliseconds, microseconds, nanoseconds)
	// This follows Java LocalTime behavior: 100000000ns -> ".100" not ".1"
	l := 18
	for l > 9 && buf[l-1] == '0' {
		l--
	}
	// Align to 3-digit boundaries (9, 12, 15, 18 total length)
	remainder := (l - 9) % 3
	if remainder != 0 {
		l += 3 - remainder
	}
	return append(b, buf[:l]...), nil
}

// MarshalText implements encoding.TextMarshaler.
func (t LocalTime) MarshalText() (text []byte, err error) {
	return marshalTextImpl(t)
}

// MarshalJSON implements json.Marshaler.
func (t LocalTime) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(t)
}

func (t LocalTime) String() string {
	return stringImpl(t)
}
