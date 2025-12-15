package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// LocalTime represents a time without a date or time zone,
// such as 14:30:45.123456789. It stores the hour, minute, second, and nanosecond.
//
// LocalTime is comparable and can be used as a map key.
// The zero value represents an unset time and IsZero returns true for it.
// Note: 00:00:00 (midnight) is a valid time and is different from the zero value.
//
// LocalTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: HH:mm:ss[.nnnnnnnnn] (e.g., "14:30:45.123456789"). Uses 24-hour format.
// Fractional seconds support nanosecond precision and are aligned to 3-digit boundaries
// (milliseconds, microseconds, nanoseconds) for Java.time compatibility.
// Timezone offsets are not supported.
//
// Internally, v uses bit 63 (the highest bit) as a validity flag.
// If bit 63 is set, the time is valid and bits 0-62 contain nanoseconds since midnight.
// If bit 63 is clear (v == 0), the time is invalid/zero.
// Max nanoseconds in a day: 86,399,999,999,999 (< 2^47), so we have plenty of space.
type LocalTime struct {
	v int64 // bit 63: valid flag; bits 0-62: nanoseconds since midnight
}

const (
	localTimeValidBit  int64 = -0x8000000000000000 // bit 63 set (sign bit in two's complement)
	localTimeValueMask int64 = 0x7FFFFFFFFFFFFFFF  // bits 0-62
)

// Hour returns the hour component (0-23).
func (t LocalTime) Hour() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Hour))
}

// Minute returns the minute component (0-59).
func (t LocalTime) Minute() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Minute) % 60)
}

// Second returns the second component (0-59).
func (t LocalTime) Second() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Second) % 60)
}

// Millisecond returns the millisecond component (0-999).
func (t LocalTime) Millisecond() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Millisecond) % 1000)
}

// Nano returns the nanosecond component (0-999999999).
func (t LocalTime) Nano() int {
	ns := t.v & localTimeValueMask
	return int(ns % int64(time.Second))
}

