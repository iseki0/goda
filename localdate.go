package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// LocalDate represents a date without a time zone in the Gregorian calendar system,
// such as 2024-03-15. It stores the year, month, and day-of-month.
//
// LocalDate is comparable and can be used as a map key.
// The zero value represents an unset date and IsZero returns true for it.
//
// LocalDate implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: yyyy-MM-dd (e.g., "2024-03-15"). Uses ISO 8601 basic calendar date format.
// Week dates (YYYY-Www-D) and ordinal dates (YYYY-DDD) are not supported.
type LocalDate struct {
	v int64
}

// Year returns the year component of this date.
func (d LocalDate) Year() Year {
	return Year(d.v >> 16)
}

// Month returns the month component of this date (1-12).
func (d LocalDate) Month() Month {
	return Month(d.v >> 8 & 0xff)
}

// DayOfMonth returns the day-of-month component of this date (1-31).
func (d LocalDate) DayOfMonth() int {
	return int(d.v & 0xff)
}

func (d LocalDate) YearMonth() YearMonth {
	if d.IsZero() {
		return YearMonth{}
	}
	return MustYearMonthOf(d.Year(), d.Month())
}

// DayOfWeek returns the day-of-week for this date.
// Returns 0 for zero value, otherwise Monday=1 through Sunday=7.
func (d LocalDate) DayOfWeek() DayOfWeek {
	if d.IsZero() {
		return 0
	}
	return DayOfWeek(floorMod(d.UnixEpochDays()+3, 7) + 1)
}

// DayOfYear returns the day-of-year for this date (1-366).
// Returns 0 for zero value.
func (d LocalDate) DayOfYear() int {
	if d.IsZero() {
		return 0
	}
	return d.Month().FirstDayOfYear(d.IsLeapYear()) - 1 + d.DayOfMonth()
}

// IsLeapYear returns true if the year of this date is a leap year.
// A leap year is divisible by 4, unless it's divisible by 100 (but not 400).
func (d LocalDate) IsLeapYear() bool {
	return d.Year().IsLeapYear()
}

// Compare compares this date with another date.
// Returns -1 if this date is before other, 0 if equal, and 1 if after.
// Zero values are considered less than non-zero values.
func (d LocalDate) Compare(other LocalDate) int {
	return doCompare(d, other, compareZero, comparing(LocalDate.Year), comparing(LocalDate.Month), comparing(LocalDate.DayOfMonth))
}

// IsBefore returns true if this date is before the specified date.
func (d LocalDate) IsBefore(other LocalDate) bool {
	return d.Compare(other) < 0
}

// IsAfter returns true if this date is after the specified date.
func (d LocalDate) IsAfter(other LocalDate) bool {
	return d.Compare(other) > 0
}

// GoTime converts this date to a time.Time at midnight UTC.
// Returns time.Time{} (zero) for zero value.
func (d LocalDate) GoTime() time.Time {
	if d.IsZero() {
		return time.Time{}
	}
	return time.Date(int(d.Year()), time.Month(d.Month()), d.DayOfMonth(), 0, 0, 0, 0, time.UTC)
}

func (d LocalDate) Chain() (chain LocalDateChain) {
	chain.value = d
	return
}

func (d LocalDate) chainWithError(e error) (chain LocalDateChain) {
	chain = d.Chain()
	chain.eError = e
	return
}

// UnixEpochDays returns the number of days since Unix epoch (1970-01-01).
// Positive values represent dates after the epoch, negative before.
// Returns 0 for zero value.
func (d LocalDate) UnixEpochDays() int64 {
	if d.IsZero() {
		return 0
	}
	const DaysPerCycle = 365*400 + 97
	const Days0000To1970 = (DaysPerCycle * 5) - (30*365 + 7)

	y := d.Year().Int64()
	m := int64(d.Month())
	day := int64(d.DayOfMonth())
	total := int64(0)

	// Calculate year contribution
	total += 365 * y
	if y >= 0 {
		total += (y+3)/4 - (y+99)/100 + (y+399)/400
	} else {
		total -= y/-4 - y/-100 + y/-400
	}

	// Calculate month contribution
	total += (367*m - 362) / 12

	// Calculate day contribution
	total += day - 1

	// Adjust for leap year if month > February
	if m > 2 {
		total--
		if !d.Year().IsLeapYear() {
			total--
		}
	}

	return total - Days0000To1970
}

