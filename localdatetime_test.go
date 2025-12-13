package goda

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDate_AtTime(t *testing.T) {
	date := MustLocalDateOf(2024, March, 15)
	time := MustLocalTimeOf(14, 30, 45, 123456789)
	dt := date.AtTime(time)

	assert.Equal(t, date, dt.LocalDate())
	assert.Equal(t, time, dt.LocalTime())
	assert.False(t, dt.IsZero())
}

func TestLocalTime_AtDate(t *testing.T) {
	date := MustLocalDateOf(2024, March, 15)
	time := MustLocalTimeOf(14, 30, 45, 123456789)
	dt := time.AtDate(date)

	assert.Equal(t, date, dt.LocalDate())
	assert.Equal(t, time, dt.LocalTime())
	assert.False(t, dt.IsZero())
}

func TestNewLocalDateTime(t *testing.T) {
	t.Run("valid components", func(t *testing.T) {
		dt, err := LocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, March, dt.Month())
		assert.Equal(t, 15, dt.DayOfMonth())
		assert.Equal(t, 14, dt.Hour())
		assert.Equal(t, 30, dt.Minute())
		assert.Equal(t, 45, dt.Second())
		assert.Equal(t, 123456789, dt.Nanosecond())
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := LocalDateTimeOf(2024, February, 30, 14, 30, 45, 0)
		assert.Error(t, err)
	})

	t.Run("invalid time", func(t *testing.T) {
		_, err := LocalDateTimeOf(2024, March, 15, 25, 30, 45, 0)
		assert.Error(t, err)
	})
}

func TestMustNewLocalDateTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)
			assert.Equal(t, Year(2024), dt.Year())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalDateTimeOf(2024, February, 30, 14, 30, 45, 0)
		})
	})
}

func TestNewLocalDateTimeByGoTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, time.UTC)
		dt := LocalDateTimeOfGoTime(goTime)

		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, March, dt.Month())
		assert.Equal(t, 15, dt.DayOfMonth())
		assert.Equal(t, 14, dt.Hour())
		assert.Equal(t, 30, dt.Minute())
		assert.Equal(t, 45, dt.Second())
		assert.Equal(t, 123456789, dt.Nanosecond())
	})

	t.Run("zero time", func(t *testing.T) {
		dt := LocalDateTimeOfGoTime(time.Time{})
		assert.True(t, dt.IsZero())
	})
}

func TestLocalDateTime_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var dt LocalDateTime
		assert.True(t, dt.IsZero())
	})

	t.Run("non-zero value", func(t *testing.T) {
		dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 0)
		assert.False(t, dt.IsZero())
	})

	t.Run("midnight is not zero", func(t *testing.T) {
		dt := MustLocalDateTimeOf(2024, March, 15, 0, 0, 0, 0)
		assert.False(t, dt.IsZero())
	})
}

func TestLocalDateTime_ComponentAccessors(t *testing.T) {
	dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)

	// LocalDate components
	assert.Equal(t, Year(2024), dt.Year())
	assert.Equal(t, March, dt.Month())
	assert.Equal(t, 15, dt.DayOfMonth())
	assert.Equal(t, Friday, dt.DayOfWeek())
	assert.Equal(t, 75, dt.DayOfYear())

	// LocalTime components
	assert.Equal(t, 14, dt.Hour())
	assert.Equal(t, 30, dt.Minute())
	assert.Equal(t, 45, dt.Second())
	assert.Equal(t, 123, dt.Millisecond())
	assert.Equal(t, 123456789, dt.Nanosecond())
}

func TestLocalDateTime_GoTime(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)
		goTime := dt.GoTime()

		assert.Equal(t, 2024, goTime.Year())
		assert.Equal(t, time.March, goTime.Month())
		assert.Equal(t, 15, goTime.Day())
		assert.Equal(t, 14, goTime.Hour())
		assert.Equal(t, 30, goTime.Minute())
		assert.Equal(t, 45, goTime.Second())
		assert.Equal(t, 123456789, goTime.Nanosecond())
		assert.Equal(t, time.UTC, goTime.Location())
	})

	t.Run("zero", func(t *testing.T) {
		var dt LocalDateTime
		goTime := dt.GoTime()
		assert.True(t, goTime.IsZero())
	})
}

func TestLocalDateTime_Compare(t *testing.T) {
	dt1 := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 0)
	dt2 := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 0)
	dt3 := MustLocalDateTimeOf(2024, March, 15, 14, 30, 46, 0)
	dt4 := MustLocalDateTimeOf(2024, March, 16, 14, 30, 45, 0)
	dt5 := MustLocalDateTimeOf(2024, March, 15, 15, 30, 45, 0)

	assert.Equal(t, 0, dt1.Compare(dt2))
	assert.Equal(t, -1, dt1.Compare(dt3))
	assert.Equal(t, 1, dt3.Compare(dt1))
	assert.Equal(t, -1, dt1.Compare(dt4))
	assert.Equal(t, 1, dt4.Compare(dt1))
	assert.Equal(t, -1, dt1.Compare(dt5))
	assert.Equal(t, 1, dt5.Compare(dt1))
}

