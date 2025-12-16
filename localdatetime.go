package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// LocalDateTime represents a date-time without a time zone,
// such as 2024-03-15T14:30:45.123456789. It combines LocalDate and LocalTime.
//
// LocalDateTime is comparable and can be used as a map key.
// The zero value represents an unset date-time and IsZero returns true for it.
//
// LocalDateTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn] (e.g., "2024-03-15T14:30:45.123456789").
// Combined date and time with 'T' separator. Lowercase 't' is accepted when parsing.
// Timezone offsets are not supported.
type LocalDateTime struct {
	date LocalDate
	time LocalTime
}

// LocalDate returns the date part of this date-time.
func (dt LocalDateTime) LocalDate() LocalDate {
	return dt.date
}

func (dt LocalDateTime) AtOffset(offset ZoneOffset) OffsetDateTime {
	if dt.IsZero() {
		return OffsetDateTime{}
	}
	return OffsetDateTime{
		datetime: dt,
		offset:   offset,
	}
}

// LocalTime returns the time part of this date-time.
func (dt LocalDateTime) LocalTime() LocalTime {
	return dt.time
}

// Year returns the year component.
func (dt LocalDateTime) Year() Year {
	return dt.date.Year()
}

// Month returns the month component.
func (dt LocalDateTime) Month() Month {
	return dt.date.Month()
}

// DayOfMonth returns the day-of-month component.
func (dt LocalDateTime) DayOfMonth() int {
	return dt.date.DayOfMonth()
}

// DayOfWeek returns the day-of-week.
func (dt LocalDateTime) DayOfWeek() DayOfWeek {
	return dt.date.DayOfWeek()
}

// DayOfYear returns the day-of-year.
func (dt LocalDateTime) DayOfYear() int {
	return dt.date.DayOfYear()
}

// Hour returns the hour component (0-23).
func (dt LocalDateTime) Hour() int {
	return dt.time.Hour()
}

// Minute returns the minute component (0-59).
func (dt LocalDateTime) Minute() int {
	return dt.time.Minute()
}

// Second returns the second component (0-59).
func (dt LocalDateTime) Second() int {
	return dt.time.Second()
}

// Millisecond returns the millisecond component (0-999).
func (dt LocalDateTime) Millisecond() int {
	return dt.time.Millisecond()
}

// Nanosecond returns the nanosecond component (0-999999999).
func (dt LocalDateTime) Nanosecond() int {
	return dt.time.Nano()
}

// IsZero returns true if this is the zero value.
func (dt LocalDateTime) IsZero() bool {
	return dt.date.IsZero() && dt.time.IsZero()
}

// IsLeapYear returns true if the year is a leap year.
func (dt LocalDateTime) IsLeapYear() bool {
	return dt.date.IsLeapYear()
}

// IsSupportedField returns true if the field is supported by LocalDateTime.
// LocalDateTime supports all fields from both LocalDate and LocalTime.
func (dt LocalDateTime) IsSupportedField(field Field) bool {
	return dt.date.IsSupportedField(field) || dt.time.IsSupportedField(field)
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the date-time for the value of the specified field.
// The returned value may be unsupported if the field is not supported by LocalDateTime.
//
// If the date-time is zero (IsZero() returns true), an unsupported TemporalValue is returned.
//
// LocalDateTime supports all fields from both LocalDate and LocalTime.
// For fields that are supported by the underlying LocalDate or LocalTime,
// this method delegates to the appropriate component.
//
// Supported fields include all date fields (FieldDayOfWeek, FieldDayOfMonth, FieldDayOfYear, FieldMonthOfYear,
// FieldYear, FieldYearOfEra, FieldEra, FieldEpochDay, FieldProlepticMonth) and all time fields (FieldNanoOfSecond,
// FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay,
// FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfDay, FieldClockHourOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldAmPmOfDay).
//
// Overflow Analysis:
// LocalDateTime delegates to LocalDate and LocalTime, both of which have no overflow issues
// for their supported fields. Therefore, LocalDateTime.GetField cannot overflow:
//   - All date fields are handled by LocalDate.GetField (see LocalDate overflow analysis)
//   - All time fields are handled by LocalTime.GetField (see LocalTime overflow analysis)
//   - No LocalDateTime-specific fields exist that could cause overflow
func (dt LocalDateTime) GetField(field Field) TemporalValue {
	if dt.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}

	// Delegate to LocalDate for date-based fields
	// LocalDate handles: FieldDayOfWeek, FieldDayOfMonth, FieldDayOfYear, FieldMonthOfYear,
	// FieldYear, FieldYearOfEra, FieldEra, FieldEpochDay, FieldProlepticMonth
	if dt.date.IsSupportedField(field) {
		return dt.date.GetField(field)
	}

	// Delegate to LocalTime for time-based fields
	// LocalTime handles: FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay,
	// FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour,
	// FieldMinuteOfDay, FieldHourOfDay, FieldClockHourOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldAmPmOfDay
	if dt.time.IsSupportedField(field) {
		return dt.time.GetField(field)
	}

	// Unsupported field (e.g., FieldInstantSeconds, FieldOffsetSeconds)
	return TemporalValue{unsupported: true}
}

// GoTime converts this date-time to a time.Time in UTC.
// Returns time.Time{} (zero) for zero value.
func (dt LocalDateTime) GoTime() time.Time {
	if dt.IsZero() {
		return time.Time{}
	}
	return time.Date(int(dt.Year()), time.Month(dt.Month()), dt.DayOfMonth(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), time.UTC)
}

