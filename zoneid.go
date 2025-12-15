package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var locLoadMap sync.Map

func loadLocation(id string) (l *time.Location, e error) {
	v, ok := locLoadMap.Load(id)
	if ok {
		return v.(*time.Location), nil
	}
	l, e = time.LoadLocation(id)
	if e != nil {
		return nil, e
	}
	locLoadMap.Store(id, l)
	return
}

// ZoneId represents a time zone identifier such as "America/New_York" or "Asia/Tokyo".
// It wraps a time.Location and provides serialization support for JSON, text, and SQL.
//
// ZoneId is comparable and can be used as a map key.
// The zero value represents an unset zone and IsZero returns true for it.
//
// ZoneId implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// encoding.TextAppender for efficient text appending,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: Uses IANA Time Zone Database names (e.g., "America/New_York", "Europe/London", "UTC").
// The string representation matches Go's time.Location.String() format.
type ZoneId struct {
	loc *time.Location
}

// ZoneIdOf creates a ZoneId from a time zone identifier string.
// The identifier must be a valid IANA Time Zone Database name (e.g., "America/New_York", "Asia/Tokyo").
// Returns an error if the zone ID is not found in the system's time zone database.
//
// Example zone IDs:
//   - "America/New_York" - Eastern Time
//   - "Europe/London" - British Time
//   - "Asia/Tokyo" - Japan Standard Time
//   - "UTC" - Coordinated Universal Time
func ZoneIdOf(id string) (ZoneId, error) {
	r, e := loadLocation(id)
	if e != nil {
		return ZoneId{}, e
	}
	return ZoneId{loc: r}, nil
}

// MustZoneIdOf creates a ZoneId from a time zone identifier string.
// Panics if the zone ID is invalid. Use ZoneIdOf for error handling.
func MustZoneIdOf(id string) ZoneId {
	return mustValue(ZoneIdOf(id))
}

// ZoneIdUTC returns a ZoneId representing UTC (Coordinated Universal Time).
func ZoneIdUTC() ZoneId {
	return ZoneId{time.UTC}
}

// ZoneIdOfGoLocation creates a ZoneId from a Go time.Location.
// This allows interoperability with Go's standard time package.
func ZoneIdOfGoLocation(l *time.Location) ZoneId {
	return ZoneId{l}
}

// ZoneIdDefault returns the system's default time zone.
// This corresponds to time.Local in Go's standard library.
// If time.Local is nil, returns ZoneIdUTC.
func ZoneIdDefault() ZoneId {
	if time.Local == nil {
		return ZoneIdUTC()
	}
	return ZoneId{time.Local}
}

// IsZero returns true if this is the zero value of ZoneId (no location set).
func (z ZoneId) IsZero() bool {
	return z.loc == nil
}

// String returns the zone ID as a string, or empty string for zero value.
func (z ZoneId) String() string {
	if z.IsZero() {
		return ""
	}
	return z.loc.String()
}

// Scan implements the sql.Scanner interface.
// It supports scanning from nil, string, and []byte.
// Nil values are converted to the zero value of ZoneId.
func (z *ZoneId) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*z = ZoneId{}
		return nil
	case string:
		return z.UnmarshalText([]byte(v))
	case []byte:
		return z.UnmarshalText(v)
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements the driver.Valuer interface.
// It returns nil for zero values, otherwise returns the zone ID as a string.
func (z ZoneId) Value() (driver.Value, error) {
	if z.IsZero() {
		return nil, nil
	}
	return z.String(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It accepts JSON strings representing a zone ID or JSON null.
func (z *ZoneId) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*z = ZoneId{}
		return nil
	}
	return unmarshalJsonImpl(z, bytes)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses zone IDs. Empty input is treated as zero value.
func (z *ZoneId) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*z = ZoneId{}
		return nil
	}
	zoneId, err := ZoneIdOf(string(text))
	if err != nil {
		return err
	}
	*z = zoneId
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the zone ID as text, or empty for zero value.
func (z ZoneId) MarshalText() (text []byte, err error) {
	return marshalTextImpl(z)
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the zone ID as a JSON string, or empty string for zero value.
func (z ZoneId) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(z)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the zone ID to b and returns the extended buffer.
func (z ZoneId) AppendText(b []byte) ([]byte, error) {
	if z.IsZero() {
		return b, nil
	}
	return append(b, z.String()...), nil
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*ZoneId)(nil)
	_ fmt.Stringer             = (*ZoneId)(nil)
	_ encoding.TextMarshaler   = (*ZoneId)(nil)
	_ encoding.TextUnmarshaler = (*ZoneId)(nil)
	_ json.Marshaler           = (*ZoneId)(nil)
	_ json.Unmarshaler         = (*ZoneId)(nil)
	_ driver.Valuer            = (*ZoneId)(nil)
	_ sql.Scanner              = (*ZoneId)(nil)
)

// Compile-time check that ZoneId is comparable
func _assertZoneIdIsComparable[T comparable](t T) {}

var _ = _assertZoneIdIsComparable[ZoneId]
