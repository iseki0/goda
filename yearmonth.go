package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
)

type YearMonth struct {
	v int64
}

func (y YearMonth) IsZero() bool {
	return y.v == 0
}

func (y YearMonth) Compare(other YearMonth) int {
	return doCompare(y, other, compareZero, comparing(YearMonth.Year), comparing(YearMonth.Month))
}

func (y YearMonth) IsLeapYear() bool {
	return y.Year().IsLeapYear()
}

func (y YearMonth) LengthOfMonth() int {
	return y.Month().Length(y.IsLeapYear())
}

func (y YearMonth) LengthOfYear() int {
	return y.Year().Length()
}

func (y YearMonth) Year() Year {
	return Year(y.v >> 16)
}

func (y YearMonth) ProlepticMonth() int64 {
	return y.Year().Int64()*12 + int64(y.Month()) - 1
}

func (y YearMonth) Month() Month {
	return Month(y.v & 0xffff)
}

func (y YearMonth) Chain() (chain YearMonthChain) {
	chain.value = y
	return
}

func YearMonthOf(year Year, month Month) (y YearMonth, e error) {
	FieldYear.checkSetE(year.Int64(), &e)
	FieldMonthOfYear.checkSetE(int64(month), &e)
	if e != nil {
		return
	}
	return YearMonth{int64(year)<<16 | int64(month)}, nil
}

func MustYearMonthOf(year Year, month Month) YearMonth {
	return mustValue(YearMonthOf(year, month))
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*YearMonth)(nil)
	_ fmt.Stringer             = (*YearMonth)(nil)
	_ encoding.TextMarshaler   = (*YearMonth)(nil)
	_ encoding.TextUnmarshaler = (*YearMonth)(nil)
	_ json.Marshaler           = (*YearMonth)(nil)
	_ json.Unmarshaler         = (*YearMonth)(nil)
	_ driver.Valuer            = (*YearMonth)(nil)
	_ sql.Scanner              = (*YearMonth)(nil)
)

// Compile-time check that YearMonth is comparable
func _assertYearMonthIsComparable[T comparable](t T) {}

var _ = _assertYearMonthIsComparable[YearMonth]
