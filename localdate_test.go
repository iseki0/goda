package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDate_Epoch(t *testing.T) {
	var check = func(i int64, tt time.Time) bool {
		var st = MustLocalDateOf(Year(tt.Year()), Month(tt.Month()), tt.Day())
		if !assert.Equal(t, i, st.UnixEpochDays(), tt) {
			return false
		}
		if !assert.Equal(t, st, MustLocalDateOfUnixEpochDays(i), tt) {
			return false
		}
		if !assert.Equal(t, st.DayOfWeek().GoWeekday(), tt.Weekday(), tt) {
			return false
		}
		return assert.Equal(t, tt.Unix()/(24*60*60), st.UnixEpochDays(), tt)
	}
	var begin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := int64(0); i < 100_0000; i++ {
		if !check(i, begin) {
			break
		}
		begin = begin.AddDate(0, 0, 1)
	}
	begin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	// negative
	for i := int64(0); i > -100_0000; i-- {
		if !check(i, begin) {
			break
		}
		begin = begin.AddDate(0, 0, -1)
	}
}

func TestNewLocalDate(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		d, err := LocalDateOf(2024, January, 1)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), d.Year())
		assert.Equal(t, January, d.Month())
		assert.Equal(t, 1, d.DayOfMonth())

		d, err = LocalDateOf(2024, February, 29) // leap year
		require.NoError(t, err)
		assert.Equal(t, 29, d.DayOfMonth())
	})

	t.Run("invalid day of month", func(t *testing.T) {
		_, err := LocalDateOf(2024, January, 32)
		assert.Error(t, err)

		_, err = LocalDateOf(2023, February, 29) // not a leap year
		assert.Error(t, err)

		_, err = LocalDateOf(2024, February, 30)
		assert.Error(t, err)

		_, err = LocalDateOf(2024, April, 31)
		assert.Error(t, err)
	})

	t.Run("invalid month", func(t *testing.T) {
		_, err := LocalDateOf(2024, Month(0), 1)
		assert.Error(t, err)

		_, err = LocalDateOf(2024, Month(13), 1)
		assert.Error(t, err)
	})
}

func TestMustNewLocalDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		assert.NotPanics(t, func() {
			d := MustLocalDateOf(2024, March, 15)
			assert.Equal(t, Year(2024), d.Year())
		})
	})

	t.Run("invalid date panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalDateOf(2024, January, 32)
		})
	})
}

func TestLocalDate_IsZero(t *testing.T) {
	var zero LocalDate
	assert.True(t, zero.IsZero())

	d := MustLocalDateOf(2024, January, 1)
	assert.False(t, d.IsZero())
}

func TestLocalDate_IsLeapYear(t *testing.T) {
	tests := []struct {
		year   Year
		isLeap bool
	}{
		{2000, true},  // divisible by 400
		{2004, true},  // divisible by 4
		{2100, false}, // divisible by 100 but not 400
		{2023, false}, // not divisible by 4
		{2024, true},  // divisible by 4
		{1900, false}, // divisible by 100 but not 400
	}

	for _, tt := range tests {
		d := MustLocalDateOf(tt.year, January, 1)
		assert.Equal(t, tt.isLeap, d.IsLeapYear(), "year %d", tt.year)
	}
}

func TestLocalDate_DayOfYear(t *testing.T) {
	tests := []struct {
		date      LocalDate
		dayOfYear int
	}{
		{MustLocalDateOf(2024, January, 1), 1},
		{MustLocalDateOf(2024, January, 31), 31},
		{MustLocalDateOf(2024, February, 1), 32},
		{MustLocalDateOf(2024, February, 29), 60}, // leap year
		{MustLocalDateOf(2023, February, 28), 59}, // non-leap year
		{MustLocalDateOf(2024, March, 1), 61},     // leap year
		{MustLocalDateOf(2023, March, 1), 60},     // non-leap year
		{MustLocalDateOf(2024, December, 31), 366},
		{MustLocalDateOf(2023, December, 31), 365},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.dayOfYear, tt.date.DayOfYear(), "date: %s", tt.date.String())
	}

	var zero LocalDate
	assert.Equal(t, 0, zero.DayOfYear())
}

