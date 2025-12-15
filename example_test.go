package goda_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/iseki0/goda"
)

// Example demonstrates basic usage of the goda package.
func Example() {
	// Create a specific date
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	fmt.Println("LocalDate:", date)

	// Create a specific time
	timeOfDay := goda.MustLocalTimeOf(14, 30, 45, 0)
	fmt.Println("LocalTime:", timeOfDay)

	// Combine date and time
	datetime := date.AtTime(timeOfDay)
	fmt.Println("LocalDateTime:", datetime)

	// With timezone offset
	offset := goda.MustZoneOffsetOfHours(8) // +08:00
	offsetDateTime := datetime.AtOffset(offset)
	fmt.Println("OffsetDateTime:", offsetDateTime)

	// Parse from string
	parsedDate := goda.MustLocalDateParse("2024-03-15")
	parsedTime := goda.MustLocalTimeParse("14:30:45.123456789")
	parsedDateTime := goda.MustLocalDateTimeParse("2024-03-15T14:30:45.123456789")
	parsedOffsetDateTime := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45+08:00")

	fmt.Println("Parsed LocalDate:", parsedDate)
	fmt.Println("Parsed LocalTime:", parsedTime)
	fmt.Println("Parsed LocalDateTime:", parsedDateTime)
	fmt.Println("Parsed OffsetDateTime:", parsedOffsetDateTime)

	// LocalDate arithmetic
	tomorrow := date.Chain().PlusDays(1).MustGet()
	nextMonth := date.Chain().PlusMonths(1).MustGet()
	nextYear := date.Chain().PlusYears(1).MustGet()

	fmt.Println("Tomorrow:", tomorrow)
	fmt.Println("Next month:", nextMonth)
	fmt.Println("Next year:", nextYear)

	// Comparisons
	date1 := goda.MustLocalDateOf(2024, goda.March, 15)
	date2 := goda.MustLocalDateOf(2024, goda.March, 20)
	if date1.IsBefore(date2) {
		fmt.Println("date1 is before date2")
	}

	// Serialization
	jsonBytes, _ := json.Marshal(date)
	fmt.Println("JSON:", string(jsonBytes))

	str := date.String()
	fmt.Println("String:", str)

	// Output:
	// LocalDate: 2024-03-15
	// LocalTime: 14:30:45
	// LocalDateTime: 2024-03-15T14:30:45
	// OffsetDateTime: 2024-03-15T14:30:45+08:00
	// Parsed LocalDate: 2024-03-15
	// Parsed LocalTime: 14:30:45.123456789
	// Parsed LocalDateTime: 2024-03-15T14:30:45.123456789
	// Parsed OffsetDateTime: 2024-03-15T14:30:45+08:00
	// Tomorrow: 2024-03-16
	// Next month: 2024-04-15
	// Next year: 2025-03-15
	// date1 is before date2
	// JSON: "2024-03-15"
	// String: 2024-03-15
}

// ExampleLocalDateNow demonstrates how to get the current date.
func ExampleLocalDateNow() {
	// Get current date in local timezone
	today := goda.LocalDateNow()

	// check that we got a valid date
	fmt.Printf("Got valid date: %v\n", !today.IsZero())
	fmt.Printf("Has year component: %v\n", today.Year() != 0)

	// Output:
	// Got valid date: true
	// Has year component: true
}

// ExampleLocalDateNowUTC demonstrates how to get the current date in UTC.
func ExampleLocalDateNowUTC() {
	// Get current date in UTC
	todayUTC := goda.LocalDateNowUTC()

	// Verify it's a valid date
	fmt.Printf("Valid: %v\n", !todayUTC.IsZero())

	// Output:
	// Valid: true
}

// ExampleLocalDateNowIn demonstrates how to get the current date in a specific timezone.
func ExampleLocalDateNowIn() {
	// Get current date in Tokyo timezone
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	todayTokyo := goda.LocalDateNowIn(tokyo)

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !todayTokyo.IsZero())

	// Output:
	// Valid: true
}