// IsSupportedField returns true if the field is supported by LocalTime.
func (t LocalTime) IsSupportedField(field Field) bool {
	switch field {
	case FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldHourOfDay, FieldClockHourOfDay, FieldAmPmOfDay:
		return true
	default:
		return false
	}
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the time for the value of the specified field.
// The returned value may be unsupported if the field is not supported by LocalTime.
//
// If the time is zero (IsZero() returns true), an unsupported TemporalValue is returned.
// For fields not supported by LocalTime (such as date fields), an unsupported TemporalValue is returned.
//
// Supported fields include:
//   - FieldNanoOfSecond: returns the nanosecond within the second (0-999,999,999)
//   - FieldNanoOfDay: returns the nanosecond of day (0-86,399,999,999,999)
//   - FieldMicroOfSecond: returns the microsecond within the second (0-999,999)
//   - FieldMicroOfDay: returns the microsecond of day (0-86,399,999,999)
//   - FieldMilliOfSecond: returns the millisecond within the second (0-999)
//   - FieldMilliOfDay: returns the millisecond of day (0-86,399,999)
//   - FieldSecondOfMinute: returns the second within the minute (0-59)
//   - FieldSecondOfDay: returns the second of day (0-86,399)
//   - FieldMinuteOfHour: returns the minute within the hour (0-59)
//   - FieldMinuteOfDay: returns the minute of day (0-1,439)
//   - FieldHourOfDay: returns the hour of day (0-23)
//   - FieldClockHourOfDay: returns the clock hour of day (1-24)
//   - FieldHourOfAmPm: returns the hour within AM/PM (0-11)
//   - FieldClockHourOfAmPm: returns the clock hour within AM/PM (1-12)
//   - FieldAmPmOfDay: returns AM/PM indicator (0=AM, 1=PM)
//
// Overflow Analysis:
// None of the supported fields can overflow int64:
//   - FieldNanoOfSecond: range 0-999,999,999, cannot overflow
//   - FieldNanoOfDay: max 86,399,999,999,999 (< 2^47), cannot overflow
//   - FieldMicroOfSecond: range 0-999,999, cannot overflow
//   - FieldMicroOfDay: max 86,399,999,999 (< 2^37), cannot overflow
//   - FieldMilliOfSecond: range 0-999, cannot overflow
//   - FieldMilliOfDay: max 86,399,999 (< 2^27), cannot overflow
//   - FieldSecondOfMinute: range 0-59, cannot overflow
//   - FieldSecondOfDay: max 86,399 (< 2^17), cannot overflow
//   - FieldMinuteOfHour: range 0-59, cannot overflow
//   - FieldMinuteOfDay: max 1,439 (< 2^11), cannot overflow
//   - FieldHourOfDay: range 0-23, cannot overflow
//   - FieldClockHourOfDay: range 1-24, cannot overflow
//   - FieldHourOfAmPm: range 0-11, cannot overflow
//   - FieldClockHourOfAmPm: range 1-12, cannot overflow
//   - FieldAmPmOfDay: values 0 or 1, cannot overflow
func (t LocalTime) GetField(field Field) TemporalValue {
	if t.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}
	var v int64
	ns := t.v & localTimeValueMask
	switch field {
	case FieldNanoOfSecond:
		// Range: 0-999,999,999, no overflow possible
		v = ns % int64(time.Second)
	case FieldNanoOfDay:
		// Max: 86,399,999,999,999 (< 2^47), no overflow possible
		v = ns
	case FieldMicroOfSecond:
		// Range: 0-999,999, no overflow possible
		v = ns % int64(time.Second) / int64(time.Microsecond)
	case FieldMicroOfDay:
		// Max: 86,399,999,999 (< 2^37), no overflow possible
		v = ns / int64(time.Microsecond)
	case FieldMilliOfSecond:
		// Range: 0-999, no overflow possible
		v = ns % int64(time.Second) / int64(time.Millisecond)
	case FieldMilliOfDay:
		// Max: 86,399,999 (< 2^27), no overflow possible
		v = ns / int64(time.Millisecond)
	case FieldSecondOfMinute:
		// Range: 0-59, no overflow possible
		v = ns / int64(time.Second) % 60
	case FieldSecondOfDay:
		// Max: 86,399 (< 2^17), no overflow possible
		v = ns / int64(time.Second)
	case FieldMinuteOfHour:
		// Range: 0-59, no overflow possible
		v = ns / int64(time.Minute) % 60
	case FieldMinuteOfDay:
		// Max: 1,439 (< 2^11), no overflow possible
		v = ns / int64(time.Minute)
	case FieldHourOfDay:
		// Range: 0-23, no overflow possible
		v = ns / int64(time.Hour)
	case FieldClockHourOfDay:
		// Range: 1-24, no overflow possible
		h := ns / int64(time.Hour)
		if h == 0 {
			v = 24
		} else {
			v = h
		}
	case FieldHourOfAmPm:
		// Range: 0-11, no overflow possible
		v = ns / int64(time.Hour) % 12
	case FieldClockHourOfAmPm:
		// Range: 1-12, no overflow possible
		h := ns / int64(time.Hour) % 12
		if h == 0 {
			v = 12
		} else {
			v = h
		}
	case FieldAmPmOfDay:
		// Values: 0 or 1, no overflow possible
		if ns/int64(time.Hour) < 12 {
			v = 0 // AM
		} else {
			v = 1 // PM
		}
	default:
		return TemporalValue{v: 0, unsupported: true}
	}
	return TemporalValue{v: v}
}

// IsZero returns true if this is the zero value of LocalTime.
func (t LocalTime) IsZero() bool {
	return t.v == 0
}

// GoTime converts this time to a time.Time at the Unix epoch date (1970-01-01) in UTC.
// Returns time.Time{} (zero) for zero value.
func (t LocalTime) GoTime() time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	nanos := t.v & localTimeValueMask
	return time.Unix(0, nanos).UTC()
}

// Compare compares this time with another time.
// Returns -1 if this time is before other, 0 if equal, and 1 if after.
// Zero values are considered less than non-zero values.
func (t LocalTime) Compare(other LocalTime) int {
	return doCompare(t, other, compareZero, comparing(LocalTime.NanoOfDay))
}