func TestLocalDateTime_IsBefore_IsAfter(t *testing.T) {
	dt1 := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 0)
	dt2 := MustLocalDateTimeOf(2024, March, 15, 14, 30, 46, 0)
	dt3 := MustLocalDateTimeOf(2024, March, 16, 14, 30, 45, 0)

	assert.True(t, dt1.IsBefore(dt2))
	assert.False(t, dt2.IsBefore(dt1))
	assert.False(t, dt1.IsBefore(dt1))

	assert.True(t, dt2.IsAfter(dt1))
	assert.False(t, dt1.IsAfter(dt2))
	assert.False(t, dt1.IsAfter(dt1))

	assert.True(t, dt1.IsBefore(dt3))
	assert.True(t, dt3.IsAfter(dt1))
}

func TestLocalDateTime_String(t *testing.T) {
	tests := []struct {
		name     string
		dt       LocalDateTime
		expected string
	}{
		{
			name:     "full nanoseconds",
			dt:       MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789),
			expected: "2024-03-15T14:30:45.123456789",
		},
		{
			name:     "milliseconds",
			dt:       MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123000000),
			expected: "2024-03-15T14:30:45.123",
		},
		{
			name:     "no fractional seconds",
			dt:       MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 0),
			expected: "2024-03-15T14:30:45",
		},
		{
			name:     "midnight",
			dt:       MustLocalDateTimeOf(2024, March, 15, 0, 0, 0, 0),
			expected: "2024-03-15T00:00:00",
		},
		{
			name:     "zero value",
			dt:       LocalDateTime{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.dt.String())
		})
	}
}

func TestLocalDateTime_MarshalText(t *testing.T) {
	dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)
	text, err := dt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15T14:30:45.123456789", string(text))
}

func TestLocalDateTime_UnmarshalText(t *testing.T) {
	t.Run("valid datetime", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-03-15T14:30:45.123456789"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, March, dt.Month())
		assert.Equal(t, 15, dt.DayOfMonth())
		assert.Equal(t, 14, dt.Hour())
		assert.Equal(t, 30, dt.Minute())
		assert.Equal(t, 45, dt.Second())
		assert.Equal(t, 123456789, dt.Nanosecond())
	})

	t.Run("lowercase t separator", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-03-15t14:30:45"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("empty string", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})

	t.Run("invalid date", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-02-30T14:30:45"))
		assert.Error(t, err)
	})

	t.Run("invalid time", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-03-15T25:30:45"))
		assert.Error(t, err)
	})
}

func TestLocalDateTime_MarshalJSON(t *testing.T) {
	dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)
	jsonBytes, err := json.Marshal(dt)
	require.NoError(t, err)
	assert.Equal(t, `"2024-03-15T14:30:45.123456789"`, string(jsonBytes))
}

func TestLocalDateTime_UnmarshalJSON(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		var dt LocalDateTime
		err := json.Unmarshal([]byte(`"2024-03-15T14:30:45.123456789"`), &dt)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("null", func(t *testing.T) {
		var dt LocalDateTime
		err := json.Unmarshal([]byte(`null`), &dt)
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})
}

func TestLocalDateTime_Scan(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.Scan("2024-03-15T14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("from bytes", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.Scan([]byte("2024-03-15T14:30:45"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
	})

	t.Run("from time.LocalTime", func(t *testing.T) {
		var dt LocalDateTime
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, time.UTC)
		err := dt.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("from nil", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.Scan(nil)
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})
}

func TestLocalDateTime_Value(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)
		val, err := dt.Value()
		require.NoError(t, err)
		assert.Equal(t, "2024-03-15T14:30:45.123456789", val)
	})

	t.Run("zero", func(t *testing.T) {
		var dt LocalDateTime
		val, err := dt.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestParseLocalDateTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dt, err := LocalDateTimeParse("2024-03-15T14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())

		dt, err = LocalDateTimeParse("2024-03-15T14:30:45")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
	})

	t.Run("empty", func(t *testing.T) {
		dt, err := LocalDateTimeParse("")
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})
}

func TestMustParseLocalDateTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			dt := MustLocalDateTimeParse("2024-03-15T14:30:45.123456789")
			assert.Equal(t, Year(2024), dt.Year())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalDateTimeParse("invalid")
		})
	})
}

func TestLocalDateTimeNow(t *testing.T) {
	now := LocalDateTimeNow()
	assert.False(t, now.IsZero())

	// Verify components are in valid ranges
	assert.True(t, now.Year() >= 1970 && now.Year() <= 9999)
	assert.True(t, now.Hour() >= 0 && now.Hour() < 24)
}