// ExampleLocalDateParse demonstrates parsing a date from a string.
func ExampleLocalDateParse() {
	// Parse a date string
	date, err := goda.LocalDateParse("2024-03-15")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleMustLocalDateParse demonstrates parsing a date that panics on error.
func ExampleMustLocalDateParse() {
	// Parse a date string (panics if invalid)
	date := goda.MustLocalDateParse("2024-03-15")
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleNewLocalDate demonstrates how to create a date.
func ExampleLocalDateOf() {
	// Create a valid date
	date, err := goda.LocalDateOf(2024, goda.January, 15)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(date)

	// Try to create an invalid date
	_, err = goda.LocalDateOf(2024, goda.February, 30)
	fmt.Println("Error:", err)

	// Output:
	// 2024-01-15
	// Error: goda: invalid date February 30
}

// ExampleMustNewLocalDate demonstrates how to create a date that panics on error.
func ExampleMustLocalDateOf() {
	// Create a date (panics if invalid)
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleLocalDateOfGoTime demonstrates converting from time.Time to LocalDate.
func ExampleLocalDateOfGoTime() {
	// Convert from time.LocalTime
	goTime := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
	date := goda.LocalDateOfGoTime(goTime)
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleLocalDateChain_PlusDays demonstrates adding days to a date.
func ExampleLocalDateChain_PlusDays() {
	date := goda.MustLocalDateOf(2024, goda.January, 15)
	fmt.Println("Original:", date)
	fmt.Println("Plus 10 days:", date.Chain().PlusDays(10).MustGet())
	fmt.Println("Minus 10 days:", date.Chain().PlusDays(-10).MustGet())

	// Output:
	// Original: 2024-01-15
	// Plus 10 days: 2024-01-25
	// Minus 10 days: 2024-01-05
}

// ExampleLocalDateChain_PlusMonths demonstrates adding months to a date.
func ExampleLocalDateChain_PlusMonths() {
	date := goda.MustLocalDateOf(2024, goda.January, 31)
	fmt.Println("Original:", date)
	fmt.Println("Plus 1 month:", date.Chain().PlusMonths(1).MustGet())
	fmt.Println("Plus 2 months:", date.Chain().PlusMonths(2).MustGet())

	// Output:
	// Original: 2024-01-31
	// Plus 1 month: 2024-02-29
	// Plus 2 months: 2024-03-31
}

// ExampleLocalDate_Compare demonstrates comparing dates.
func ExampleLocalDate_Compare() {
	date1 := goda.MustLocalDateOf(2024, goda.March, 15)
	date2 := goda.MustLocalDateOf(2024, goda.March, 20)
	date3 := goda.MustLocalDateOf(2024, goda.March, 15)

	fmt.Println("date1 < date2:", date1.IsBefore(date2))
	fmt.Println("date1 > date2:", date1.IsAfter(date2))
	fmt.Println("date1 == date3:", date1.Compare(date3) == 0)

	// Output:
	// date1 < date2: true
	// date1 > date2: false
	// date1 == date3: true
}

// ExampleLocalDate_DayOfWeek demonstrates getting the day of week.
func ExampleLocalDate_DayOfWeek() {
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	fmt.Println("Day of week:", date.DayOfWeek())
	fmt.Println("Is Friday?", date.DayOfWeek() == goda.Friday)

	// Output:
	// Day of week: Friday
	// Is Friday? true
}

// ExampleLocalTimeNow demonstrates how to get the current time.
func ExampleLocalTimeNow() {
	// Get current time in local timezone
	now := goda.LocalTimeNow()

	// Verify it's valid and within range
	fmt.Printf("Valid: %v\n", !now.IsZero())
	fmt.Printf("Hour in range: %v\n", now.Hour() >= 0 && now.Hour() < 24)

	// Output:
	// Valid: true
	// Hour in range: true
}

// ExampleLocalTimeNowUTC demonstrates how to get the current time in UTC.
func ExampleLocalTimeNowUTC() {
	// Get current time in UTC
	nowUTC := goda.LocalTimeNowUTC()

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !nowUTC.IsZero())

	// Output:
	// Valid: true
}

// ExampleLocalTimeNowIn demonstrates how to get the current time in a specific timezone.
func ExampleLocalTimeNowIn() {
	// Get current time in Tokyo timezone
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	nowTokyo := goda.LocalTimeNowIn(tokyo)

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !nowTokyo.IsZero())

	// Output:
	// Valid: true
}

// ExampleLocalTimeParse demonstrates parsing a time from a string.
func ExampleLocalTimeParse() {
	// Parse a time string
	t, err := goda.LocalTimeParse("14:30:45.123456789")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)

	// Output:
	// 14:30:45.123456789
}

// ExampleMustLocalTimeParse demonstrates parsing a time that panics on error.
func ExampleMustLocalTimeParse() {
	// Parse a time string (panics if invalid)
	t := goda.MustLocalTimeParse("14:30:45")
	fmt.Println(t)

	// Output:
	// 14:30:45
}

// ExampleLocalTimeOf demonstrates how to create a time.
func ExampleLocalTimeOf() {
	// Create a valid time
	t, err := goda.LocalTimeOf(14, 30, 45, 123456789)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)

	// Try to create an invalid time
	_, err = goda.LocalTimeOf(25, 0, 0, 0)
	fmt.Println("Error:", err)

	// Output:
	// 14:30:45.123456789
	// Error: goda: invalid value of HourOfDay (valid range 0 - 23): 25
}

// ExampleMustLocalTimeOf demonstrates how to create a time that panics on error.
func ExampleMustLocalTimeOf() {
	// Create a time (panics if invalid)
	t := goda.MustLocalTimeOf(14, 30, 45, 0)
	fmt.Println(t)

	// Midnight
	midnight := goda.MustLocalTimeOf(0, 0, 0, 0)
	fmt.Println("Midnight:", midnight)

	// Output:
	// 14:30:45
	// Midnight: 00:00:00
}

// ExampleLocalTimeOfGoTime demonstrates converting from time.Time to LocalTime.
func ExampleLocalTimeOfGoTime() {
	// Convert from time.LocalTime
	goTime := time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
	t := goda.LocalTimeOfGoTime(goTime)
	fmt.Println(t)

	// Output:
	// 14:30:45.123456789
}

// ExampleLocalTime_hour demonstrates accessing time components.
func ExampleLocalTime_hour() {
	t := goda.MustLocalTimeOf(14, 30, 45, 123456789)
	fmt.Println("Hour:", t.Hour())
	fmt.Println("Minute:", t.Minute())
	fmt.Println("Second:", t.Second())
	fmt.Println("Millisecond:", t.Millisecond())
	fmt.Println("Nano:", t.Nano())

	// Output:
	// Hour: 14
	// Minute: 30
	// Second: 45
	// Millisecond: 123
	// Nano: 123456789
}

// ExampleLocalTime_Compare demonstrates comparing times.
func ExampleLocalTime_Compare() {
	t1 := goda.MustLocalTimeOf(14, 30, 0, 0)
	t2 := goda.MustLocalTimeOf(15, 0, 0, 0)
	t3 := goda.MustLocalTimeOf(14, 30, 0, 0)

	fmt.Println("t1 < t2:", t1.IsBefore(t2))
	fmt.Println("t1 > t2:", t1.IsAfter(t2))
	fmt.Println("t1 == t3:", t1.Compare(t3) == 0)

	// Output:
	// t1 < t2: true
	// t1 > t2: false
	// t1 == t3: true
}

// ExampleLocalTime_String demonstrates the string format with fractional seconds.
func ExampleLocalTime_String() {
	// LocalTime without fractional seconds
	t1 := goda.MustLocalTimeOf(14, 30, 45, 0)
	fmt.Println(t1)

	// LocalTime with milliseconds
	t2 := goda.MustLocalTimeOf(14, 30, 45, 123000000)
	fmt.Println(t2)

	// LocalTime with microseconds
	t3 := goda.MustLocalTimeOf(14, 30, 45, 123456000)
	fmt.Println(t3)

	// LocalTime with nanoseconds
	t4 := goda.MustLocalTimeOf(14, 30, 45, 123456789)
	fmt.Println(t4)

	// LocalTime with trailing zeros aligned to 3-digit boundaries
	t5 := goda.MustLocalTimeOf(14, 30, 45, 100000000)
	fmt.Println(t5)

	// Output:
	// 14:30:45
	// 14:30:45.123
	// 14:30:45.123456
	// 14:30:45.123456789
	// 14:30:45.100
}

