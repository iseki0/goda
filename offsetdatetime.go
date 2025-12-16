package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// OffsetDateTime represents a date-time with a time-zone offset from UTC,
// such as 2024-03-15T14:30:45.123456789+01:00.
//
// OffsetDateTime stores an offset from UTC, but does not store or use a time-zone.
// Use ZonedDateTime when you need full time-zone support including DST transitions.
//
// OffsetDateTime is comparable and can be used as a map key.
// The zero value represents an unset date-time and IsZero returns true for it.
//
// OffsetDateTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm (e.g., "2024-03-15T14:30:45.123456789+01:00").
type OffsetDateTime struct {
	datetime LocalDateTime
	offset   ZoneOffset
}

// LocalDateTime returns the local date-time part without the offset.
func (odt OffsetDateTime) LocalDateTime() LocalDateTime {
	return odt.datetime
}

// LocalDate returns the date part of this date-time.
func (odt OffsetDateTime) LocalDate() LocalDate {
	return odt.datetime.LocalDate()
}

// LocalTime returns the time part of this date-time.
func (odt OffsetDateTime) LocalTime() LocalTime {
	return odt.datetime.LocalTime()
}

// Offset returns the zone offset.
func (odt OffsetDateTime) Offset() ZoneOffset {
	return odt.offset
}

// Year returns the year component.
func (odt OffsetDateTime) Year() Year {
	return odt.datetime.Year()
}

// Month returns the month component.
func (odt OffsetDateTime) Month() Month {
	return odt.datetime.Month()
}

// DayOfMonth returns the day-of-month component.
func (odt OffsetDateTime) DayOfMonth() int {
	return odt.datetime.DayOfMonth()
}

// DayOfWeek returns the day-of-week.
func (odt OffsetDateTime) DayOfWeek() DayOfWeek {
	return odt.datetime.DayOfWeek()
}

// DayOfYear returns the day-of-year.
func (odt OffsetDateTime) DayOfYear() int {
	return odt.datetime.DayOfYear()
}

// Hour returns the hour component (0-23).
func (odt OffsetDateTime) Hour() int {
	return odt.datetime.Hour()
}

// Minute returns the minute component (0-59).
func (odt OffsetDateTime) Minute() int {
	return odt.datetime.Minute()
}

// Second returns the second component (0-59).
func (odt OffsetDateTime) Second() int {
	return odt.datetime.Second()
}

// Millisecond returns the millisecond component (0-999).
func (odt OffsetDateTime) Millisecond() int {
	return odt.datetime.Millisecond()
}

// Nanosecond returns the nanosecond component (0-999999999).
func (odt OffsetDateTime) Nanosecond() int {
	return odt.datetime.Nanosecond()
}

// IsZero returns true if this is the zero value.
func (odt OffsetDateTime) IsZero() bool {
	return odt.datetime.IsZero() && odt.offset.IsZero()
}

// IsLeapYear returns true if the year is a leap year.
func (odt OffsetDateTime) IsLeapYear() bool {
	return odt.datetime.IsLeapYear()
}

// IsSupportedField returns true if the field is supported by OffsetDateTime.
// OffsetDateTime supports all fields from LocalDateTime plus FieldOffsetSeconds and FieldInstantSeconds.
func (odt OffsetDateTime) IsSupportedField(field Field) bool {
	return odt.datetime.IsSupportedField(field) || odt.offset.IsSupportedField(field) || field == FieldInstantSeconds
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the date-time for the value of the specified field.
// The returned value may be unsupported if the field is not supported by OffsetDateTime.
//
// If the date-time is zero (IsZero() returns true), an unsupported TemporalValue is returned.
//
// OffsetDateTime supports all fields from LocalDateTime plus:
//   - FieldOffsetSeconds: the offset in seconds
//   - FieldInstantSeconds: the epoch seconds (Unix timestamp)
func (odt OffsetDateTime) GetField(field Field) TemporalValue {
	if odt.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}
	if odt.offset.IsSupportedField(field) {
		return odt.offset.GetField(field)
	}

	// Handle offset-specific fields
	switch field {
	case FieldInstantSeconds:
		v, o := odt.epochSecondOverflow()
		return TemporalValue{v: v, overflow: o}
	case FieldOffsetSeconds:
		return TemporalValue{v: int64(odt.Offset().totalSeconds)}
	}

	// Delegate to LocalDateTime for date and time fields
	return odt.datetime.GetField(field)
}

// GoTime converts this offset date-time to a time.Time.
// Returns time.Time{} (zero) for zero value.
func (odt OffsetDateTime) GoTime() time.Time {
	if odt.IsZero() {
		return time.Time{}
	}
	// Create a fixed location with the offset
	loc := time.FixedZone("", odt.offset.TotalSeconds())
	return time.Date(int(odt.Year()), time.Month(odt.Month()), odt.DayOfMonth(), odt.Hour(), odt.Minute(), odt.Second(), odt.Nanosecond(), loc)
}

// EpochSecond returns the number of seconds since Unix epoch (1970-01-01T00:00:00Z).
func (odt OffsetDateTime) EpochSecond() int64 {
	i, _ := odt.epochSecondOverflow()
	return i
}

