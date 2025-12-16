package goda

import (
	"math"
	"strconv"
)

// Field represents a date-time field, such as month-of-year or hour-of-day.
// This is similar to Java's ChronoField.
type Field int

// Field constants representing date and time components.
const (
	// FieldNanoOfSecond represents the nano-of-second field (0-999,999,999).
	FieldNanoOfSecond Field = iota + 1

	// FieldNanoOfDay represents the nano-of-day field (0-86,399,999,999,999).
	FieldNanoOfDay

	// FieldMicroOfSecond represents the micro-of-second field (0-999,999).
	FieldMicroOfSecond

	// FieldMicroOfDay represents the micro-of-day field (0-86,399,999,999).
	FieldMicroOfDay

	// FieldMilliOfSecond represents the milli-of-second field (0-999).
	FieldMilliOfSecond

	// FieldMilliOfDay represents the milli-of-day field (0-86,399,999).
	FieldMilliOfDay

	// FieldSecondOfMinute represents the second-of-minute field (0-59).
	FieldSecondOfMinute

	// FieldSecondOfDay represents the second-of-day field (0-86,399).
	FieldSecondOfDay

	// FieldMinuteOfHour represents the minute-of-hour field (0-59).
	FieldMinuteOfHour

	// FieldMinuteOfDay represents the minute-of-day field (0-1,439).
	FieldMinuteOfDay

	// FieldHourOfAmPm represents the hour-of-am-pm field (0-11).
	FieldHourOfAmPm

	// FieldClockHourOfAmPm represents the clock-hour-of-am-pm field (1-12).
	FieldClockHourOfAmPm

	// FieldHourOfDay represents the hour-of-day field (0-23).
	FieldHourOfDay

	// FieldClockHourOfDay represents the clock-hour-of-day field (1-24).
	FieldClockHourOfDay

	// FieldAmPmOfDay represents the am-pm-of-day field (0=AM, 1=PM).
	FieldAmPmOfDay

	// FieldDayOfWeek represents the day-of-week field (1=Monday, 7=Sunday).
	FieldDayOfWeek

	// FieldAlignedDayOfWeekInMonth represents the aligned day-of-week within a month.
	FieldAlignedDayOfWeekInMonth

	// FieldAlignedDayOfWeekInYear represents the aligned day-of-week within a year.
	FieldAlignedDayOfWeekInYear

	// FieldDayOfMonth represents the day-of-month field (1-31).
	FieldDayOfMonth

	// FieldDayOfYear represents the day-of-year field (1-366).
	FieldDayOfYear

	// FieldEpochDay represents the epoch-day field, based on the Unix epoch of 1970-01-01.
	FieldEpochDay

	// FieldAlignedWeekOfMonth represents the aligned week within a month.
	FieldAlignedWeekOfMonth

	// FieldAlignedWeekOfYear represents the aligned week within a year.
	FieldAlignedWeekOfYear

	// FieldMonthOfYear represents the month-of-year field (1=January, 12=December).
	FieldMonthOfYear

	// FieldProlepticMonth represents the proleptic-month, counting months sequentially from year 0.
	FieldProlepticMonth

	// FieldYearOfEra represents the year within the era.
	FieldYearOfEra

	// FieldYear represents the proleptic year, such as 2024.
	FieldYear

	// FieldEra represents the era field.
	FieldEra

	// FieldInstantSeconds represents the instant epoch-seconds.
	FieldInstantSeconds

	// FieldOffsetSeconds represents the offset from UTC/Greenwich in seconds.
	FieldOffsetSeconds
)

type fieldRange struct {
	Min   int64
	Max   int64
	Valid bool
}

type fieldDescriptor struct {
	name     string
	javaName string
	based    int8
	fieldRange
}

const (
	timeBased = iota + 1
	dateBased
)

func makeRange(min, max int64) fieldRange {
	return fieldRange{min, max, true}
}