// ExampleMonth demonstrates working with months.
func ExampleMonth() {
	// Months are constants from January (1) to December (12)
	fmt.Println("January:", goda.January)
	fmt.Println("December:", goda.December)

	// Get month length
	fmt.Println("Days in February (non-leap):", goda.February.Length(false))
	fmt.Println("Days in February (leap):", goda.February.Length(true))

	// Output:
	// January: January
	// December: December
	// Days in February (non-leap): 28
	// Days in February (leap): 29
}

// ExampleYear_IsLeapYear demonstrates checking for leap years.
func ExampleYear_IsLeapYear() {
	fmt.Println("2024 is leap:", goda.Year(2024).IsLeapYear())
	fmt.Println("2023 is leap:", goda.Year(2023).IsLeapYear())
	fmt.Println("2000 is leap:", goda.Year(2000).IsLeapYear())
	fmt.Println("1900 is leap:", goda.Year(1900).IsLeapYear())

	// Output:
	// 2024 is leap: true
	// 2023 is leap: false
	// 2000 is leap: true
	// 1900 is leap: false
}

// ExampleLocalDate_MarshalJSON demonstrates JSON serialization.
func ExampleLocalDate_MarshalJSON() {
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	jsonBytes, _ := json.Marshal(date)
	fmt.Println(string(jsonBytes))

	// Output:
	// "2024-03-15"
}

// ExampleLocalDate_UnmarshalJSON demonstrates JSON deserialization.
func ExampleLocalDate_UnmarshalJSON() {
	var date goda.LocalDate
	jsonData := []byte(`"2024-03-15"`)
	err := json.Unmarshal(jsonData, &date)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleLocalTime_MarshalJSON demonstrates JSON serialization.
func ExampleLocalTime_MarshalJSON() {
	t := goda.MustLocalTimeOf(14, 30, 45, 123456789)
	jsonBytes, _ := json.Marshal(t)
	fmt.Println(string(jsonBytes))

	// Output:
	// "14:30:45.123456789"
}

// ExampleLocalTime_UnmarshalJSON demonstrates JSON deserialization.
func ExampleLocalTime_UnmarshalJSON() {
	var t goda.LocalTime
	jsonData := []byte(`"14:30:45.123456789"`)
	err := json.Unmarshal(jsonData, &t)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)

	// Output:
	// 14:30:45.123456789
}

// ExampleLocalDateTime demonstrates basic LocalDateTime usage.
func ExampleLocalDateTime() {
	// Create from components
	dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 123456789)
	fmt.Println(dt)

	// Access date and time parts
	fmt.Printf("LocalDate: %s\n", dt.LocalDate())
	fmt.Printf("LocalTime: %s\n", dt.LocalTime())

	// Access individual components
	fmt.Printf("Year: %d, Hour: %d\n", dt.Year(), dt.Hour())

	// Output:
	// 2024-03-15T14:30:45.123456789
	// LocalDate: 2024-03-15
	// LocalTime: 14:30:45.123456789
	// Year: 2024, Hour: 14
}