func TestLocalDate_DayOfWeek(t *testing.T) {
	tests := []struct {
		date      LocalDate
		dayOfWeek DayOfWeek
	}{
		{MustLocalDateOf(2024, 11, 5), Tuesday},  // Known date
		{MustLocalDateOf(2024, 1, 1), Monday},    // New Year 2024
		{MustLocalDateOf(2023, 1, 1), Sunday},    // New Year 2023
		{MustLocalDateOf(2000, 1, 1), Saturday},  // Y2K
		{MustLocalDateOf(1970, 1, 1), Thursday},  // Unix epoch
		{MustLocalDateOf(2024, 2, 29), Thursday}, // Leap day 2024
	}

	for _, tt := range tests {
		assert.Equal(t, tt.dayOfWeek, tt.date.DayOfWeek(), "date: %s", tt.date.String())
	}

	var zero LocalDate
	assert.Equal(t, DayOfWeek(0), zero.DayOfWeek())
}

func TestLocalDate_Compare(t *testing.T) {
	d1 := MustLocalDateOf(2024, March, 15)
	d2 := MustLocalDateOf(2024, March, 15)
	d3 := MustLocalDateOf(2024, March, 16)
	d4 := MustLocalDateOf(2024, February, 15)
	d5 := MustLocalDateOf(2023, March, 15)

	assert.Equal(t, 0, d1.Compare(d2))
	assert.Equal(t, -1, d1.Compare(d3))
	assert.Equal(t, 1, d3.Compare(d1))
	assert.Equal(t, 1, d1.Compare(d4))
	assert.Equal(t, 1, d1.Compare(d5))
}

func TestLocalDate_IsBefore(t *testing.T) {
	d1 := MustLocalDateOf(2024, March, 15)
	d2 := MustLocalDateOf(2024, March, 16)
	d3 := MustLocalDateOf(2024, March, 15)

	assert.True(t, d1.IsBefore(d2))
	assert.False(t, d2.IsBefore(d1))
	assert.False(t, d1.IsBefore(d3))
}

func TestLocalDate_IsAfter(t *testing.T) {
	d1 := MustLocalDateOf(2024, March, 15)
	d2 := MustLocalDateOf(2024, March, 16)
	d3 := MustLocalDateOf(2024, March, 15)

	assert.False(t, d1.IsAfter(d2))
	assert.True(t, d2.IsAfter(d1))
	assert.False(t, d1.IsAfter(d3))
}

func TestLocalDate_GoTime(t *testing.T) {
	d := MustLocalDateOf(2024, March, 15)
	goTime := d.GoTime()

	assert.Equal(t, 2024, goTime.Year())
	assert.Equal(t, time.March, goTime.Month())
	assert.Equal(t, 15, goTime.Day())
	assert.Equal(t, 0, goTime.Hour())
	assert.Equal(t, 0, goTime.Minute())
	assert.Equal(t, 0, goTime.Second())
	assert.Equal(t, time.UTC, goTime.Location())

	var zero LocalDate
	assert.True(t, zero.GoTime().IsZero())
}

func TestNewLocalDateByGoTime(t *testing.T) {
	goTime := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	d := LocalDateOfGoTime(goTime)

	assert.Equal(t, Year(2024), d.Year())
	assert.Equal(t, March, d.Month())
	assert.Equal(t, 15, d.DayOfMonth())

	// Test with zero time
	d = LocalDateOfGoTime(time.Time{})
	assert.True(t, d.IsZero())
}

func TestLocalDate_String(t *testing.T) {
	tests := []struct {
		date     LocalDate
		expected string
	}{
		{MustLocalDateOf(2024, March, 15), "2024-03-15"},
		{MustLocalDateOf(2024, January, 1), "2024-01-01"},
		{MustLocalDateOf(2024, December, 31), "2024-12-31"},
		{MustLocalDateOf(1999, June, 5), "1999-06-05"},
		{LocalDate{}, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.date.String())
	}
}

func TestLocalDate_MarshalText(t *testing.T) {
	d := MustLocalDateOf(2024, March, 15)
	text, err := d.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15", string(text))

	var zero LocalDate
	text, err = zero.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(text))
}

func TestLocalDate_UnmarshalText(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte("2024-03-15"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2024, March, 15), d)

		err = d.UnmarshalText([]byte("1999-12-31"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(1999, December, 31), d)
	})

	t.Run("empty string", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("invalid format", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte("2024/03/15"))
		assert.Error(t, err)

		err = d.UnmarshalText([]byte("2024-3-15"))
		assert.Error(t, err)

		err = d.UnmarshalText([]byte("not-a-date"))
		assert.Error(t, err)
	})

	t.Run("invalid date", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte("2024-02-30"))
		assert.Error(t, err)

		err = d.UnmarshalText([]byte("2024-13-01"))
		assert.Error(t, err)
	})
}