// Compare compares this date-time with another.
// Returns -1 if this is before other, 0 if equal, 1 if after.
func (dt LocalDateTime) Compare(other LocalDateTime) int {
	return doCompare(dt, other, compareZero, comparing1(LocalDateTime.LocalDate), comparing1(LocalDateTime.LocalTime))
}

// IsBefore returns true if this date-time is before the specified date-time.
func (dt LocalDateTime) IsBefore(other LocalDateTime) bool {
	return dt.Compare(other) < 0
}

// IsAfter returns true if this date-time is after the specified date-time.
func (dt LocalDateTime) IsAfter(other LocalDateTime) bool {
	return dt.Compare(other) > 0
}

func (dt LocalDateTime) Chain() (chain LocalDateTimeChain) {
	chain.value = dt
	return
}

func (dt LocalDateTime) chainWithError(e error) (chain LocalDateTimeChain) {
	chain = dt.Chain()
	chain.eError = e
	return
}

// LocalDateTimeOf creates a new LocalDateTime from individual components.
// Returns an error if any component is invalid.
func LocalDateTimeOf(year Year, month Month, day, hour, minute, second, nanosecond int) (LocalDateTime, error) {
	date, err := LocalDateOf(year, month, day)
	if err != nil {
		return LocalDateTime{}, err
	}
	time, err := LocalTimeOf(hour, minute, second, nanosecond)
	if err != nil {
		return LocalDateTime{}, err
	}
	return LocalDateTime{
		date: date,
		time: time,
	}, nil
}

// MustLocalDateTimeOf creates a new LocalDateTime from individual components.
// Panics if any component is invalid.
func MustLocalDateTimeOf(year Year, month Month, day, hour, minute, second, nanosecond int) LocalDateTime {
	return mustValue(LocalDateTimeOf(year, month, day, hour, minute, second, nanosecond))
}

func LocalDateTimeOfEpochSecond(epochSecond int64, nanoOfSecond int64, offset ZoneOffset) (r LocalDateTime, e error) {
	FieldNanoOfSecond.checkSetE(nanoOfSecond, &e)
	if e != nil {
		return
	}
	var localSecond = epochSecond + int64(offset.TotalSeconds())
	var localEpochDay = floorDiv(localSecond, 86400)
	var secsOfDay = floorMod(localSecond, 86400)
	r.date, e = LocalDateOfEpochDays(localEpochDay)
	if e != nil {
		return
	}
	r.time, e = LocalTimeOfNanoOfDay(secsOfDay*time.Second.Nanoseconds() + nanoOfSecond)
	return
}

// LocalDateTimeNow returns the current date-time in the system's local time zone.
func LocalDateTimeNow() LocalDateTime {
	return LocalDateTimeOfGoTime(time.Now())
}

// LocalDateTimeNowIn returns the current date-time in the specified time zone.
func LocalDateTimeNowIn(loc *time.Location) LocalDateTime {
	return LocalDateTimeOfGoTime(time.Now().In(loc))
}

// LocalDateTimeNowUTC returns the current date-time in UTC.
func LocalDateTimeNowUTC() LocalDateTime {
	return LocalDateTimeOfGoTime(time.Now().UTC())
}

// LocalDateTimeOfGoTime creates a LocalDateTime from a time.Time.
// Returns zero value if t.IsZero().
func LocalDateTimeOfGoTime(t time.Time) LocalDateTime {
	if t.IsZero() {
		return LocalDateTime{}
	}
	return LocalDateOfGoTime(t).AtTime(LocalTimeOfGoTime(t))
}

// LocalDateTimeParse parses a date-time string in RFC3339-compatible format.
// The date must be in yyyy-MM-dd form, and the time must be in HH:mm:ss or
// HH:mm:ss[.nnnnnnnnn] form.
//
// The separator between the date and time may be 'T', 't', or a single space.
//
// Examples:
//
//	dt, err := LocalDateTimeParse("2024-03-15T14:30:45.123456789")
//	dt, err := LocalDateTimeParse("2024-03-15 14:30:45")
//	dt, err := LocalDateTimeParse("2024-03-15t14:30:45")
func LocalDateTimeParse(s string) (LocalDateTime, error) {
	var dt LocalDateTime
	err := dt.UnmarshalText([]byte(s))
	return dt, err
}

// MustLocalDateTimeParse parses a date-time string in yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn] format.
// Panics if the string is invalid.
func MustLocalDateTimeParse(s string) LocalDateTime {
	return mustValue(LocalDateTimeParse(s))
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*LocalDateTime)(nil)
	_ fmt.Stringer             = (*LocalDateTime)(nil)
	_ encoding.TextMarshaler   = (*LocalDateTime)(nil)
	_ encoding.TextUnmarshaler = (*LocalDateTime)(nil)
	_ json.Marshaler           = (*LocalDateTime)(nil)
	_ json.Unmarshaler         = (*LocalDateTime)(nil)
	_ driver.Valuer            = (*LocalDateTime)(nil)
	_ sql.Scanner              = (*LocalDateTime)(nil)
)

// Compile-time check that LocalDateTime is comparable
func _assertLocalDateTimeIsComparable[T comparable](t T) {}

var _ = _assertLocalDateTimeIsComparable[LocalDateTime]

var localDateTimeMinEpochSecond = mustValue(LocalDateTimeOfEpochSecond(math.MinInt64, 0, ZoneOffsetUTC()))
var localDateTimeMaxEpochSecond = mustValue(LocalDateTimeOfEpochSecond(math.MaxInt64, 0, ZoneOffsetUTC()))
