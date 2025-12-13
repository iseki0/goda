package goda

import (
	"math"
)

type YearMonthChained struct {
	Chain[YearMonth]
}

func (y YearMonthChained) PlusMonths(months int64) YearMonthChained {
	if !y.ok() {
		return y
	}
	if months == 0 {
		return y
	}
	if _, overflow := mulExact(y.value.Year().Int64(), 12); overflow {
		y.Error = overflowError()
		return y
	}
	var monthCount = y.value.Year().Int64()*12 + (int64(y.value.Month()) - 1)
	var calcMonth = monthCount + months
	var newYear = floorDiv(calcMonth, 12)
	var newMonth = floorMod(calcMonth, 12) + 1
	y.value, y.Error = YearMonthOf(Year(newYear), Month(newMonth))
	return y
}

func (y YearMonthChained) MinusMonths(months int64) YearMonthChained {
	if months == math.MinInt64 {
		return y.PlusMonths(math.MaxInt64).PlusMonths(1)
	}
	return y.PlusMonths(-months)
}

func (y YearMonthChained) PlusYears(years int64) YearMonthChained {
	if !y.ok() {
		return y
	}
	newYear, overflow := addExactly(y.value.Year().Int64(), years)
	if !overflow {
		y.Error = overflowError()
	}
	return y.WithField(FieldYear, TemporalValueOf(newYear))
}

func (y YearMonthChained) MinusYears(years int64) YearMonthChained {
	if years == math.MinInt64 {
		return y.PlusYears(math.MaxInt64).PlusYears(1)
	}
	return y.PlusYears(-years)
}

func (y YearMonthChained) WithMonth(month Month) YearMonthChained {
	return y.WithField(FieldMonthOfYear, TemporalValueOf(int64(month)))
}

func (y YearMonthChained) WithYear(year Year) YearMonthChained {
	return y.WithField(FieldYear, TemporalValueOf(int64(year)))
}

func (y YearMonthChained) WithField(field Field, value TemporalValue) YearMonthChained {
	field.checkSetE(value.Int64(), &y.Error)
	if !y.ok() {
		return y
	}
	if field == FieldProlepticMonth {
		return y.PlusMonths(value.Int64() - y.value.ProlepticMonth())
	}
	var year = y.value.Year()
	var month = y.value.Month()
	switch field {
	case FieldYear:
		year = Year(value.v)
	case FieldYearOfEra:
		if year >= 1 {
			// CE
			year = Year(value.v)
		} else {
			// BCE: convert YearOfEra back to negative year
			year = Year(-(value.v - 1))
		}
	case FieldMonthOfYear:
		month = Month(value.v)
	default:
		y.Error = unsupportedField(field)
		return y
	}
	y.value, y.Error = YearMonthOf(year, month)
	return y
}