func (odt OffsetDateTime) epochSecondOverflow() (i int64, overflow bool) {
	if odt.IsZero() {
		return 0, false
	}
	epochDay := odt.datetime.LocalDate().UnixEpochDays()
	secondsOfDay := odt.datetime.LocalTime().GetField(FieldSecondOfDay).Int64()
	i, overflow = addExactly(epochDay*86400+secondsOfDay, -int64(odt.offset.TotalSeconds()))
	overflow = overflow || odt.LocalDateTime().Compare(localDateTimeMinEpochSecond) < 0 || odt.LocalDateTime().Compare(localDateTimeMaxEpochSecond) > 0
	return
}

// Compare compares this offset date-time with another.
// The comparison is based on the instant then on the local date-time.
// Returns -1 if this is before other, 0 if equal, 1 if after.
func (odt OffsetDateTime) Compare(other OffsetDateTime) int {
	return doCompare(odt, other, compareZero, comparing(OffsetDateTime.EpochSecond), comparing(OffsetDateTime.Nanosecond), comparing1(OffsetDateTime.LocalDateTime))
}

// IsBefore returns true if this offset date-time is before the specified offset date-time.
func (odt OffsetDateTime) IsBefore(other OffsetDateTime) bool {
	return doCompare(odt, other, compareZero, comparing(OffsetDateTime.EpochSecond), comparing(OffsetDateTime.Nanosecond)) < 0
}

// IsAfter returns true if this offset date-time is after the specified offset date-time.
func (odt OffsetDateTime) IsAfter(other OffsetDateTime) bool {
	return doCompare(odt, other, compareZero, comparing(OffsetDateTime.EpochSecond), comparing(OffsetDateTime.Nanosecond)) > 0
}

func (odt OffsetDateTime) Chain() (chain OffsetDateTimeChain) {
	chain.value = odt
	return
}

// OffsetDateTimeOf creates a new OffsetDateTime from individual components.
// Returns an error if any component is invalid.
func OffsetDateTimeOf(year Year, month Month, day int, hour int, minute int, second int, nanosecond int, offset ZoneOffset) (OffsetDateTime, error) {
	dt, err := LocalDateTimeOf(year, month, day, hour, minute, second, nanosecond)
	if err != nil {
		return OffsetDateTime{}, err
	}
	return OffsetDateTime{
		datetime: dt,
		offset:   offset,
	}, nil
}

// MustOffsetDateTimeOf creates a new OffsetDateTime from individual components.
// Panics if any component is invalid.
func MustOffsetDateTimeOf(year Year, month Month, day int, hour int, minute int, second int, nanosecond int, offset ZoneOffset) OffsetDateTime {
	return mustValue(OffsetDateTimeOf(year, month, day, hour, minute, second, nanosecond, offset))
}

// OffsetDateTimeNow returns the current date-time with offset in the system's local time zone.
func OffsetDateTimeNow() OffsetDateTime {
	return OffsetDateTimeOfGoTime(time.Now())
}

// OffsetDateTimeNowUTC returns the current date-time with offset in UTC.
func OffsetDateTimeNowUTC() OffsetDateTime {
	return OffsetDateTimeOfGoTime(time.Now().UTC())
}

// OffsetDateTimeOfGoTime creates an OffsetDateTime from a time.Time.
// Returns zero value if t.IsZero().
func OffsetDateTimeOfGoTime(t time.Time) OffsetDateTime {
	if t.IsZero() {
		return OffsetDateTime{}
	}
	_, offsetSeconds := t.Zone()
	offset := MustZoneOffsetOfSeconds(offsetSeconds)
	return LocalDateTimeOfGoTime(t).AtOffset(offset)
}

// OffsetDateTimeParse parses a date-time string in RFC3339-compatible format.
// The format is yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss] or with 'Z' for UTC.
//
// Examples:
//
//	odt, err := OffsetDateTimeParse("2024-03-15T14:30:45.123456789+01:00")
//	odt, err := OffsetDateTimeParse("2024-03-15T14:30:45Z")
//	odt, err := OffsetDateTimeParse("2024-03-15T14:30:45+05:30")
func OffsetDateTimeParse(s string) (OffsetDateTime, error) {
	var odt OffsetDateTime
	err := odt.UnmarshalText([]byte(s))
	return odt, err
}

// MustOffsetDateTimeParse parses a date-time string with offset.
// Panics if the string is invalid.
func MustOffsetDateTimeParse(s string) OffsetDateTime {
	odt, err := OffsetDateTimeParse(s)
	if err != nil {
		panic(err)
	}
	return odt
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*OffsetDateTime)(nil)
	_ fmt.Stringer             = (*OffsetDateTime)(nil)
	_ encoding.TextMarshaler   = (*OffsetDateTime)(nil)
	_ encoding.TextUnmarshaler = (*OffsetDateTime)(nil)
	_ json.Marshaler           = (*OffsetDateTime)(nil)
	_ json.Unmarshaler         = (*OffsetDateTime)(nil)
	_ driver.Valuer            = (*OffsetDateTime)(nil)
	_ sql.Scanner              = (*OffsetDateTime)(nil)
)

// Compile-time check that OffsetDateTime is comparable
func _assertOffsetDateTimeIsComparable[T comparable](t T) {}

var _ = _assertOffsetDateTimeIsComparable[OffsetDateTime]