// ExampleLocalDateTimeParse demonstrates parsing a datetime from a string.
func ExampleLocalDateTimeParse() {
	dt, err := goda.LocalDateTimeParse("2024-03-15T14:30:45.123456789")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleMustLocalDateTimeParse demonstrates parsing that panics on error.
func ExampleMustLocalDateTimeParse() {
	dt := goda.MustLocalDateTimeParse("2024-03-15T14:30:45")
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45
}

// ExampleLocalDate_AtTime demonstrates creating a datetime from date and time.
func ExampleLocalDate_AtTime() {
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	time := goda.MustLocalTimeOf(14, 30, 45, 123456789)
	dt := date.AtTime(time)
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleLocalTime_AtDate demonstrates creating a datetime from time and date.
func ExampleLocalTime_AtDate() {
	time := goda.MustLocalTimeOf(14, 30, 45, 123456789)
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	dt := time.AtDate(date)
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleLocalDateTimeNow demonstrates getting the current datetime.
func ExampleLocalDateTimeNow() {
	now := goda.LocalDateTimeNow()

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !now.IsZero())
	fmt.Printf("Has components: %v\n", now.Year() != 0 && now.Hour() >= 0)

	// Output:
	// Valid: true
	// Has components: true
}

// ExampleLocalDateTime_Compare demonstrates comparing datetimes.
func ExampleLocalDateTime_Compare() {
	dt1 := goda.MustLocalDateTimeParse("2024-03-15T14:30:45")
	dt2 := goda.MustLocalDateTimeParse("2024-03-15T14:30:45")
	dt3 := goda.MustLocalDateTimeParse("2024-03-15T15:30:45")

	fmt.Printf("dt1 == dt2: %v\n", dt1.Compare(dt2) == 0)
	fmt.Printf("dt1 < dt3: %v\n", dt1.IsBefore(dt3))
	fmt.Printf("dt3 > dt1: %v\n", dt3.IsAfter(dt1))

	// Output:
	// dt1 == dt2: true
	// dt1 < dt3: true
	// dt3 > dt1: true
}

// ExampleLocalDateTimeChain_PlusDays demonstrates adding days.
func ExampleLocalDateTimeChain_PlusDays() {
	dt := goda.MustLocalDateTimeParse("2024-03-15T14:30:45")
	future := dt.Chain().PlusDays(10).MustGet()
	fmt.Println(future)

	// Output:
	// 2024-03-25T14:30:45
}

// ExampleLocalDateTime_MarshalJSON demonstrates JSON serialization.
func ExampleLocalDateTime_MarshalJSON() {
	dt := goda.MustLocalDateTimeParse("2024-03-15T14:30:45.123456789")
	jsonBytes, _ := json.Marshal(dt)
	fmt.Println(string(jsonBytes))

	// Output:
	// "2024-03-15T14:30:45.123456789"
}

// ExampleLocalDateTime_UnmarshalJSON demonstrates JSON deserialization.
func ExampleLocalDateTime_UnmarshalJSON() {
	var dt goda.LocalDateTime
	jsonData := []byte(`"2024-03-15T14:30:45.123456789"`)
	err := json.Unmarshal(jsonData, &dt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleLocalDateTime_IsSupportedField demonstrates checking field support.
func ExampleLocalDateTime_IsSupportedField() {
	dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 0)

	fmt.Printf("Supports FieldHourOfDay: %v\n", dt.IsSupportedField(goda.FieldHourOfDay))
	fmt.Printf("Supports FieldDayOfMonth: %v\n", dt.IsSupportedField(goda.FieldDayOfMonth))
	fmt.Printf("Supports FieldOffsetSeconds: %v\n", dt.IsSupportedField(goda.FieldOffsetSeconds))

	// Output:
	// Supports FieldHourOfDay: true
	// Supports FieldDayOfMonth: true
	// Supports FieldOffsetSeconds: false
}

// ExampleLocalDate_GetField demonstrates querying date fields with TemporalValue.
func ExampleLocalDate_GetField() {
	date := goda.MustLocalDateOf(2024, goda.March, 15) // Friday

	// Query various date fields
	dayOfWeek := date.GetField(goda.FieldDayOfWeek)
	if dayOfWeek.Valid() {
		fmt.Printf("Day of week: %d (1=Monday, 7=Sunday)\n", dayOfWeek.Int())
	}

	month := date.GetField(goda.FieldMonthOfYear)
	if month.Valid() {
		fmt.Printf("Month: %d\n", month.Int())
	}

	year := date.GetField(goda.FieldYear)
	if year.Valid() {
		fmt.Printf("Year: %d\n", year.Int())
	}

	// Query unsupported field (time field on date)
	hour := date.GetField(goda.FieldHourOfDay)
	if hour.Unsupported() {
		fmt.Println("Hour field is not supported for LocalDate")
	}

	// Output:
	// Day of week: 5 (1=Monday, 7=Sunday)
	// Month: 3
	// Year: 2024
	// Hour field is not supported for LocalDate
}

// ExampleLocalDate_GetField_advancedFields demonstrates advanced date field queries.
func ExampleLocalDate_GetField_advancedFields() {
	date := goda.MustLocalDateOf(2024, goda.March, 15)

	// Day of year (1-366)
	dayOfYear := date.GetField(goda.FieldDayOfYear)
	fmt.Printf("Day of year: %d\n", dayOfYear.Int())

	// Epoch days (days since 1970-01-01)
	epochDay := date.GetField(goda.FieldEpochDay)
	fmt.Printf("Days since Unix epoch: %d\n", epochDay.Int64())

	// Proleptic month (months since year 0)
	prolepticMonth := date.GetField(goda.FieldProlepticMonth)
	fmt.Printf("Proleptic month: %d\n", prolepticMonth.Int64())

	// FieldEra (0=BCE, 1=CE)
	era := date.GetField(goda.FieldEra)
	if era.Int() == 1 {
		fmt.Println("FieldEra: CE (Common FieldEra)")
	}

	// Output:
	// Day of year: 75
	// Days since Unix epoch: 19797
	// Proleptic month: 24290
	// FieldEra: CE (Common FieldEra)
}

// ExampleLocalTime_GetField demonstrates querying time fields with TemporalValue.
func ExampleLocalTime_GetField() {
	t := goda.MustLocalTimeOf(14, 30, 45, 123456789)

	// Query various time fields
	hour := t.GetField(goda.FieldHourOfDay)
	if hour.Valid() {
		fmt.Printf("Hour of day (0-23): %d\n", hour.Int())
	}

	minute := t.GetField(goda.FieldMinuteOfHour)
	if minute.Valid() {
		fmt.Printf("Minute of hour: %d\n", minute.Int())
	}

	second := t.GetField(goda.FieldSecondOfMinute)
	if second.Valid() {
		fmt.Printf("Second of minute: %d\n", second.Int())
	}

	nanos := t.GetField(goda.FieldNanoOfSecond)
	if nanos.Valid() {
		fmt.Printf("Nanoseconds: %d\n", nanos.Int())
	}

	// Query unsupported field (date field on time)
	dayOfMonth := t.GetField(goda.FieldDayOfMonth)
	if dayOfMonth.Unsupported() {
		fmt.Println("FieldDayOfMonth field is not supported for LocalTime")
	}

	// Output:
	// Hour of day (0-23): 14
	// Minute of hour: 30
	// Second of minute: 45
	// Nanoseconds: 123456789
	// FieldDayOfMonth field is not supported for LocalTime
}

// ExampleLocalTime_GetField_clockHours demonstrates 12-hour clock field queries.
func ExampleLocalTime_GetField_clockHours() {
	// Afternoon time (2:30 PM)
	afternoon := goda.MustLocalTimeOf(14, 30, 0, 0)

	// 24-hour format
	hourOfDay := afternoon.GetField(goda.FieldHourOfDay)
	fmt.Printf("24-hour format: %d:30\n", hourOfDay.Int())

	// 12-hour format components
	hourOfAmPm := afternoon.GetField(goda.FieldHourOfAmPm)
	amPm := afternoon.GetField(goda.FieldAmPmOfDay)
	amPmStr := "AM"
	if amPm.Int() == 1 {
		amPmStr = "PM"
	}
	fmt.Printf("12-hour format: %d:30 %s\n", hourOfAmPm.Int(), amPmStr)

	// Clock hour (1-12 instead of 0-11)
	clockHour := afternoon.GetField(goda.FieldClockHourOfAmPm)
	fmt.Printf("Clock hour (1-12): %d:30 %s\n", clockHour.Int(), amPmStr)

	// Midnight special case
	midnight := goda.MustLocalTimeOf(0, 0, 0, 0)
	midnightClock := midnight.GetField(goda.FieldClockHourOfDay)
	fmt.Printf("Midnight clock hour: %d:00\n", midnightClock.Int())

	// Output:
	// 24-hour format: 14:30
	// 12-hour format: 2:30 PM
	// Clock hour (1-12): 2:30 PM
	// Midnight clock hour: 24:00
}

// ExampleLocalTime_GetField_ofDayFields demonstrates querying cumulative daily values.
func ExampleLocalTime_GetField_ofDayFields() {
	t := goda.MustLocalTimeOf(14, 30, 45, 500000000) // 2:30:45.5 PM

	// Total seconds elapsed since midnight
	secondOfDay := t.GetField(goda.FieldSecondOfDay)
	fmt.Printf("Seconds since midnight: %d\n", secondOfDay.Int())

	// Total minutes elapsed since midnight
	minuteOfDay := t.GetField(goda.FieldMinuteOfDay)
	fmt.Printf("Minutes since midnight: %d\n", minuteOfDay.Int())

	// Total milliseconds elapsed since midnight
	milliOfDay := t.GetField(goda.FieldMilliOfDay)
	fmt.Printf("Milliseconds since midnight: %d\n", milliOfDay.Int64())

	// Total nanoseconds elapsed since midnight
	nanoOfDay := t.GetField(goda.FieldNanoOfDay)
	fmt.Printf("Nanoseconds since midnight: %d\n", nanoOfDay.Int64())

	// Output:
	// Seconds since midnight: 52245
	// Minutes since midnight: 870
	// Milliseconds since midnight: 52245500
	// Nanoseconds since midnight: 52245500000000
}

// ExampleLocalDateTime_GetField demonstrates querying fields from a date-time.
func ExampleLocalDateTime_GetField() {
	dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 123456789)

	// Query date fields
	year := dt.GetField(goda.FieldYear)
	month := dt.GetField(goda.FieldMonthOfYear)
	day := dt.GetField(goda.FieldDayOfMonth)

	if year.Valid() && month.Valid() && day.Valid() {
		fmt.Printf("Date: %04d-%02d-%02d\n", year.Int(), month.Int(), day.Int())
	}

	// Query time fields
	hour := dt.GetField(goda.FieldHourOfDay)
	minute := dt.GetField(goda.FieldMinuteOfHour)
	second := dt.GetField(goda.FieldSecondOfMinute)

	if hour.Valid() && minute.Valid() && second.Valid() {
		fmt.Printf("Time: %02d:%02d:%02d\n", hour.Int(), minute.Int(), second.Int())
	}

	// Query day of week
	dayOfWeek := dt.GetField(goda.FieldDayOfWeek)
	if dayOfWeek.Valid() {
		fmt.Printf("Day of week: %d (Friday)\n", dayOfWeek.Int())
	}

	// Output:
	// Date: 2024-03-15
	// Time: 14:30:45
	// Day of week: 5 (Friday)
}

// ExampleLocalDateTime_GetField_delegation demonstrates field delegation.
func ExampleLocalDateTime_GetField_delegation() {
	dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 0)

	// LocalDateTime delegates date fields to LocalDate
	dayOfYear := dt.GetField(goda.FieldDayOfYear)
	fmt.Printf("Day of year: %d\n", dayOfYear.Int())

	// LocalDateTime delegates time fields to LocalTime
	nanoOfDay := dt.GetField(goda.FieldNanoOfDay)
	fmt.Printf("Nanoseconds of day: %d\n", nanoOfDay.Int64())

	// Unsupported fields return unsupported TemporalValue
	offsetSeconds := dt.GetField(goda.FieldOffsetSeconds)
	if offsetSeconds.Unsupported() {
		fmt.Println("FieldOffsetSeconds is not supported for LocalDateTime")
	}

	// Output:
	// Day of year: 75
	// Nanoseconds of day: 52245000000000
	// FieldOffsetSeconds is not supported for LocalDateTime
}

// ExampleZoneOffset demonstrates basic ZoneOffset usage.
func ExampleZoneOffset() {
	// Create zone offsets
	utc := goda.ZoneOffsetUTC()
	fmt.Println("UTC:", utc)

	// From hours
	est := goda.MustZoneOffsetOfHours(-5)
	fmt.Println("EST:", est)

	// From hours and minutes
	ist := goda.MustZoneOffsetOfHoursMinutes(5, 30)
	fmt.Println("IST:", ist)

	// From total seconds
	jst := goda.MustZoneOffsetOfSeconds(9 * 3600)
	fmt.Println("JST:", jst)

	// Output:
	// UTC: Z
	// EST: -05:00
	// IST: +05:30
	// JST: +09:00
}

// ExampleZoneOffsetParse demonstrates parsing zone offsets.
func ExampleZoneOffsetParse() {
	// Parse various formats
	z1 := goda.MustZoneOffsetParse("Z")
	fmt.Println("Z:", z1.TotalSeconds())

	z2 := goda.MustZoneOffsetParse("+02:00")
	fmt.Println("+02:00:", z2.TotalSeconds())

	z3 := goda.MustZoneOffsetParse("-05:30")
	fmt.Println("-05:30:", z3.TotalSeconds())

	z4 := goda.MustZoneOffsetParse("+0930")
	fmt.Println("+0930:", z4.TotalSeconds())

	// Output:
	// Z: 0
	// +02:00: 7200
	// -05:30: -19800
	// +0930: 34200
}

// ExampleZoneOffset_TotalSeconds demonstrates accessing offset components.
func ExampleZoneOffset_TotalSeconds() {
	offset := goda.MustZoneOffsetParse("+05:30")

	fmt.Println("Total seconds:", offset.TotalSeconds())
	fmt.Println("Hours:", offset.Hours())
	fmt.Println("Minutes:", offset.Minutes())
	fmt.Println("Seconds:", offset.Seconds())

	// Output:
	// Total seconds: 19800
	// Hours: 5
	// Minutes: 30
	// Seconds: 0
}

// ExampleZoneOffset_Compare demonstrates comparing zone offsets.
func ExampleZoneOffset_Compare() {
	utc := goda.ZoneOffsetUTC()
	est := goda.MustZoneOffsetOfHours(-5)
	cet := goda.MustZoneOffsetOfHours(1)

	fmt.Printf("EST < UTC: %v\n", est.Compare(utc) < 0)
	fmt.Printf("CET > UTC: %v\n", cet.Compare(utc) > 0)
	fmt.Printf("EST < CET: %v\n", est.Compare(cet) < 0)

	// Output:
	// EST < UTC: true
	// CET > UTC: true
	// EST < CET: true
}

// ExampleZoneOffset_MarshalJSON demonstrates JSON serialization.
func ExampleZoneOffset_MarshalJSON() {
	offset := goda.MustZoneOffsetParse("+05:30")
	jsonBytes, _ := json.Marshal(offset)
	fmt.Println(string(jsonBytes))

	// Output:
	// "+05:30"
}

// ExampleOffsetDateTime demonstrates basic OffsetDateTime usage.
func ExampleOffsetDateTime() {
	// Create from components
	offset := goda.MustZoneOffsetOfHours(8) // +08:00 (China Standard Time)
	odt := goda.MustOffsetDateTimeOf(2024, goda.March, 15, 14, 30, 45, 0, offset)
	fmt.Println(odt)

	// Create from LocalDateTime
	dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 0)
	odt2 := dt.AtOffset(offset)
	fmt.Println(odt2)

	// Parse from string
	odt3 := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45+08:00")
	fmt.Println(odt3)

	// UTC
	utcOdt := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45Z")
	fmt.Println(utcOdt)

	// Output:
	// 2024-03-15T14:30:45+08:00
	// 2024-03-15T14:30:45+08:00
	// 2024-03-15T14:30:45+08:00
	// 2024-03-15T14:30:45Z
}