func TestLocalDate_MarshalJSON(t *testing.T) {
	d := MustLocalDateOf(2024, March, 15)
	data, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Equal(t, `"2024-03-15"`, string(data))

	var zero LocalDate
	data, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `""`, string(data))
}

func TestLocalDate_UnmarshalJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`"2024-03-15"`), &d)
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2024, March, 15), d)
	})

	t.Run("empty string", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`""`), &d)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("null", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`null`), &d)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`"invalid-date"`), &d)
		assert.Error(t, err)
	})
}

func TestLocalDate_Scan(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		var d LocalDate
		err := d.Scan(nil)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("string value", func(t *testing.T) {
		var d LocalDate
		err := d.Scan("2024-03-15")
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2024, March, 15), d)
	})

	t.Run("byte slice value", func(t *testing.T) {
		var d LocalDate
		err := d.Scan([]byte("2024-03-15"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2024, March, 15), d)
	})

	t.Run("time.LocalTime value", func(t *testing.T) {
		var d LocalDate
		goTime := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
		err := d.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2024, March, 15), d)
	})
}

func TestLocalDate_Value(t *testing.T) {
	d := MustLocalDateOf(2024, March, 15)
	val, err := d.Value()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15", val)

	var zero LocalDate
	val, err = zero.Value()
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestLocalDate_AppendText(t *testing.T) {
	d := MustLocalDateOf(2024, March, 15)
	buf := []byte("LocalDate: ")
	buf, err := d.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalDate: 2024-03-15", string(buf))

	var zero LocalDate
	buf = []byte("LocalDate: ")
	buf, err = zero.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalDate: ", string(buf))
}

func TestLocalDateNow(t *testing.T) {
	// Test that LocalDateNow() returns a valid date
	today := LocalDateNow()
	assert.False(t, today.IsZero(), "LocalDateNow should not be zero")

	// Test that it matches time.Now()
	now := time.Now()
	expected := LocalDateOfGoTime(now)

	// Allow for the possibility that the date changed between calls
	// (very unlikely but possible at midnight)
	diff := today.Compare(expected)
	assert.True(t, diff >= -1 && diff <= 1, "LocalDateNow should be within 1 day of current time")
}

func TestLocalDateNowUTC(t *testing.T) {
	todayUTC := LocalDateNowUTC()
	assert.False(t, todayUTC.IsZero(), "LocalDateNowUTC should not be zero")

	// Test that it matches time.Now().UTC()
	now := time.Now().UTC()
	expected := LocalDateOfGoTime(now)

	diff := todayUTC.Compare(expected)
	assert.True(t, diff >= -1 && diff <= 1, "LocalDateNowUTC should be within 1 day of current UTC time")
}

func TestParseLocalDate(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		date, err := LocalDateParse("2024-03-15")
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2024, March, 15), date)

		date, err = LocalDateParse("1999-12-31")
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(1999, December, 31), date)

		date, err = LocalDateParse("2000-02-29") // leap year
		require.NoError(t, err)
		assert.Equal(t, MustLocalDateOf(2000, February, 29), date)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := LocalDateParse("2024/03/15")
		assert.Error(t, err)

		_, err = LocalDateParse("2024-3-15")
		assert.Error(t, err)

		_, err = LocalDateParse("not-a-date")
		assert.Error(t, err)
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := LocalDateParse("2024-02-30")
		assert.Error(t, err)

		_, err = LocalDateParse("2024-13-01")
		assert.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		date, err := LocalDateParse("")
		require.NoError(t, err)
		assert.True(t, date.IsZero())
	})
}

func TestMustParseLocalDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		assert.NotPanics(t, func() {
			date := MustLocalDateParse("2024-03-15")
			assert.Equal(t, MustLocalDateOf(2024, March, 15), date)
		})
	})

	t.Run("invalid date panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalDateParse("2024-02-30")
		})

		assert.Panics(t, func() {
			MustLocalDateParse("invalid")
		})
	})
}

