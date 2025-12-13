package goda

import (
	"encoding"
	"fmt"
)

// Year represents a year in the ISO-8601 calendar system.
// It can represent any year from math.MinInt64 to math.MaxInt64.
type Year int64

const YearMax = 1<<47 - 1
const YearMin = -YearMax - 1

// Int returns this year as an int.
func (y Year) Int() int {
	return int(y)
}

// Int64 returns this year as an int64.
func (y Year) Int64() int64 {
	return int64(y)
}

// IsLeapYear returns true if this year is a leap year.
// A leap year is divisible by 4, unless it's divisible by 100 (but not 400).
// For example: 2024 is a leap year, 1900 is not, but 2000 is.
func (y Year) IsLeapYear() bool {
	return (y%4 == 0 && y%100 != 0) || (y%400 == 0)
}

// Length returns the number of days in this year (365 or 366).
func (y Year) Length() int {
	if y.IsLeapYear() {
		return 366
	}
	return 365
}

var _ encoding.TextAppender = Year(0)
var _ fmt.Stringer = Year(0)