// ExampleOffsetDateTimeParse demonstrates parsing offset date-times.
func ExampleOffsetDateTimeParse() {
	// With positive offset
	odt1, _ := goda.OffsetDateTimeParse("2024-03-15T14:30:45+08:00")
	fmt.Println(odt1)

	// With negative offset
	odt2, _ := goda.OffsetDateTimeParse("2024-03-15T14:30:45-05:00")
	fmt.Println(odt2)

	// UTC (Z)
	odt3, _ := goda.OffsetDateTimeParse("2024-03-15T14:30:45Z")
	fmt.Println(odt3)

	// With minutes offset
	odt4, _ := goda.OffsetDateTimeParse("2024-03-15T14:30:45+05:30")
	fmt.Println(odt4)

	// Output:
	// 2024-03-15T14:30:45+08:00
	// 2024-03-15T14:30:45-05:00
	// 2024-03-15T14:30:45Z
	// 2024-03-15T14:30:45+05:30
}

// ExampleOffsetDateTimeNow demonstrates getting the current time with offset.
func ExampleOffsetDateTimeNow() {
	// Get current time with system's local offset
	now := goda.OffsetDateTimeNow()
	fmt.Printf("Has current time: %v\n", !now.IsZero())

	// Get current UTC time
	utcNow := goda.OffsetDateTimeNowUTC()
	fmt.Printf("UTC offset: %d seconds\n", utcNow.Offset().TotalSeconds())

	// Output:
	// Has current time: true
	// UTC offset: 0 seconds
}