func TestLocalDateNowIn(t *testing.T) {
	// Test with different time zones
	locations := []struct {
		name string
		loc  *time.Location
	}{
		{"UTC", time.UTC},
		{"Local", time.Local},
	}

	for _, tt := range locations {
		t.Run(tt.name, func(t *testing.T) {
			todayIn := LocalDateNowIn(tt.loc)
			assert.False(t, todayIn.IsZero(), "LocalDateNowIn should not be zero")

			now := time.Now().In(tt.loc)
			expected := LocalDateOfGoTime(now)

			diff := todayIn.Compare(expected)
			assert.True(t, diff >= -1 && diff <= 1, "LocalDateNowIn should be within 1 day of current time in specified zone")
		})
	}
}

func TestLocalDate_ValuePostgres(t *testing.T) {
	var pg = GetPG(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustLocalDateParse("2000-12-29")
		var actual LocalDate
		var expectedTrue bool
		var e = pg.QueryRow("SELECT $1::date, $1::date = '2000-12-29'", expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalDate
		var expectedTrue bool
		var e = pg.QueryRow("SELECT NULL::date, $1::date is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalDate_ValueMySQL(t *testing.T) {
	var mysql = GetMySQL(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustLocalDateParse("2000-12-29")
		var actual LocalDate
		var expectedTrue bool
		var e = mysql.QueryRow("SELECT CAST(? AS DATE), CAST(? AS DATE) = '2000-12-29'", expected, expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalDate
		var expectedTrue bool
		var e = mysql.QueryRow("SELECT CAST(NULL AS DATE), CAST(? AS DATE) is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalDate_LengthOfMonth(t *testing.T) {
	t.Run("31-day months", func(t *testing.T) {
		months31 := []Month{January, March, May, July, August, October, December}
		for _, month := range months31 {
			date := MustLocalDateOf(2024, month, 1)
			assert.Equal(t, 31, date.LengthOfMonth(), "month: %d", month)
		}
	})

	t.Run("30-day months", func(t *testing.T) {
		months30 := []Month{April, June, September, November}
		for _, month := range months30 {
			date := MustLocalDateOf(2024, month, 1)
			assert.Equal(t, 30, date.LengthOfMonth(), "month: %d", month)
		}
	})

	t.Run("February in leap year", func(t *testing.T) {
		date := MustLocalDateOf(2024, February, 1)
		assert.Equal(t, 29, date.LengthOfMonth())
	})

	t.Run("February in non-leap year", func(t *testing.T) {
		date := MustLocalDateOf(2023, February, 1)
		assert.Equal(t, 28, date.LengthOfMonth())
	})

	t.Run("various leap years", func(t *testing.T) {
		tests := []struct {
			year     Year
			expected int
		}{
			{2000, 29}, // divisible by 400
			{2004, 29}, // divisible by 4
			{2100, 28}, // divisible by 100 but not 400
			{2023, 28}, // not divisible by 4
		}

		for _, tt := range tests {
			date := MustLocalDateOf(tt.year, February, 1)
			assert.Equal(t, tt.expected, date.LengthOfMonth(), "year: %d", tt.year)
		}
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalDate
		assert.Equal(t, 0, zero.LengthOfMonth())
	})
}

func TestLocalDate_LengthOfYear(t *testing.T) {
	t.Run("leap years", func(t *testing.T) {
		leapYears := []Year{2000, 2004, 2024, 2400}
		for _, year := range leapYears {
			date := MustLocalDateOf(year, January, 1)
			assert.Equal(t, 366, date.LengthOfYear(), "year: %d", year)
		}
	})

	t.Run("non-leap years", func(t *testing.T) {
		nonLeapYears := []Year{1900, 2001, 2023, 2100}
		for _, year := range nonLeapYears {
			date := MustLocalDateOf(year, January, 1)
			assert.Equal(t, 365, date.LengthOfYear(), "year: %d", year)
		}
	})

	t.Run("century years", func(t *testing.T) {
		tests := []struct {
			year     Year
			expected int
		}{
			{1600, 366}, // divisible by 400
			{1700, 365}, // divisible by 100 but not 400
			{1800, 365}, // divisible by 100 but not 400
			{1900, 365}, // divisible by 100 but not 400
			{2000, 366}, // divisible by 400
			{2100, 365}, // divisible by 100 but not 400
			{2400, 366}, // divisible by 400
		}

		for _, tt := range tests {
			date := MustLocalDateOf(tt.year, June, 15)
			assert.Equal(t, tt.expected, date.LengthOfYear(), "year: %d", tt.year)
		}
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalDate
		assert.Equal(t, 0, zero.LengthOfYear())
	})
}