// IsSupportedField returns true if the field is supported by LocalDate.
func (d LocalDate) IsSupportedField(field Field) bool {
	switch field {
	case FieldDayOfWeek, FieldDayOfMonth, FieldDayOfYear, FieldEpochDay, FieldMonthOfYear, FieldProlepticMonth, FieldYearOfEra, FieldYear, FieldEra:
		return true
	default:
		return false
	}
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the date for the value of the specified field.
// The returned value may be unsupported if the field is not supported by LocalDate.
//
// If the date is zero (IsZero() returns true), an unsupported TemporalValue is returned.
// For fields not supported by LocalDate (such as time fields), an unsupported TemporalValue is returned.
//
// Supported fields include:
//   - FieldDayOfWeek: returns the day of week (1=Monday, 7=Sunday)
//   - FieldDayOfMonth: returns the day of month (1-31)
//   - FieldDayOfYear: returns the day of year (1-366)
//   - FieldMonthOfYear: returns the month (1=January, 12=December)
//   - FieldYear: returns the proleptic year
//   - FieldYearOfEra: returns the year within the era (same as FieldYear for CE dates)
//   - FieldEra: returns the era (0=BCE, 1=CE)
//   - FieldEpochDay: returns the number of days since Unix epoch (1970-01-01)
//   - FieldProlepticMonth: returns the number of months since year 0
//
// Overflow Analysis:
// None of the supported fields can overflow int64 in practice:
//   - FieldDayOfWeek: range 1-7, cannot overflow
//   - FieldDayOfMonth: range 1-31, cannot overflow
//   - FieldDayOfYear: range 1-366, cannot overflow
//   - FieldMonthOfYear: range 1-12, cannot overflow
//   - FieldYear/FieldYearOfEra: Year is int64, direct cast, cannot overflow
//   - FieldEra: values 0 or 1, cannot overflow
//   - FieldEpochDay: int64, calculated from year/month/day which are bounded by LocalDate's internal representation
//   - FieldProlepticMonth: Year * 12 + Month.
//     Year is int64, so max value is approximately int64_max * 12,
//     which would overflow.
//     However, LocalDate stores Year in the upper 48 bits of a 64-bit value,
//     limiting the practical range to approximately ±140 trillion years, making overflow impossible
//     in any realistic scenario.
func (d LocalDate) GetField(field Field) TemporalValue {
	if d.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}
	var v int64
	switch field {
	case FieldDayOfWeek:
		// Range: 1-7, no overflow possible
		v = int64(d.DayOfWeek())
	case FieldDayOfMonth:
		// Range: 1-31, no overflow possible
		v = int64(d.DayOfMonth())
	case FieldDayOfYear:
		// Range: 1-366, no overflow possible
		v = int64(d.DayOfYear())
	case FieldMonthOfYear:
		// Range: 1-12, no overflow possible
		v = int64(d.Month())
	case FieldYear:
		// Year is already int64, direct cast, no overflow possible
		v = int64(d.Year())
	case FieldYearOfEra:
		v = int64(d.Year())
		if v < 0 {
			v = -v + 1
		}
	case FieldEra:
		// Values: 0 or 1, no overflow possible
		if d.Year() >= 1 {
			v = 1 // CE (Common FieldEra)
		} else {
			v = 0 // BCE (Before Common FieldEra)
		}
	case FieldEpochDay:
		// Returns int64, computed from bounded year/month/day values, no overflow in practice
		v = d.UnixEpochDays()
	case FieldProlepticMonth:
		// Calculate proleptic month (months since year 0)
		// Year is stored in 48 bits internally, limiting range to ±140 trillion years
		// Year * 12 cannot overflow int64 in this constrained range
		v = int64(d.Year())*12 + int64(d.Month()) - 1
	default:
		return TemporalValue{unsupported: true}
	}
	return TemporalValue{v: v}
}

