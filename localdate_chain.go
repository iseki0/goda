package goda

import (
	"math"
)

type LocalDateChain struct {
	Chain[LocalDate]
}

var localDateMaxEpoch = LocalDateMax().UnixEpochDays()
var localDateMinEpoch = LocalDateMin().UnixEpochDays()

func (l LocalDateChain) PlusDays(days int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "PlusDays"))
	if !l.ok() {
		return l
	}
	if days == 0 {
		return l
	}
	r, overflow := addExactly(l.value.UnixEpochDays(), days)
	if overflow || r > localDateMaxEpoch || r < localDateMinEpoch {
		l.eError = overflowError()
		return l
	}
	l.value, l.eError = LocalDateOfEpochDays(r)
	return l
}

func (l LocalDateChain) MinusDays(days int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "MinusDays"))
	if days == math.MinInt64 {
		return l.PlusDays(math.MaxInt64).PlusDays(1)
	}
	return l.PlusDays(-days)
}

func (l LocalDateChain) PlusMonths(months int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "PlusMonths"))
	if !l.ok() {
		return l
	}
	if months == 0 {
		return l
	}
	newYearMonth := l.value.YearMonth().Chain().PlusMonths(months).mergeError(&l.eError)
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalDateOf(newYearMonth.Year(), newYearMonth.Month(), min(l.value.DayOfMonth(), newYearMonth.LengthOfMonth()))
	return l
}

func (l LocalDateChain) MinusMonths(months int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "MinusMonths"))
	if months == math.MinInt64 {
		return l.PlusMonths(math.MaxInt64).PlusMonths(1)
	}
	return l.PlusMonths(-months)
}

func (l LocalDateChain) PlusYears(years int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "PlusYears"))
	if !l.ok() {
		return l
	}
	if years == 0 {
		return l
	}
	newYearMonth := l.value.YearMonth().Chain().PlusYears(years).mergeError(&l.eError)
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalDateOf(newYearMonth.Year(), newYearMonth.Month(), min(l.value.DayOfMonth(), newYearMonth.LengthOfMonth()))
	return l
}

func (l LocalDateChain) MinusYears(years int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "MinusYears"))
	if years == math.MinInt64 {
		return l.PlusYears(math.MaxInt64).PlusYears(1)
	}
	return l.PlusYears(-years)
}

func (l LocalDateChain) PlusWeeks(weeks int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "PlusWeeks"))
	if !l.ok() {
		return l
	}
	if weeks == 0 {
		return l
	}
	r, overflow := mulExact(7, weeks)
	if overflow {
		l.eError = overflowError()
		return l
	}
	return l.PlusDays(r)
}

func (l LocalDateChain) MinusWeeks(weeks int64) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "MinusWeeks"))
	if weeks == math.MinInt64 {
		return l.PlusWeeks(math.MaxInt64).PlusWeeks(1)
	}
	return l.PlusWeeks(-weeks)
}

func (l LocalDateChain) WithDayOfMonth(dayOfMonth int) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "WithDayOfMonth"))
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalDateOf(l.value.Year(), l.value.Month(), dayOfMonth)
	return l
}

func (l LocalDateChain) WithDayOfYear(dayOfYear int) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "WithDayOfYear"))
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalDateOfYearDay(l.value.Year(), dayOfYear)
	return l
}

func (l LocalDateChain) WithMonth(month Month) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "WithMonth"))
	l.eError = FieldMonthOfYear.check(int64(month))
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalDateOf(l.value.Year(), month, min(MustYearMonthOf(l.value.Year(), month).LengthOfMonth(), l.value.DayOfMonth()))
	return l
}

func (l LocalDateChain) WithYear(year Year) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "WithYear"))
	l.eError = FieldYear.check(int64(year))
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalDateOf(year, l.value.Month(), min(MustYearMonthOf(year, l.value.Month()).LengthOfMonth(), l.value.DayOfMonth()))
	return l
}

// WithField returns a copy of this LocalDate with the specified field replaced.
// Zero values return zero immediately.
//
// Supported fields mirror Java's LocalDate#with(TemporalField, long):
//   - FieldDayOfMonth: sets the day-of-month while keeping year and month.
//   - FieldDayOfYear: sets the day-of-year while keeping year.
//   - FieldMonthOfYear: sets the month-of-year while keeping year and day-of-month (adjusted if necessary).
//   - FieldYear: sets the year while keeping month and day-of-month (adjusted if necessary).
//   - FieldYearOfEra: sets the year within the current era (same as FieldYear for CE dates).
//   - FieldEra: switches between BCE/CE eras while preserving year, month, and day.
//   - FieldEpochDay: sets the date based on days since Unix epoch (1970-01-01).
//   - FieldProlepticMonth: sets the date based on months since year 0.
//
// Fields outside this list return an error. Range violations propagate the validation error.
func (l LocalDateChain) WithField(field Field, value TemporalValue) LocalDateChain {
	defer l.leaveFunction(l.enterFunction("LocalDate", "WithField"))
	field.checkSetE(value.Int64(), &l.eError)
	if !l.ok() {
		return l
	}
	newValue := value.v
	switch field {
	case FieldDayOfWeek:
		return l.PlusDays(newValue - int64(l.value.DayOfWeek()))
	case FieldAlignedDayOfWeekInMonth:
		return l.PlusDays(newValue - l.value.GetField(FieldAlignedDayOfWeekInMonth).Int64())
	case FieldAlignedDayOfWeekInYear:
		return l.PlusDays(newValue - l.value.GetField(FieldAlignedDayOfWeekInYear).Int64())
	case FieldDayOfMonth:
		return l.WithDayOfMonth(int(newValue))
	case FieldDayOfYear:
		return l.WithDayOfYear(int(newValue))
	case FieldEpochDay:
		l.value, l.eError = LocalDateOfEpochDays(newValue)
	case FieldAlignedWeekOfMonth:
		return l.PlusWeeks(newValue - l.value.GetField(FieldAlignedWeekOfMonth).Int64())
	case FieldAlignedWeekOfYear:
		return l.PlusWeeks(newValue - l.value.GetField(FieldAlignedWeekOfYear).Int64())
	case FieldMonthOfYear:
		return l.WithMonth(Month(newValue))
	case FieldProlepticMonth:
		return l.PlusMonths(newValue - l.value.GetField(FieldProlepticMonth).Int64())
	case FieldYearOfEra:
		if l.value.Year() >= 1 {
			return l.WithYear(Year(newValue))
		}
		return l.WithYear(Year(1 - newValue))
	case FieldYear:
		return l.WithYear(Year(newValue))
	case FieldEra:
		if l.value.GetField(FieldEra).Int64() == newValue {
			return l
		}
		return l.WithYear(1 - l.value.Year())
	default:
		l.eError = unsupportedField(field)
	}
	return l
}
