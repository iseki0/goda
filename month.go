package goda

import (
	"strconv"
	"time"
)

// Month represents a month-of-year in the ISO-8601 calendar system,
// such as January or December. It is compatible with time.Month.
type Month time.Month

// Month constants.
const (
	January   Month = iota + 1 // January (month 1)
	February                   // February (month 2)
	March                      // March (month 3)
	April                      // April (month 4)
	May                        // May (month 5)
	June                       // June (month 6)
	July                       // July (month 7)
	August                     // August (month 8)
	September                  // September (month 9)
	October                    // October (month 10)
	November                   // November (month 11)
	December                   // December (month 12)
)

// IsZero returns true if this is the zero value (not a valid month).
func (m Month) IsZero() bool {
	return m == 0
}

// FirstDayOfYear returns the day-of-year (1-366) for the first day of this month.
// The isLeap parameter indicates whether the year is a leap year.
// For example, March 1st is day 60 in non-leap years and day 61 in leap years.
func (m Month) FirstDayOfYear(isLeap bool) int {
	var d int
	switch m {
	case 0:
		return 0
	case January:
		d = 0 + 1
	case February:
		d = 31 + 1
	case March:
		d = 59 + 1
	case April:
		d = 90 + 1
	case May:
		d = 120 + 1
	case June:
		d = 151 + 1
	case July:
		d = 181 + 1
	case August:
		d = 212 + 1
	case September:
		d = 243 + 1
	case October:
		d = 273 + 1
	case November:
		d = 304 + 1
	case December:
		d = 334 + 1
	default:
		panic("invalid month: " + strconv.Itoa(int(m)))
	}
	if isLeap && m >= March {
		d += 1
	}
	return d
}

// MaxDays returns the maximum days in the month.
// For February, it returns 29 (leap year maximum).
// For other months, it returns the fixed number of days.
func (m Month) MaxDays() int {
	switch m {
	case 0:
		return 0
	case January, March, May, July, August, October, December:
		return 31
	case April, June, September, November:
		return 30
	case February:
		return 29
	default:
		panic("invalid month: " + strconv.Itoa(int(m)))
	}
}

// Length returns the number of days in this month for the specified year type.
// The isLeap parameter indicates whether the year is a leap year.
// February returns 29 for leap years and 28 for non-leap years.
// Other months return their fixed number of days.
func (m Month) Length(isLeap bool) int {
	if m == February {
		if isLeap {
			return 29
		}
		return 28
	}
	return m.MaxDays()
}