// AtTime combines this date with a time to create a LocalDateTime.
func (d LocalDate) AtTime(time LocalTime) LocalDateTime {
	return LocalDateTime{
		date: d,
		time: time,
	}
}

// LengthOfMonth returns the number of days in the month of this date.
// Returns 28, 29, 30, or 31 depending on the month and whether it's a leap year.
// Returns 0 for zero value.
func (d LocalDate) LengthOfMonth() int {
	if d.IsZero() {
		return 0
	}
	return d.Month().Length(d.Year().IsLeapYear())
}

// LengthOfYear returns the number of days in the year of this date.
// Returns 365 for non-leap years and 366 for leap years.
// Returns 0 for zero value.
func (d LocalDate) LengthOfYear() int {
	if d.IsZero() {
		return 0
	}
	return d.Year().Length()
}

// IsZero returns true if this is the zero value of LocalDate.
func (d LocalDate) IsZero() bool {
	return d.v == 0
}

// LocalDateOf creates a new LocalDate from the specified year, month, and day-of-month.
// Returns an error if the date is invalid (e.g., month out of range 1-12,
// day out of range for the month, or February 29 in a non-leap year).
func LocalDateOf(year Year, month Month, dayOfMonth int) (d LocalDate, e error) {
	_, e = YearMonthOf(year, month)
	FieldDayOfMonth.checkSetE(int64(dayOfMonth), &e)
	if e != nil {
		return
	}
	if dayOfMonth > month.Length(year.IsLeapYear()) {
		if dayOfMonth == 29 {
			e = newError("invalid date February 29 in non-leap year")
		} else {
			e = newError("invalid date %s %d", month, dayOfMonth)
		}
		return
	}
	d = LocalDate{
		v: int64(year)<<16 | int64(month)<<8 | int64(dayOfMonth),
	}
	return
}

// MustLocalDateOf creates a new LocalDate from the specified year, month, and day-of-month.
// Panics if the date is invalid.
// Use LocalDateOf for error handling.
func MustLocalDateOf(year Year, month Month, dayOfMonth int) LocalDate {
	return mustValue(LocalDateOf(year, month, dayOfMonth))
}

func LocalDateOfYearDay(year Year, dayOfYear int) (r LocalDate, e error) {
	FieldYear.checkSetE(year.Int64(), &e)
	FieldDayOfYear.checkSetE(int64(dayOfYear), &e)
	if e != nil {
		return
	}
	leap := year.IsLeapYear()
	if dayOfYear == 366 && !leap {
		e = newError("invalid date DayOfYear 366 in non-leap year")
	}
	moy := Month((dayOfYear-1)/31 + 1)
	var monthEnd = moy.FirstDayOfYear(leap) + moy.Length(leap) - 1
	if dayOfYear > monthEnd {
		moy = moy + 1
	}
	var dom = dayOfYear - moy.FirstDayOfYear(leap) + 1
	return LocalDateOf(year, moy, dom)
}

// LocalDateOfGoTime creates a LocalDate from a time.Time.
// The time zone and time-of-day components are ignored.
// Returns zero value if t.IsZero().
func LocalDateOfGoTime(t time.Time) LocalDate {
	if t.IsZero() {
		return LocalDate{}
	}
	return MustLocalDateOf(Year(t.Year()), Month(t.Month()), t.Day())
}

// LocalDateNow returns the current date in the system's local time zone.
// This is equivalent to LocalDateOfGoTime(time.Now()).
// For UTC time, use LocalDateNowUTC.
// For a specific timezone, use LocalDateNowIn.
func LocalDateNow() LocalDate {
	return LocalDateOfGoTime(time.Now())
}

