# goda

English | [中文](README.zh.md)

[![CI](https://github.com/iseki0/goda/workflows/CI/badge.svg)](https://github.com/iseki0/goda/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/iseki0/goda.svg)](https://pkg.go.dev/github.com/iseki0/goda)
[![Go Report Card](https://goreportcard.com/badge/github.com/iseki0/goda)](https://goreportcard.com/report/github.com/iseki0/goda)
[![codecov](https://codecov.io/gh/iseki0/goda/graph/badge.svg?token=TBHUZUY561)](https://codecov.io/gh/iseki0/goda)

> **ThreeTen/JSR-310** model in Go

A Go implementation inspired by Java's `java.time` package (JSR-310), providing immutable date and time types that are **type-safe** and **easy to use**.

## Features

### Core Types

- 📅 **LocalDate**: Date without time (e.g., `2024-03-15`)
- ⏰ **LocalTime**: Time without date (e.g., `14:30:45.123456789`)
- 📆 **LocalDateTime**: Date-time without timezone (e.g., `2024-03-15T14:30:45.123456789`)
- 🌐 **ZoneOffset**: Time-zone offset from Greenwich/UTC (e.g., `+08:00`, `-05:00`, `Z`)
- 🌍 **OffsetDateTime**: Date-time with offset (e.g., `2024-03-15T14:30:45.123456789+01:00`)
- 🔢 **Field**: Enumeration of date-time fields (like Java's `ChronoField`)
- 🔍 **TemporalAccessor**: Universal interface for querying temporal objects
- 📊 **TemporalValue**: Type-safe wrapper for field values with validation state

### Key Features

- ✅ **ISO 8601 basic format** support (yyyy-MM-dd, HH:mm:ss[.nnnnnnnnn], combined with 'T')
- ✅ **Java.time compatible formatting**: Fractional seconds aligned to 3-digit boundaries (milliseconds, microseconds, nanoseconds)
- ✅ **Full JSON and SQL** database integration
- ✅ **Date arithmetic**: Add/subtract days, months, years with overflow handling
- ✅ **Type-safe field access**: Query any field with `TemporalValue` return type that validates support and overflow
- ✅ **TemporalAccessor interface**: Universal query pattern across all temporal types
- ✅ **Chain operations**: Fluent API with error handling for complex mutations
- ✅ **Immutable**: All operations return new values
- ✅ **Type-safe**: Compile-time safety with distinct types
- ✅ **Zero-value friendly**: Zero values are properly handled

## Installation

```bash
go get github.com/iseki0/goda
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/iseki0/goda"
)

func main() {
    // Create dates and times
    date := goda.MustLocalDateOf(2024, goda.March, 15)
    time := goda.MustLocalTimeOf(14, 30, 45, 123456789)
    datetime := date.AtTime(time)  // or time.AtDate(date)
    
    fmt.Println(date)     // 2024-03-15
    fmt.Println(time)     // 14:30:45.123456789
    fmt.Println(datetime) // 2024-03-15T14:30:45.123456789
    
    // Create from components directly
    datetime2 := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 123456789)
    
    // With time zone offset
    offset := goda.MustZoneOffsetOfHours(1)  // +01:00
    offsetDateTime := datetime.AtOffset(offset)
    fmt.Println(offsetDateTime) // 2024-03-15T14:30:45.123456789+01:00
    
    // Parse from strings
    date, _ = goda.LocalDateParse("2024-03-15")
    time = goda.MustLocalTimeParse("14:30:45.123456789")
    datetime = goda.MustLocalDateTimeParse("2024-03-15T14:30:45")
    
    // Get current date/time
    today := goda.LocalDateNow()
    now := goda.LocalTimeNow()
    currentDateTime := goda.LocalDateTimeNow()
    currentOffsetDateTime := goda.OffsetDateTimeNow()
    
    // Date arithmetic
    tomorrow := today.Chain().PlusDays(1).MustGet()
    nextMonth := today.Chain().PlusMonths(1).MustGet()
    nextYear := today.Chain().PlusYears(1).MustGet()
    
    // Comparisons
    if tomorrow.IsAfter(today) {
        fmt.Println("Tomorrow is after today!")
    }
}
```

### Field Access with TemporalValue

Access individual date-time fields using the `Field` enumeration with type-safe `TemporalValue` returns:

```go
date := goda.MustLocalDateOf(2024, goda.March, 15)

// Check field support
fmt.Println(date.IsSupportedField(goda.FieldDayOfMonth))  // true
fmt.Println(date.IsSupportedField(goda.FieldHourOfDay))   // false

// Get field values with validation
year := date.GetField(goda.FieldYear)
if year.Valid() {
    fmt.Println("Year:", year.Int64())  // 2024
}

dayOfWeek := date.GetField(goda.FieldDayOfWeek)
if dayOfWeek.Valid() {
    fmt.Println("Day of week:", dayOfWeek.Int())  // 5 (Friday)
}

// Unsupported fields return unsupported TemporalValue
hourOfDay := date.GetField(goda.FieldHourOfDay)
if hourOfDay.Unsupported() {
    fmt.Println("Hour field is not supported for LocalDate")
}

// Time fields
time := goda.MustLocalTimeOf(14, 30, 45, 123456789)
hour := time.GetField(goda.FieldHourOfDay)
if hour.Valid() {
    fmt.Println("Hour:", hour.Int())  // 14
}

nanoOfDay := time.GetField(goda.FieldNanoOfDay)
if nanoOfDay.Valid() {
    fmt.Println("Nanoseconds since midnight:", nanoOfDay.Int64())
}
```

**TemporalValue API:**
- `Valid() bool`: Returns true if the field is supported and no overflow occurred
- `Unsupported() bool`: Returns true if the field is not supported by this temporal type
- `Overflow() bool`: Returns true if the field value overflowed (reserved for future use)
- `Int64() int64`: Get the raw value as int64
- `Int() int`: Get the value as int (for convenience)

**Why TemporalValue?**

The `TemporalValue` return type provides type-safe field queries that prevent silent errors:
- **Explicit validation**: Check `Valid()` before using the value
- **Clear error semantics**: Distinguish between unsupported fields and actual errors
- **Future-proof**: Ready for overflow detection when needed
- **No silent zeros**: Unlike raw `int64` returns, you can distinguish between "0" and "unsupported"

### TemporalAccessor Interface

All temporal types implement the `TemporalAccessor` interface, providing a uniform query pattern:

```go
// TemporalAccessor provides read-only access to temporal fields
type TemporalAccessor interface {
    IsZero() bool
    IsSupportedField(field Field) bool
    GetField(field Field) TemporalValue
}

// Write generic functions that work with any temporal type
func printYear(t goda.TemporalAccessor) {
    if year := t.GetField(goda.FieldYear); year.Valid() {
        fmt.Printf("Year: %d\n", year.Int())
    }
}

// Works with LocalDate, LocalTime, or LocalDateTime
printYear(goda.LocalDateNow())
printYear(goda.LocalDateTimeNow())
```

### Chain Operations

All temporal types support chain operations for fluent, error-handled mutations. Chain operations allow you to perform multiple modifications in a single expression with proper error handling:

```go
// Chain multiple operations fluently
dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 123456789)

// Chain date and time modifications
meetingTime := dt.Chain().
    PlusDays(7).              // Next week
    WithHour(16).             // At 4 PM
    WithMinute(0).            // On the hour
    WithSecond(0).            // No seconds
    WithNano(0).              // No nanoseconds
    MustGet()                 // Get result (panics on error)

fmt.Println("Meeting scheduled for:", meetingTime)

// Error handling with chains
result, err := dt.Chain().
    PlusMonths(1).
    WithDayOfMonth(32).       // Invalid day - will cause error
    GetResult()               // Returns (zero value, error)

if err != nil {
    fmt.Println("Invalid operation:", err)
    // Use fallback
    validTime := dt.Chain().
        PlusMonths(1).
        WithDayOfMonth(31).   // Valid day
        GetOrElse(dt)         // Fallback to original if error
}
```

### JSON Serialization

```go
type Event struct {
    Name        string                `json:"name"`
    Date        goda.LocalDate        `json:"date"`
    Time        goda.LocalTime        `json:"time"`
    CreatedAt   goda.LocalDateTime    `json:"created_at"`
    ScheduledAt goda.OffsetDateTime   `json:"scheduled_at"`  // With timezone
}

event := Event{
    Name:        "Meeting",
    Date:        goda.MustLocalDateOf(2024, goda.March, 15),
    Time:        goda.MustLocalTimeOf(14, 30, 0, 0),
    CreatedAt:   goda.MustLocalDateTimeParse("2024-03-15T14:30:00"),
    ScheduledAt: goda.MustOffsetDateTimeParse("2024-03-15T14:30:00+08:00"),
}

jsonData, _ := json.Marshal(event)
// {"name":"Meeting","date":"2024-03-15","time":"14:30:00",
//  "created_at":"2024-03-15T14:30:00","scheduled_at":"2024-03-15T14:30:00+08:00"}
```

### Database Integration

```go
type Record struct {
    ID          int64
    CreatedAt   goda.LocalDateTime
    Date        goda.LocalDate
    UpdatedAt   goda.OffsetDateTime  // With timezone for audit trails
}

// Works with database/sql - implements sql.Scanner and driver.Valuer
db.QueryRow("SELECT id, created_at, date, updated_at FROM records WHERE id = ?", 1).Scan(
    &record.ID, &record.CreatedAt, &record.Date, &record.UpdatedAt,
)

// Insert with offset datetime
offset := goda.MustZoneOffsetOfHours(8)
now := goda.OffsetDateTimeNow()
db.Exec("INSERT INTO records (created_at, updated_at) VALUES (?, ?)",
    goda.LocalDateTimeNow(), now)
```

## API Overview

### Core Types

| Type               | Description                             | Example                                |
|--------------------|-----------------------------------------|----------------------------------------|
| `LocalDate`        | Date without time                       | `2024-03-15`                           |
| `LocalTime`        | Time without date                       | `14:30:45.123456789`                   |
| `LocalDateTime`    | Date-time without timezone              | `2024-03-15T14:30:45`                  |
| `ZoneOffset`       | Time-zone offset from Greenwich/UTC     | `+08:00`, `-05:00`, `Z`                |
| `OffsetDateTime`   | Date-time with offset from UTC          | `2024-03-15T14:30:45+08:00`            |
| `Month`            | Month of year (1-12)                    | `March`                                |
| `Year`             | Year                                    | `2024`                                 |
| `DayOfWeek`        | Day of week (1=Monday, 7=Sunday)        | `Friday`                               |
| `Field`            | Date-time field enumeration             | `HourOfDay`, `DayOfMonth`              |
| `TemporalAccessor` | Interface for querying temporal objects | Implemented by all temporal types      |
| `TemporalValue`    | Type-safe field value with validation   | Returned by `GetField()`               |
| `Error`            | Structured error with context           | Provides detailed error information    |
| `LocalDateChain`   | Chain operations for LocalDate          | `date.Chain().PlusDays(1).MustGet()`   |
| `LocalTimeChain`   | Chain operations for LocalTime          | `time.Chain().PlusHours(1).MustGet()`  |
| `LocalDateTimeChain`| Chain operations for LocalDateTime      | `dt.Chain().PlusDays(1).MustGet()`     |
| `OffsetDateTimeChain`| Chain operations for OffsetDateTime     | `odt.Chain().PlusHours(1).MustGet()`   |

### Format Specification

This package uses ISO 8601 basic calendar date and time formats (not the full specification):

**LocalDate**: `yyyy-MM-dd` (e.g., "2024-03-15")  
Only Gregorian calendar dates. No week dates (YYYY-Www-D) or ordinal dates (YYYY-DDD).

**LocalTime**: `HH:mm:ss[.nnnnnnnnn]` (e.g., "14:30:45.123456789")  
24-hour format. Fractional seconds up to nanoseconds. Fractional seconds are aligned to 3-digit boundaries (milliseconds, microseconds, nanoseconds) for Java.time compatibility: 100ms → "14:30:45.100", 123.4ms → "14:30:45.123400". Parsing accepts any length of fractional seconds (e.g., "14:30:45.1" → 100ms).

**LocalDateTime**: `yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]` (e.g., "2024-03-15T14:30:45.123456789")  
Combined with 'T' separator (lowercase 't' accepted when parsing).

**ZoneOffset**: `±HH:mm[:ss]` or `Z` for UTC (e.g., "+08:00", "-05:30", "Z")  
Hours must be in range [-18, 18], minutes and seconds in [0, 59]. Compact formats (±HH, ±HHMM, ±HHMMSS) are also supported.

**OffsetDateTime**: `yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss]` (e.g., "2024-03-15T14:30:45+08:00")  
Combines LocalDateTime and ZoneOffset. 'Z' is accepted as UTC offset.

### Time Formatting

Time values use ISO 8601 format with **Java.time compatible** fractional second alignment:

| Precision     | Digits | Example                                    |
|---------------|--------|--------------------------------------------|
| Whole seconds | 0      | `14:30:45`                                 |
| Milliseconds  | 3      | `14:30:45.100`, `14:30:45.123`             |
| Microseconds  | 6      | `14:30:45.123400`, `14:30:45.123456`       |
| Nanoseconds   | 9      | `14:30:45.000000001`, `14:30:45.123456789` |

Fractional seconds are automatically aligned to 3-digit boundaries (milliseconds, microseconds, nanoseconds), matching Java's `LocalTime` behavior. Parsing accepts any length of fractional seconds.

### Field Constants (30 fields)

**Time Fields**: `NanoOfSecond`, `NanoOfDay`, `MicroOfSecond`, `MicroOfDay`, `MilliOfSecond`, `MilliOfDay`, `SecondOfMinute`, `SecondOfDay`, `MinuteOfHour`, `MinuteOfDay`, `HourOfAmPm`, `ClockHourOfAmPm`, `HourOfDay`, `ClockHourOfDay`, `AmPmOfDay`

**Date Fields**: `DayOfWeekField`, `DayOfMonth`, `DayOfYear`, `EpochDay`, `AlignedDayOfWeekInMonth`, `AlignedDayOfWeekInYear`, `AlignedWeekOfMonth`, `AlignedWeekOfYear`, `MonthOfYear`, `ProlepticMonth`, `YearOfEra`, `YearField`, `Era`

**Other Fields**: `InstantSeconds`, `OffsetSeconds`

### Implemented Interfaces

All temporal types (`LocalDate`, `LocalTime`, `LocalDateTime`, `OffsetDateTime`) implement:
- `TemporalAccessor`: Universal query interface with `GetField(field Field) TemporalValue`
- `fmt.Stringer`
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler`
- `encoding.TextAppender`
- `json.Marshaler` / `json.Unmarshaler`
- `sql.Scanner` / `driver.Valuer`

## Design Philosophy

This package follows the **ThreeTen/JSR-310** model (Java's `java.time` package), providing date and time types that are:

- **Immutable**: All operations return new values
- **Type-safe**: Distinct types for date, time, and datetime
- **Simple formats**: Uses ISO 8601 basic formats (not the full complex specification)
- **Database-friendly**: Direct SQL integration
- **Field-based access**: Universal field access pattern via `TemporalAccessor` interface
- **Safe field queries**: `TemporalValue` return type validates field support and prevents silent errors
- **Zero-value safe**: Zero values are properly handled throughout

### When to Use Each Type

**LocalDate, LocalTime, LocalDateTime** - Use when timezone is not relevant:
- **Birthdays**: "March 15" means March 15 everywhere
- **Business hours**: "9:00 AM - 5:00 PM" in local context
- **Schedules**: "Meeting at 2:30 PM" without timezone concerns
- **Calendar dates**: Historical dates, recurring events

**OffsetDateTime** - Use when you need a fixed offset from UTC:
- **API timestamps**: REST APIs often use RFC3339 with offsets
- **Audit logs**: Record exact moment with original timezone offset
- **Event scheduling**: When timezone offset matters but DST transitions don't
- **International coordination**: "The meeting is at 14:00 UTC+1"

**ZoneOffset** - Use to represent timezone offsets:
- **Fixed offsets**: +08:00, -05:00, Z (UTC)
- **No DST handling**: Use when you don't need daylight saving time rules
- **Simple offset arithmetic**: Convert between different offsets

For full timezone support with DST transitions, use `ZonedDateTime` (coming soon).

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/iseki0/goda).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