// ExampleOffsetDateTime_Compare demonstrates comparing offset date-times.
func ExampleOffsetDateTime_Compare() {
	// These represent the same instant in time
	odt1 := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45+08:00")
	odt2 := goda.MustOffsetDateTimeParse("2024-03-15T06:30:45Z")

	if odt1.Compare(odt2) == 0 {
		fmt.Println("Same instant")
	}

	// Different instants
	odt3 := goda.MustOffsetDateTimeParse("2024-03-15T15:30:45+08:00")
	if odt3.IsAfter(odt1) {
		fmt.Println("odt3 is later")
	}

	// Output:
	// odt3 is later
}

// ExampleOffsetDateTimeChain_PlusHours demonstrates time arithmetic with hours.
func ExampleOffsetDateTimeChain_PlusHours() {
	odt := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45+08:00")

	// Add hours
	later := odt.Chain().PlusHours(5).MustGet()
	fmt.Println("5 hours later:", later)

	// Subtract hours
	earlier := odt.Chain().MinusHours(2).MustGet()
	fmt.Println("2 hours earlier:", earlier)

	// Output:
	// 5 hours later: 2024-03-15T19:30:45+08:00
	// 2 hours earlier: 2024-03-15T12:30:45+08:00
}

// ExampleOffsetDateTime_EpochSecond demonstrates Unix timestamp conversion.
func ExampleOffsetDateTime_EpochSecond() {
	// Different offsets, same instant
	odt1 := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45+08:00")
	odt2 := goda.MustOffsetDateTimeParse("2024-03-15T06:30:45Z")

	epoch1 := odt1.EpochSecond()
	epoch2 := odt2.EpochSecond()

	fmt.Printf("Has epoch seconds: %v\n", epoch1 > 0)
	fmt.Printf("Same instant: %v\n", epoch1 == epoch2)

	// Output:
	// Has epoch seconds: true
	// Same instant: true
}

// ExampleOffsetDateTime_MarshalJSON demonstrates JSON serialization.
func ExampleOffsetDateTime_MarshalJSON() {
	odt := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45+08:00")
	jsonBytes, _ := json.Marshal(odt)
	fmt.Println(string(jsonBytes))

	// UTC
	utcOdt := goda.MustOffsetDateTimeParse("2024-03-15T14:30:45Z")
	utcBytes, _ := json.Marshal(utcOdt)
	fmt.Println(string(utcBytes))

	// Output:
	// "2024-03-15T14:30:45+08:00"
	// "2024-03-15T14:30:45Z"
}

// ExampleOffsetDateTime_UnmarshalJSON demonstrates JSON deserialization.
func ExampleOffsetDateTime_UnmarshalJSON() {
	jsonData := []byte(`"2024-03-15T14:30:45+08:00"`)

	var odt goda.OffsetDateTime
	_ = json.Unmarshal(jsonData, &odt)

	fmt.Printf("Year: %d\n", odt.Year())
	fmt.Printf("Offset: %d hours\n", odt.Offset().TotalSeconds()/3600)

	// Output:
	// Year: 2024
	// Offset: 8 hours
}