// LocalDateNowIn returns the current date in the specified time zone.
// This is equivalent to LocalDateOfGoTime(time.Now().In(loc)).
func LocalDateNowIn(loc *time.Location) LocalDate {
	return LocalDateOfGoTime(time.Now().In(loc))
}

// LocalDateNowUTC returns the current date in UTC.
// This is equivalent to LocalDateOfGoTime(time.Now().UTC()).
func LocalDateNowUTC() LocalDate {
	return LocalDateOfGoTime(time.Now().UTC())
}

// LocalDateParse parses a date string in yyyy-MM-dd format.
// Returns an error if the string is invalid or represents an invalid date.
//
// Supported format: yyyy-MM-dd (e.g., "2024-03-15")
// Week dates and ordinal dates are not supported.
//
// Example:
//
//	date, err := LocalDateParse("2024-03-15")
//	if err != nil {
//	    // handle error
//	}
func LocalDateParse(s string) (LocalDate, error) {
	var d LocalDate
	err := d.UnmarshalText([]byte(s))
	return d, err
}

// MustLocalDateParse parses a date string in yyyy-MM-dd format.
// Panics if the string is invalid.
// Use LocalDateParse for error handling.
//
// Example:
//
//	date := MustLocalDateParse("2024-03-15")
func MustLocalDateParse(s string) LocalDate {
	return mustValue(LocalDateParse(s))
}

func MustLocalDateOfUnixEpochDays(days int64) LocalDate {
	return mustValue(LocalDateOfEpochDays(days))
}

// LocalDateOfEpochDays creates a LocalDate from the number of days since Unix epoch (1970-01-01).
// Positive values represent dates after the epoch, negative before.
func LocalDateOfEpochDays(days int64) (LocalDate, error) {
	const DaysPerCycle = 365*400 + 97
	const Days0000To1970 = (DaysPerCycle * 5) - (30*365 + 7)
	zeroDay := days + Days0000To1970

	// Adjust to March-based year (makes leap day fall at end of 4-year cycle)
	zeroDay -= 60 // Shift to 0000-03-01 as base

	var adjust int64
	if zeroDay < 0 {
		// For negative years, adjust to positive for calculation
		adjustCycles := (zeroDay+1)/DaysPerCycle - 1
		adjust = adjustCycles * 400
		zeroDay += -adjustCycles * DaysPerCycle
	}

	// Estimate the year
	yearEst := (400*zeroDay + 591) / DaysPerCycle
	// Calculate day-of-year for estimated year
	doyEst := zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	if doyEst < 0 {
		// Correct the estimate
		yearEst--
		doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	}

	yearEst += adjust // Restore negative year adjustment

	// Convert from March-based day-of-year
	marchDoy0 := int(doyEst)

	// Convert back to January-based month and day
	marchMonth0 := (marchDoy0*5 + 2) / 153
	month := Month(marchMonth0 + 3)
	if month > 12 {
		month -= 12
	}
	dom := marchDoy0 - (marchMonth0*306+5)/10 + 1
	if marchDoy0 >= 306 {
		yearEst++
	}
	return LocalDateOf(Year(yearEst), month, dom)
}

func LocalDateMin() LocalDate {
	return MustLocalDateOf(YearMin, January, 1)
}

func LocalDateMax() LocalDate {
	return MustLocalDateOf(YearMax, December, 31)
}

var (
	_ encoding.TextAppender    = (*LocalDate)(nil)
	_ fmt.Stringer             = (*LocalDate)(nil)
	_ encoding.TextMarshaler   = (*LocalDate)(nil)
	_ encoding.TextUnmarshaler = (*LocalDate)(nil)
	_ json.Marshaler           = (*LocalDate)(nil)
	_ json.Unmarshaler         = (*LocalDate)(nil)
	_ driver.Valuer            = (*LocalDate)(nil)
	_ sql.Scanner              = (*LocalDate)(nil)
)

// Compile-time check that LocalDate is comparable
func _assertLocalDateIsComparable[T comparable](t T) {}

var _ = _assertLocalDateIsComparable[LocalDate]
