package goda

import (
	"math"
)

type YearMonthChain struct {
	Chain[YearMonth]
}

func (y YearMonthChain) PlusMonths(months int64) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "PlusMonths"))
	if !y.ok() {
		return y
	}
	if months == 0 {
		return y
	}
	if _, overflow := mulExact(y.value.Year().Int64(), 12); overflow {
		y.eError = overflowError()
		return y
	}
	var monthCount = y.value.Year().Int64()*12 + (int64(y.value.Month()) - 1)
	var calcMonth = monthCount + months
	var newYear = floorDiv(calcMonth, 12)
	var newMonth = floorMod(calcMonth, 12) + 1
	y.value, y.eError = YearMonthOf(Year(newYear), Month(newMonth))
	return y
}

func (y YearMonthChain) MinusMonths(months int64) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "MinusMonths"))
	if months == math.MinInt64 {
		return y.PlusMonths(math.MaxInt64).PlusMonths(1)
	}
	return y.PlusMonths(-months)
}

func (y YearMonthChain) PlusYears(years int64) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "PlusYears"))
	if !y.ok() {
		return y
	}
	newYear, overflow := addExactly(y.value.Year().Int64(), years)
	if overflow {
		y.eError = overflowError()
	}
	return y.WithField(FieldYear, TemporalValueOf(newYear))
}

func (y YearMonthChain) MinusYears(years int64) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "MinusYears"))
	if years == math.MinInt64 {
		return y.PlusYears(math.MaxInt64).PlusYears(1)
	}
	return y.PlusYears(-years)
}

func (y YearMonthChain) WithMonth(month Month) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "WithMonth"))
	return y.WithField(FieldMonthOfYear, TemporalValueOf(int64(month)))
}

func (y YearMonthChain) WithYear(year Year) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "WithYear"))
	return y.WithField(FieldYear, TemporalValueOf(int64(year)))
}

func (y YearMonthChain) WithField(field Field, value TemporalValue) YearMonthChain {
	defer y.leaveFunction(y.enterFunction("YearMonth", "WithField"))
	field.checkSetE(value.Int64(), &y.eError)
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
		y.eError = unsupportedField(field)
		return y
	}
	y.value, y.eError = YearMonthOf(year, month)
	return y
}