// ExampleZoneId demonstrates basic ZoneId usage.
func ExampleZoneId() {
	// Create zone IDs for different time zones
	newYork, _ := goda.ZoneIdOf("America/New_York")
	fmt.Println("New York:", newYork)

	tokyo := goda.MustZoneIdOf("Asia/Tokyo")
	fmt.Println("Tokyo:", tokyo)

	// UTC is a common zone
	utc := goda.ZoneIdUTC()
	fmt.Println("UTC:", utc)

	// Get system's default zone
	defaultZone := goda.ZoneIdDefault()
	fmt.Printf("Has default zone: %v\n", !defaultZone.IsZero())

	// Output:
	// New York: America/New_York
	// Tokyo: Asia/Tokyo
	// UTC: UTC
	// Has default zone: true
}

// ExampleZoneIdOf demonstrates creating a ZoneId from a string.
func ExampleZoneIdOf() {
	// Valid zone ID
	zone, err := goda.ZoneIdOf("Europe/London")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Zone:", zone)

	// Invalid zone ID returns error
	_, err = goda.ZoneIdOf("Invalid/Zone")
	if err != nil {
		fmt.Println("Got expected error for invalid zone")
	}

	// Output:
	// Zone: Europe/London
	// Got expected error for invalid zone
}

// ExampleMustZoneIdOf demonstrates creating a ZoneId that panics on error.
func ExampleMustZoneIdOf() {
	// Valid zone ID
	zone := goda.MustZoneIdOf("Asia/Shanghai")
	fmt.Println(zone)

	// Output:
	// Asia/Shanghai
}

// ExampleZoneIdUTC demonstrates getting the UTC zone.
func ExampleZoneIdUTC() {
	utc := goda.ZoneIdUTC()
	fmt.Println("Zone:", utc)
	fmt.Println("Is UTC:", utc.String() == "UTC")

	// Output:
	// Zone: UTC
	// Is UTC: true
}

// ExampleZoneIdOfGoLocation demonstrates creating a ZoneId from time.Location.
func ExampleZoneIdOfGoLocation() {
	// From Go's time.UTC
	utcZone := goda.ZoneIdOfGoLocation(time.UTC)
	fmt.Println("UTC:", utcZone)

	// From Go's time.Local
	localZone := goda.ZoneIdOfGoLocation(time.Local)
	fmt.Printf("Has local zone: %v\n", !localZone.IsZero())

	// From loaded location
	paris, _ := time.LoadLocation("Europe/Paris")
	parisZone := goda.ZoneIdOfGoLocation(paris)
	fmt.Println("Paris:", parisZone)

	// Output:
	// UTC: UTC
	// Has local zone: true
	// Paris: Europe/Paris
}

// ExampleZoneIdDefault demonstrates getting the system's default zone.
func ExampleZoneIdDefault() {
	defaultZone := goda.ZoneIdDefault()
	fmt.Printf("Has default zone: %v\n", !defaultZone.IsZero())
	fmt.Printf("Zone name not empty: %v\n", defaultZone.String() != "")

	// Output:
	// Has default zone: true
	// Zone name not empty: true
}

// ExampleZoneId_IsZero demonstrates checking for zero value.
func ExampleZoneId_IsZero() {
	// Zero value
	var zero goda.ZoneId
	fmt.Println("Zero is zero:", zero.IsZero())

	// Non-zero value
	zone := goda.ZoneIdUTC()
	fmt.Println("UTC is zero:", zone.IsZero())

	// Output:
	// Zero is zero: true
	// UTC is zero: false
}

// ExampleZoneId_String demonstrates getting the string representation.
func ExampleZoneId_String() {
	// Various zones
	zones := []string{
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
	}

	for _, name := range zones {
		zone := goda.MustZoneIdOf(name)
		fmt.Println(zone.String())
	}

	// Output:
	// America/New_York
	// Europe/London
	// Asia/Tokyo
	// Australia/Sydney
}

// ExampleZoneId_MarshalJSON demonstrates JSON serialization.
func ExampleZoneId_MarshalJSON() {
	zone := goda.MustZoneIdOf("Asia/Singapore")
	jsonBytes, _ := json.Marshal(zone)
	fmt.Println(string(jsonBytes))

	// Zero value
	var zero goda.ZoneId
	zeroBytes, _ := json.Marshal(zero)
	fmt.Println(string(zeroBytes))

	// Output:
	// "Asia/Singapore"
	// ""
}

// ExampleZoneId_UnmarshalJSON demonstrates JSON deserialization.
func ExampleZoneId_UnmarshalJSON() {
	// From JSON string
	var zone goda.ZoneId
	err := json.Unmarshal([]byte(`"Europe/Berlin"`), &zone)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Zone:", zone)

	// From JSON null
	var nullZone goda.ZoneId
	_ = json.Unmarshal([]byte(`null`), &nullZone)
	fmt.Println("Null is zero:", nullZone.IsZero())

	// Output:
	// Zone: Europe/Berlin
	// Null is zero: true
}

// ExampleZoneId_roundTrip demonstrates serialization round-trip.
func ExampleZoneId_roundTrip() {
	// Original zone
	original := goda.MustZoneIdOf("Pacific/Auckland")

	// Serialize to JSON
	jsonData, _ := json.Marshal(original)

	// Deserialize from JSON
	var restored goda.ZoneId
	_ = json.Unmarshal(jsonData, &restored)

	// Compare
	fmt.Println("Original:", original)
	fmt.Println("Restored:", restored)
	fmt.Println("Same:", original.String() == restored.String())

	// Output:
	// Original: Pacific/Auckland
	// Restored: Pacific/Auckland
	// Same: true
}

// ExampleLocalDate_PlusWeeks demonstrates adding weeks to a date.
func ExampleLocalDateChain_PlusWeeks() {
	date := goda.MustLocalDateOf(2024, goda.January, 15)
	fmt.Println("Original:", date)
	fmt.Println("Plus 2 weeks:", date.Chain().PlusWeeks(2).MustGet())
	fmt.Println("Minus 1 week:", date.Chain().PlusWeeks(-1).MustGet())

	// Output:
	// Original: 2024-01-15
	// Plus 2 weeks: 2024-01-29
	// Minus 1 week: 2024-01-08
}