// IsBefore returns true if this time is before the specified time.
func (t LocalTime) IsBefore(other LocalTime) bool {
	return t.Compare(other) < 0
}

// IsAfter returns true if this time is after the specified time.
func (t LocalTime) IsAfter(other LocalTime) bool {
	return t.Compare(other) > 0
}

// AtDate combines this time with a date to create a LocalDateTime.
func (t LocalTime) AtDate(date LocalDate) LocalDateTime {
	return LocalDateTime{
		date: date,
		time: t,
	}
}

func (t LocalTime) SecondOfDay() int {
	return int(t.NanoOfDay() / 1000_000_000)
}

func (t LocalTime) NanoOfDay() int64 {
	return t.v & localTimeValueMask
}

func (t LocalTime) Chain() (chain LocalTimeChain) {
	chain.value = t
	return
}

func (t LocalTime) chainWithError(e error) (chain LocalTimeChain) {
	chain = t.Chain()
	chain.eError = e
	return
}

// LocalTimeOf creates a new LocalTime from the specified hour, minute, second, and nanosecond.
// Returns an error if any component is out of range:
// - hour must be 0-23
// - minute must be 0-59
// - second must be 0-59
// - nanosecond must be 0-999999999
func LocalTimeOf(hour, minute, second, nanosecond int) (r LocalTime, e error) {
	FieldHourOfDay.checkSetE(int64(hour), &e)
	FieldMinuteOfHour.checkSetE(int64(minute), &e)
	FieldSecondOfMinute.checkSetE(int64(second), &e)
	FieldNanoOfSecond.checkSetE(int64(nanosecond), &e)
	if e != nil {
		return
	}
	nanos := int64(hour)*int64(time.Hour) +
		int64(minute)*int64(time.Minute) +
		int64(second)*int64(time.Second) +
		int64(nanosecond)
	return LocalTime{
		v: nanos | localTimeValidBit,
	}, nil
}

// MustLocalTimeOf creates a new LocalTime from the specified hour, minute, second, and nanosecond.
// Panics if any component is out of range. Use LocalTimeOf for error handling.
func MustLocalTimeOf(hour, minute, second, nanosecond int) LocalTime {
	return mustValue(LocalTimeOf(hour, minute, second, nanosecond))
}

// LocalTimeOfGoTime creates a LocalTime from a time.Time.
// The date and time zone components are ignored.
// Returns zero value if t.IsZero().
func LocalTimeOfGoTime(t time.Time) LocalTime {
	if t.IsZero() {
		return LocalTime{}
	}
	nanos := time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC).UnixNano()
	return LocalTime{
		v: nanos | localTimeValidBit,
	}
}

// LocalTimeOfNanoOfDay creates a LocalTime from the nanosecond-of-day value.
// Returns an error if nanoOfDay is out of the valid range (0 to 86,399,999,999,999).
// The valid range represents 00:00:00.000000000 to 23:59:59.999999999.
//
// Example:
//
//	// Create time for 12:00:00
//	lt, err := LocalTimeOfNanoOfDay(12 * 60 * 60 * 1000000000)
func LocalTimeOfNanoOfDay(nanoOfDay int64) (r LocalTime, e error) {
	FieldNanoOfDay.checkSetE(nanoOfDay, &e)
	if e != nil {
		return
	}
	return LocalTime{
		v: nanoOfDay | localTimeValidBit,
	}, nil
}

// MustLocalTimeOfNanoOfDay creates a LocalTime from the nanosecond-of-day value.
// Panics if nanoOfDay is out of range. Use LocalTimeOfNanoOfDay for error handling.
//
// Example:
//
//	// Create time for 12:00:00
//	lt := MustLocalTimeOfNanoOfDay(12 * 60 * 60 * 1000000000)
func MustLocalTimeOfNanoOfDay(nanoOfDay int64) LocalTime {
	return mustValue(LocalTimeOfNanoOfDay(nanoOfDay))
}