var fieldDescriptors = []fieldDescriptor{
	FieldNanoOfSecond:            {name: "NanoOfSecond", javaName: "NANO_OF_SECOND", based: timeBased, fieldRange: makeRange(0, 999_999_999)},
	FieldNanoOfDay:               {name: "NanoOfDay", javaName: "NANO_OF_DAY", based: timeBased, fieldRange: makeRange(0, 86_399_999_999_999)},
	FieldMicroOfSecond:           {name: "MicroOfSecond", javaName: "MICRO_OF_SECOND", based: timeBased, fieldRange: makeRange(0, 999_999)},
	FieldMicroOfDay:              {name: "MicroOfDay", javaName: "MICRO_OF_DAY", based: timeBased, fieldRange: makeRange(0, 86_399_999_999)},
	FieldMilliOfSecond:           {name: "MilliOfSecond", javaName: "MILLI_OF_SECOND", based: timeBased, fieldRange: makeRange(0, 999)},
	FieldMilliOfDay:              {name: "MilliOfDay", javaName: "MILLI_OF_DAY", based: timeBased, fieldRange: makeRange(0, 86_399_999)},
	FieldSecondOfMinute:          {name: "SecondOfMinute", javaName: "SECOND_OF_MINUTE", based: timeBased, fieldRange: makeRange(0, 59)},
	FieldSecondOfDay:             {name: "SecondOfDay", javaName: "SECOND_OF_DAY", based: timeBased, fieldRange: makeRange(0, 86_399)},
	FieldMinuteOfHour:            {name: "MinuteOfHour", javaName: "MINUTE_OF_HOUR", based: timeBased, fieldRange: makeRange(0, 59)},
	FieldMinuteOfDay:             {name: "MinuteOfDay", javaName: "MINUTE_OF_DAY", based: timeBased, fieldRange: makeRange(0, 60*24-1)},
	FieldHourOfAmPm:              {name: "HourOfAmPm", javaName: "HOUR_OF_AMPM", based: timeBased, fieldRange: makeRange(0, 11)},
	FieldClockHourOfAmPm:         {name: "ClockHourOfAmPm", javaName: "CLOCK_HOUR_OF_AMPM", based: timeBased, fieldRange: makeRange(1, 12)},
	FieldHourOfDay:               {name: "HourOfDay", javaName: "HOUR_OF_DAY", based: timeBased, fieldRange: makeRange(0, 23)},
	FieldClockHourOfDay:          {name: "ClockHourOfDay", javaName: "CLOCK_HOUR_OF_DAY", based: timeBased, fieldRange: makeRange(1, 24)},
	FieldAmPmOfDay:               {name: "AmPmOfDay", javaName: "AMPM_OF_DAY", based: timeBased, fieldRange: makeRange(0, 1)},
	FieldDayOfWeek:               {name: "DayOfWeek", javaName: "DAY_OF_WEEK", based: dateBased, fieldRange: makeRange(1, 7)},
	FieldAlignedDayOfWeekInMonth: {name: "AlignedDayOfWeekInMonth", javaName: "ALIGNED_DAY_OF_WEEK_IN_MONTH", based: dateBased, fieldRange: makeRange(1, 7)},
	FieldAlignedDayOfWeekInYear:  {name: "AlignedDayOfWeekInYear", javaName: "ALIGNED_DAY_OF_WEEK_IN_YEAR", based: dateBased, fieldRange: makeRange(1, 7)},
	FieldDayOfMonth:              {name: "DayOfMonth", javaName: "DAY_OF_MONTH", based: dateBased, fieldRange: makeRange(1, 31)},
	FieldDayOfYear:               {name: "DayOfYear", javaName: "DAY_OF_YEAR", based: dateBased, fieldRange: makeRange(1, 366)},
	FieldEpochDay:                {name: "EpochDay", javaName: "EPOCH_DAY", based: dateBased, fieldRange: makeRange(math.MinInt64, math.MaxInt64)},
	FieldAlignedWeekOfMonth:      {name: "AlignedWeekOfMonth", javaName: "ALIGNED_WEEK_OF_MONTH", based: dateBased, fieldRange: makeRange(1, 5)},
	FieldAlignedWeekOfYear:       {name: "AlignedWeekOfYear", javaName: "ALIGNED_WEEK_OF_YEAR", based: dateBased, fieldRange: makeRange(1, 53)},
	FieldMonthOfYear:             {name: "MonthOfYear", javaName: "MONTH_OF_YEAR", based: dateBased, fieldRange: makeRange(1, 12)},
	FieldProlepticMonth:          {name: "ProlepticMonth", javaName: "PROLEPTIC_MONTH", based: dateBased, fieldRange: makeRange(YearMin*12, YearMax*12+11)},
	FieldYearOfEra:               {name: "YearOfEra", javaName: "YEAR_OF_ERA", based: dateBased, fieldRange: makeRange(1, math.MaxInt64)},
	FieldYear:                    {name: "Year", javaName: "YEAR", based: dateBased, fieldRange: makeRange(YearMin, YearMax)},
	FieldEra:                     {name: "Era", javaName: "ERA", based: dateBased, fieldRange: makeRange(0, 1)},
	FieldInstantSeconds:          {name: "InstantSeconds", javaName: "INSTANT_SECONDS", fieldRange: makeRange(math.MinInt64, math.MaxInt64)},
	FieldOffsetSeconds:           {name: "OffsetSeconds", javaName: "OFFSET_SECONDS", fieldRange: makeRange(-18*3600, 18*3600)},
}

func (f Field) check(value int64) error {
	if !f.Valid() {
		return invalidFieldError(f)
	}
	var r = fieldDescriptors[f].fieldRange
	if r.Valid && (value < r.Min || value > r.Max) {
		return fieldOutOfRangeError(f, value)
	}
	return nil
}

func (f Field) fieldRange() fieldRange {
	return fieldDescriptors[f].fieldRange
}

func (f Field) checkSetE(value int64, e *error) {
	if *e != nil {
		return
	}
	*e = f.check(value)
}

func (f Field) Valid() bool {
	return f > 0 && f <= FieldOffsetSeconds
}

// String returns the name of the field.
func (f Field) String() string {
	if f.Valid() {
		return fieldDescriptors[f].name
	}
	return "UnknownField(" + strconv.Itoa(int(f)) + ")"
}

// IsDateBased checks if this field represents a component of a date.
//
// A field is date-based if it can be derived from FieldEpochDay.
// Note that it is valid for both IsDateBased() and IsTimeBased() to return false,
// such as when representing a field like minute-of-week.
//
// Returns true if this field is a component of a date.
func (f Field) IsDateBased() bool {
	return f.Valid() && fieldDescriptors[f].based == dateBased
}

// IsTimeBased checks if this field represents a component of a time.
//
// A field is time-based if it can be derived from FieldNanoOfDay.
// Note that it is valid for both IsDateBased() and IsTimeBased() to return false,
// such as when representing a field like minute-of-week.
//
// Returns true if this field is a component of a time.
func (f Field) IsTimeBased() bool {
	return f.Valid() && fieldDescriptors[f].based == timeBased
}

func (f Field) JavaName() string {
	if !f.Valid() {
		return ""
	}
	var j = fieldDescriptors[f].javaName
	return "ChronoField." + j
}