// ExampleLocalDateChain_MinusWeeks demonstrates subtracting weeks from a date.
func ExampleLocalDateChain_MinusWeeks() {
	date := goda.MustLocalDateOf(2024, goda.February, 15)
	fmt.Println("Original:", date)
	fmt.Println("Minus 2 weeks:", date.Chain().MinusWeeks(2).MustGet())

	// Output:
	// Original: 2024-02-15
	// Minus 2 weeks: 2024-02-01
}

// ExampleLocalDateChain_WithDayOfMonth demonstrates changing the day of month.
func ExampleLocalDateChain_WithDayOfMonth() {
	date := goda.MustLocalDateOf(2024, goda.March, 15)

	// Change to the 1st of the month
	date2, err := date.Chain().WithDayOfMonth(1).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("First of month:", date2)

	// Change to the last day of the month
	date3, err := date.Chain().WithDayOfMonth(31).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Last of month:", date3)

	// Output:
	// First of month: 2024-03-01
	// Last of month: 2024-03-31
}

// ExampleLocalDateChain_WithDayOfYear demonstrates changing the day of year.
func ExampleLocalDateChain_WithDayOfYear() {
	date := goda.MustLocalDateOf(2024, goda.March, 15)
	fmt.Println("Original:", date)

	// Change to the 100th day of the year
	date2, err := date.Chain().WithDayOfYear(100).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Day 100:", date2)

	// Change to the 1st day of the year
	date3, err := date.Chain().WithDayOfYear(1).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Day 1:", date3)

	// Output:
	// Original: 2024-03-15
	// Day 100: 2024-04-09
	// Day 1: 2024-01-01
}

// ExampleLocalDateChain_WithMonth demonstrates changing the month.
func ExampleLocalDateChain_WithMonth() {
	date := goda.MustLocalDateOf(2024, goda.January, 31)
	fmt.Println("Original:", date)

	// Change to February (day will be clamped to 29 in leap year)
	date2, err := date.Chain().WithMonth(goda.February).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("February:", date2)

	// Change to March (day 31 is valid)
	date3, err := date.Chain().WithMonth(goda.March).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("March:", date3)

	// Output:
	// Original: 2024-01-31
	// February: 2024-02-29
	// March: 2024-03-31
}

// ExampleLocalDateChain_WithYear demonstrates changing the year.
func ExampleLocalDateChain_WithYear() {
	// Leap year date (Feb 29)
	date := goda.MustLocalDateOf(2024, goda.February, 29)
	fmt.Println("Original (leap year):", date)

	// Change to non-leap year (day will be clamped to 28)
	date2, err := date.Chain().WithYear(2023).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Non-leap year:", date2)

	// Change to another leap year (day 29 is valid)
	date3, err := date.Chain().WithYear(2020).GetResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Another leap year:", date3)

	// Output:
	// Original (leap year): 2024-02-29
	// Non-leap year: 2023-02-28
	// Another leap year: 2020-02-29
}

// ExampleLocalDate_LengthOfMonth demonstrates getting the month length.
func ExampleLocalDate_LengthOfMonth() {
	// February in leap year
	date1 := goda.MustLocalDateOf(2024, goda.February, 15)
	fmt.Printf("February 2024: %d days\n", date1.LengthOfMonth())

	// February in non-leap year
	date2 := goda.MustLocalDateOf(2023, goda.February, 15)
	fmt.Printf("February 2023: %d days\n", date2.LengthOfMonth())

	// April (30-day month)
	date3 := goda.MustLocalDateOf(2024, goda.April, 15)
	fmt.Printf("April 2024: %d days\n", date3.LengthOfMonth())

	// March (31-day month)
	date4 := goda.MustLocalDateOf(2024, goda.March, 15)
	fmt.Printf("March 2024: %d days\n", date4.LengthOfMonth())

	// Output:
	// February 2024: 29 days
	// February 2023: 28 days
	// April 2024: 30 days
	// March 2024: 31 days
}

// ExampleLocalDate_LengthOfYear demonstrates getting the year length.
func ExampleLocalDate_LengthOfYear() {
	// Leap year
	date1 := goda.MustLocalDateOf(2024, goda.March, 15)
	fmt.Printf("Year 2024: %d days\n", date1.LengthOfYear())

	// Non-leap year
	date2 := goda.MustLocalDateOf(2023, goda.March, 15)
	fmt.Printf("Year 2023: %d days\n", date2.LengthOfYear())

	// Century year divisible by 400 (leap year)
	date3 := goda.MustLocalDateOf(2000, goda.June, 15)
	fmt.Printf("Year 2000: %d days\n", date3.LengthOfYear())

	// Century year not divisible by 400 (non-leap year)
	date4 := goda.MustLocalDateOf(1900, goda.June, 15)
	fmt.Printf("Year 1900: %d days\n", date4.LengthOfYear())

	// Output:
	// Year 2024: 366 days
	// Year 2023: 365 days
	// Year 2000: 366 days
	// Year 1900: 365 days
}

// ExampleLocalTimeChain_WithField demonstrates mutating LocalTime components using WithField.
func ExampleLocalTimeChain_WithField() {
	morning := goda.MustLocalTimeOf(7, 30, 45, 123000000)

	toEvening, _ := morning.Chain().WithField(goda.FieldAmPmOfDay, goda.TemporalValueOf(1)).GetResult()
	withDifferentNanos, _ := morning.Chain().WithField(goda.FieldNanoOfSecond, goda.TemporalValueOf(987654321)).GetResult()
	exactTime, _ := morning.Chain().WithField(goda.FieldNanoOfDay, goda.TemporalValueOf(int64(22*time.Hour+1*time.Minute))).GetResult()

	fmt.Println(toEvening)
	fmt.Println(withDifferentNanos)
	fmt.Println(exactTime)

	// Output:
	// 19:30:45.123
	// 07:30:45.987654321
	// 22:01:00
}