// LocalTimeOfSecondOfDay creates a LocalTime from the second-of-day value.
// Returns an error if secondOfDay is out of the valid range (0 to 86,399).
// The valid range represents 00:00:00 to 23:59:59.
// The nanosecond component will be set to zero.
//
// Example:
//
//	// Create time for 12:00:00
//	lt, err := LocalTimeOfSecondOfDay(12 * 60 * 60)
func LocalTimeOfSecondOfDay(secondOfDay int) (r LocalTime, e error) {
	FieldSecondOfDay.checkSetE(int64(secondOfDay), &e)
	if e != nil {
		return
	}
	nanos := int64(secondOfDay) * int64(time.Second)
	return LocalTime{
		v: nanos | localTimeValidBit,
	}, nil
}

// MustLocalTimeOfSecondOfDay creates a LocalTime from the second-of-day value.
// Panics if secondOfDay is out of range. Use LocalTimeOfSecondOfDay for error handling.
//
// Example:
//
//	// Create time for 12:00:00
//	lt := MustLocalTimeOfSecondOfDay(12 * 60 * 60)
func MustLocalTimeOfSecondOfDay(secondOfDay int) LocalTime {
	return mustValue(LocalTimeOfSecondOfDay(secondOfDay))
}

// LocalTimeNow returns the current time in the system's local time zone.
// This is equivalent to LocalTimeOfGoTime(time.Now()).
// For UTC time, use LocalTimeNowUTC. For a specific timezone, use LocalTimeNowIn.
func LocalTimeNow() LocalTime {
	return LocalTimeOfGoTime(time.Now())
}

// LocalTimeNowIn returns the current time in the specified time zone.
// This is equivalent to LocalTimeOfGoTime(time.Now().In(loc)).
func LocalTimeNowIn(loc *time.Location) LocalTime {
	return LocalTimeOfGoTime(time.Now().In(loc))
}

// LocalTimeNowUTC returns the current time in UTC.
// This is equivalent to LocalTimeOfGoTime(time.Now().UTC()).
func LocalTimeNowUTC() LocalTime {
	return LocalTimeOfGoTime(time.Now().UTC())
}

// LocalTimeParse parses a time string in HH:mm:ss[.nnnnnnnnn] format (24-hour).
// Returns an error if the string is invalid or represents an invalid time.
//
// Accepts fractional seconds of any length (1-9 digits):
//   - HH:mm:ss (e.g., "14:30:45")
//   - HH:mm:ss.f (e.g., "14:30:45.1" → 100 milliseconds)
//   - HH:mm:ss.fff (e.g., "14:30:45.123" → 123 milliseconds)
//   - HH:mm:ss.ffffff (e.g., "14:30:45.123456" → 123.456 milliseconds)
//   - HH:mm:ss.nnnnnnnnn (e.g., "14:30:45.123456789" → full nanosecond precision)
//
// Note: Output formatting aligns to 3-digit boundaries (milliseconds, microseconds, nanoseconds).
// Timezone offsets are not supported.
//
// Example:
//
//	time, err := LocalTimeParse("14:30:45.123")
//	if err != nil {
//	    // handle error
//	}
func LocalTimeParse(s string) (LocalTime, error) {
	var t LocalTime
	err := t.UnmarshalText([]byte(s))
	return t, err
}

// MustLocalTimeParse parses a time string in HH:mm:ss[.nnnnnnnnn] format (24-hour).
// Panics if the string is invalid. Use LocalTimeParse for error handling.
//
// Example:
//
//	time := MustLocalTimeParse("14:30:45.123456789")
func MustLocalTimeParse(s string) LocalTime {
	return mustValue(LocalTimeParse(s))
}

var (
	_ encoding.TextAppender    = (*LocalTime)(nil)
	_ fmt.Stringer             = (*LocalTime)(nil)
	_ encoding.TextMarshaler   = (*LocalTime)(nil)
	_ encoding.TextUnmarshaler = (*LocalTime)(nil)
	_ json.Marshaler           = (*LocalTime)(nil)
	_ json.Unmarshaler         = (*LocalTime)(nil)
	_ driver.Valuer            = (*LocalTime)(nil)
	_ sql.Scanner              = (*LocalTime)(nil)
)

// Compile-time check that LocalTime is comparable
func _assertLocalTimeIsComparable[T comparable](T) {}

var _ = _assertLocalTimeIsComparable[LocalTime]