func TestLocalDateTimeNowUTC(t *testing.T) {
	nowUTC := LocalDateTimeNowUTC()
	assert.False(t, nowUTC.IsZero())
}

func TestLocalDateTimeNowIn(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	nowTokyo := LocalDateTimeNowIn(tokyo)
	assert.False(t, nowTokyo.IsZero())
}

func TestLocalDateTime_IsLeapYear(t *testing.T) {
	dt2024 := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 0)
	dt2023 := MustLocalDateTimeOf(2023, March, 15, 14, 30, 45, 0)

	assert.True(t, dt2024.IsLeapYear())
	assert.False(t, dt2023.IsLeapYear())
}

//go:embed TestFieldGetter.txt
var TestlocaldatetimeGetfieldJavamatchdata string

func TestLocalDateTime_GetField_JavaMatch(t *testing.T) {
	var fields = []Field{FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldHourOfDay, FieldClockHourOfDay, FieldAmPmOfDay, FieldDayOfWeek, FieldAlignedDayOfWeekInMonth, FieldAlignedDayOfWeekInYear, FieldDayOfMonth, FieldDayOfYear, FieldEpochDay, FieldAlignedWeekOfMonth, FieldAlignedWeekOfYear, FieldMonthOfYear, FieldProlepticMonth, FieldYearOfEra, FieldYear, FieldEra}
	for lineNumber, line := range strings.Split(TestlocaldatetimeGetfieldJavamatchdata, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")
		t.Run(fmt.Sprintf("line_%d", lineNumber+1), func(t *testing.T) {
			dt, e := LocalDateTimeParse(cols[0])
			if e != nil {
				t.Fatal(e)
			}
			for i, v := range cols[1:] {
				if dt.GetField(fields[i]).Unsupported() {
					continue
				}
				t.Run(fields[i].String(), func(t *testing.T) {
					if !assert.Equal(t, v, fmt.Sprint(dt.GetField(fields[i]).Int64())) {
						t.Log(fields[i].String(), line)
						t.Failed()
						return
					}
				})
			}
		})
	}
}

func TestLocalDateTime_WithField_JavaMatch(t *testing.T) {
	var fields = []Field{FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldHourOfDay, FieldClockHourOfDay, FieldAmPmOfDay, FieldDayOfWeek, FieldAlignedDayOfWeekInMonth, FieldAlignedDayOfWeekInYear, FieldDayOfMonth, FieldDayOfYear, FieldEpochDay, FieldAlignedWeekOfMonth, FieldAlignedWeekOfYear, FieldMonthOfYear, FieldProlepticMonth, FieldYearOfEra, FieldYear, FieldEra}
	for lineNumber, line := range strings.Split(TestlocaldatetimeGetfieldJavamatchdata, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")
		dt, e := LocalDateTimeParse(cols[0])
		if e != nil {
			t.Fatal(e)
		}
		t.Run(fmt.Sprintf("line_%d", lineNumber+1), func(t *testing.T) {
			for i, v := range cols[1:] {
				var field = fields[i]
				if dt.GetField(field).Unsupported() {
					continue
				}
				switch field {
				case FieldMicroOfSecond, FieldMilliOfSecond, FieldMilliOfDay, FieldMicroOfDay:
					// subfield cannot be tested by this testcase
					continue
				}
				t.Run(field.String(), func(t *testing.T) {
					var ndt LocalDateTime
					ndt, e = dt.Chain().WithField(field, TemporalValueOf(1)).GetResult()
					if e != nil {
						t.Fatal(e)
					}
					if ndt == dt && v != "1" {
						t.Fatalf("data unchanged, %s, %s", dt, v)
					}
					var nv int
					nv, e = strconv.Atoi(v)
					if e != nil {
						t.Fatal(e)
					}
					if fmt.Sprint(nv) != v {
						panic("invalid value: " + v)
					}
					ndt, e = ndt.Chain().WithField(field, TemporalValueOf(nv)).GetResult()
					if e != nil {
						t.Fatal(e)
					}
					if ndt != dt {
						t.Fatalf("%s != %s (%s)", ndt, dt, v)
						return
					}
				})
			}
		})
	}
}

func TestLocalDateTime_ValuePostgres(t *testing.T) {
	var pg = GetPG(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustLocalDateTimeParse("2000-12-29 12:00:00")
		var actual LocalDateTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT $1::timestamp without time zone, $1::timestamp without time zone = '2000-12-29 12:00:00'", expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalDateTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT NULL::timestamp without time zone, $1::timestamp without time zone is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalDateTime_ValueMySQL(t *testing.T) {
	var pg = GetMySQL(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustLocalDateTimeParse("2000-12-29 12:00:00")
		var actual LocalDateTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT CAST(? AS DATETIME), CAST(? AS DATETIME)  = '2000-12-29 12:00:00'", expected, expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalDateTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT CAST(NULL AS DATETIME), CAST(? AS DATETIME) is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}
