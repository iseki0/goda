package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOffsetDateTime(t *testing.T) {
	t.Run("valid components", func(t *testing.T) {
		offset := MustZoneOffsetOfHours(1)
		odt, err := OffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, March, odt.Month())
		assert.Equal(t, 15, odt.DayOfMonth())
		assert.Equal(t, 14, odt.Hour())
		assert.Equal(t, 30, odt.Minute())
		assert.Equal(t, 45, odt.Second())
		assert.Equal(t, 123456789, odt.Nanosecond())
		assert.Equal(t, 3600, odt.Offset().TotalSeconds())
	})

	t.Run("invalid date", func(t *testing.T) {
		offset := ZoneOffsetUTC()
		_, err := OffsetDateTimeOf(2024, February, 30, 14, 30, 45, 0, offset)
		assert.Error(t, err)
	})

	t.Run("invalid time", func(t *testing.T) {
		offset := ZoneOffsetUTC()
		_, err := OffsetDateTimeOf(2024, March, 15, 25, 30, 45, 0, offset)
		assert.Error(t, err)
	})
}

func TestMustNewOffsetDateTime(t *testing.T) {
	t.Run("valid components", func(t *testing.T) {
		offset := MustZoneOffsetOfHours(-5)
		assert.NotPanics(t, func() {
			odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
			assert.Equal(t, Year(2024), odt.Year())
		})
	})

	t.Run("invalid components panic", func(t *testing.T) {
		offset := ZoneOffsetUTC()
		assert.Panics(t, func() {
			MustOffsetDateTimeOf(2024, March, 32, 14, 30, 45, 0, offset)
		})
	})
}

func TestOffsetDateTimeOfGoTime(t *testing.T) {
	t.Run("with offset", func(t *testing.T) {
		loc := time.FixedZone("TEST", 3600)
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, loc)
		odt := OffsetDateTimeOfGoTime(goTime)

		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, March, odt.Month())
		assert.Equal(t, 15, odt.DayOfMonth())
		assert.Equal(t, 14, odt.Hour())
		assert.Equal(t, 30, odt.Minute())
		assert.Equal(t, 45, odt.Second())
		assert.Equal(t, 123456789, odt.Nanosecond())
		assert.Equal(t, 3600, odt.Offset().TotalSeconds())
	})

	t.Run("UTC", func(t *testing.T) {
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
		odt := OffsetDateTimeOfGoTime(goTime)
		assert.Equal(t, 0, odt.Offset().TotalSeconds())
	})

	t.Run("zero time", func(t *testing.T) {
		odt := OffsetDateTimeOfGoTime(time.Time{})
		assert.True(t, odt.IsZero())
	})
}

func TestOffsetDateTimeNow(t *testing.T) {
	now := OffsetDateTimeNow()
	assert.False(t, now.IsZero())
	assert.True(t, now.Year() >= 1970 && now.Year() <= 9999)
}

func TestOffsetDateTimeNowUTC(t *testing.T) {
	nowUTC := OffsetDateTimeNowUTC()
	assert.False(t, nowUTC.IsZero())
	assert.Equal(t, 0, nowUTC.Offset().TotalSeconds())
}

func TestParseOffsetDateTime(t *testing.T) {
	t.Run("with positive offset", func(t *testing.T) {
		odt, err := OffsetDateTimeParse("2024-03-15T14:30:45.123456789+01:00")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, 14, odt.Hour())
		assert.Equal(t, 3600, odt.Offset().TotalSeconds())
	})

	t.Run("with negative offset", func(t *testing.T) {
		odt, err := OffsetDateTimeParse("2024-03-15T14:30:45-05:00")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, -18000, odt.Offset().TotalSeconds())
	})

	t.Run("with Z (UTC)", func(t *testing.T) {
		odt, err := OffsetDateTimeParse("2024-03-15T14:30:45Z")
		require.NoError(t, err)
		assert.Equal(t, 0, odt.Offset().TotalSeconds())
	})

	t.Run("with minutes offset", func(t *testing.T) {
		odt, err := OffsetDateTimeParse("2024-03-15T14:30:45+05:30")
		require.NoError(t, err)
		assert.Equal(t, 19800, odt.Offset().TotalSeconds())
	})

	t.Run("empty string", func(t *testing.T) {
		odt, err := OffsetDateTimeParse("")
		require.NoError(t, err)
		assert.True(t, odt.IsZero())
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := OffsetDateTimeParse("invalid")
		assert.Error(t, err)
	})

	t.Run("missing offset", func(t *testing.T) {
		_, err := OffsetDateTimeParse("2024-03-15T14:30:45")
		assert.Error(t, err)
	})
}

func TestMustParseOffsetDateTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			odt := MustOffsetDateTimeParse("2024-03-15T14:30:45+01:00")
			assert.Equal(t, Year(2024), odt.Year())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustOffsetDateTimeParse("invalid")
		})
	})
}

func TestOffsetDateTime_Accessors(t *testing.T) {
	offset := MustZoneOffsetOfHours(2)
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)

	assert.Equal(t, Year(2024), odt.Year())
	assert.Equal(t, March, odt.Month())
	assert.Equal(t, 15, odt.DayOfMonth())
	assert.Equal(t, Friday, odt.DayOfWeek())
	assert.Equal(t, 75, odt.DayOfYear())
	assert.Equal(t, 14, odt.Hour())
	assert.Equal(t, 30, odt.Minute())
	assert.Equal(t, 45, odt.Second())
	assert.Equal(t, 123, odt.Millisecond())
	assert.Equal(t, 123456789, odt.Nanosecond())
	assert.Equal(t, 7200, odt.Offset().TotalSeconds())

	assert.Equal(t, MustLocalDateOf(2024, March, 15), odt.LocalDate())
	assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), odt.LocalTime())
	assert.Equal(t, MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789), odt.LocalDateTime())
}

func TestOffsetDateTime_IsZero(t *testing.T) {
	var zero OffsetDateTime
	assert.True(t, zero.IsZero())

	offset := ZoneOffsetUTC()
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
	assert.False(t, odt.IsZero())
}

func TestOffsetDateTime_IsLeapYear(t *testing.T) {
	offset := ZoneOffsetUTC()
	odt2024 := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
	odt2023 := MustOffsetDateTimeOf(2023, March, 15, 14, 30, 45, 0, offset)

	assert.True(t, odt2024.IsLeapYear())
	assert.False(t, odt2023.IsLeapYear())
}

func TestOffsetDateTime_ToEpochSecond(t *testing.T) {
	t.Run("UTC", func(t *testing.T) {
		offset := ZoneOffsetUTC()
		// 2024-03-15T14:30:45Z
		odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
		epochSec := odt.EpochSecond()

		// Verify by converting to Go time
		goTime := odt.GoTime()
		assert.Equal(t, goTime.Unix(), epochSec)
	})

	t.Run("with positive offset", func(t *testing.T) {
		offset := MustZoneOffsetOfHours(1)
		// 2024-03-15T14:30:45+01:00
		odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
		epochSec := odt.EpochSecond()

		// Should be same as 13:30:45 UTC
		offsetUTC := ZoneOffsetUTC()
		odtUTC := MustOffsetDateTimeOf(2024, March, 15, 13, 30, 45, 0, offsetUTC)
		assert.Equal(t, odtUTC.EpochSecond(), epochSec)
	})

	t.Run("zero value", func(t *testing.T) {
		var odt OffsetDateTime
		assert.Equal(t, int64(0), odt.EpochSecond())
	})
}

func TestOffsetDateTime_GoTime(t *testing.T) {
	offset := MustZoneOffsetOfHours(5)
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)

	goTime := odt.GoTime()
	assert.Equal(t, 2024, goTime.Year())
	assert.Equal(t, time.Month(3), goTime.Month())
	assert.Equal(t, 15, goTime.Day())
	assert.Equal(t, 14, goTime.Hour())
	assert.Equal(t, 30, goTime.Minute())
	assert.Equal(t, 45, goTime.Second())
	assert.Equal(t, 123456789, goTime.Nanosecond())

	_, offsetSec := goTime.Zone()
	assert.Equal(t, 5*3600, offsetSec)
}

func TestOffsetDateTime_Compare(t *testing.T) {
	offset1 := MustZoneOffsetOfHours(1)
	offset2 := MustZoneOffsetOfHours(5)

	odt1 := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset1)
	odt2 := MustOffsetDateTimeOf(2024, March, 15, 18, 30, 45, 0, offset2) // Same instant
	odt3 := MustOffsetDateTimeOf(2024, March, 15, 19, 30, 45, 0, offset2) // Later

	assert.Equal(t, -1, odt1.Compare(odt2)) // Same instant
	assert.Equal(t, -1, odt1.Compare(odt3))
	assert.Equal(t, 1, odt3.Compare(odt1))

	// Test with zero values
	var zero OffsetDateTime
	assert.Equal(t, -1, zero.Compare(odt1))
	assert.Equal(t, 1, odt1.Compare(zero))
	assert.Equal(t, 0, zero.Compare(zero))
}

func TestOffsetDateTime_IsBefore(t *testing.T) {
	offset := ZoneOffsetUTC()
	odt1 := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
	odt2 := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 46, 0, offset)

	assert.True(t, odt1.IsBefore(odt2))
	assert.False(t, odt2.IsBefore(odt1))
	assert.False(t, odt1.IsBefore(odt1))
}

func TestOffsetDateTime_IsAfter(t *testing.T) {
	offset := ZoneOffsetUTC()
	odt1 := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)
	odt2 := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 46, 0, offset)

	assert.True(t, odt2.IsAfter(odt1))
	assert.False(t, odt1.IsAfter(odt2))
	assert.False(t, odt1.IsAfter(odt1))
}

func TestOffsetDateTime_IsSupportedField(t *testing.T) {
	offset := ZoneOffsetUTC()
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offset)

	// Date and time fields should be supported
	assert.True(t, odt.IsSupportedField(FieldYear))
	assert.True(t, odt.IsSupportedField(FieldMonthOfYear))
	assert.True(t, odt.IsSupportedField(FieldDayOfMonth))
	assert.True(t, odt.IsSupportedField(FieldHourOfDay))
	assert.True(t, odt.IsSupportedField(FieldMinuteOfHour))
	assert.True(t, odt.IsSupportedField(FieldSecondOfMinute))

	// Offset-specific fields
	assert.True(t, odt.IsSupportedField(FieldOffsetSeconds))
	assert.True(t, odt.IsSupportedField(FieldInstantSeconds))
}

func TestOffsetDateTime_GetField(t *testing.T) {
	offset := MustZoneOffsetOfHours(2)
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)

	// Date fields
	assert.Equal(t, int64(2024), odt.GetField(FieldYear).Int64())
	assert.Equal(t, int64(3), odt.GetField(FieldMonthOfYear).Int64())
	assert.Equal(t, int64(15), odt.GetField(FieldDayOfMonth).Int64())

	// Time fields
	assert.Equal(t, int64(14), odt.GetField(FieldHourOfDay).Int64())
	assert.Equal(t, int64(30), odt.GetField(FieldMinuteOfHour).Int64())
	assert.Equal(t, int64(45), odt.GetField(FieldSecondOfMinute).Int64())

	// Offset field
	assert.Equal(t, int64(7200), odt.GetField(FieldOffsetSeconds).Int64())

	// Instant field
	assert.True(t, odt.GetField(FieldInstantSeconds).Valid())

	// Zero value
	var zero OffsetDateTime
	assert.True(t, zero.GetField(FieldYear).Unsupported())
}

func TestOffsetDateTime_String(t *testing.T) {
	offset := MustZoneOffsetOfHours(1)
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)
	assert.Equal(t, "2024-03-15T14:30:45.123456789+01:00", odt.String())

	// With UTC
	offsetUTC := ZoneOffsetUTC()
	odtUTC := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 0, offsetUTC)
	assert.Equal(t, "2024-03-15T14:30:45Z", odtUTC.String())

	// Zero value
	var zero OffsetDateTime
	assert.Equal(t, "", zero.String())
}

func TestOffsetDateTime_MarshalText(t *testing.T) {
	offset := MustZoneOffsetOfHours(5)
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)

	text, err := odt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15T14:30:45.123456789+05:00", string(text))
}

func TestOffsetDateTime_UnmarshalText(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var odt OffsetDateTime
		err := odt.UnmarshalText([]byte("2024-03-15T14:30:45.123456789+01:00"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, 3600, odt.Offset().TotalSeconds())
	})

	t.Run("empty", func(t *testing.T) {
		var odt OffsetDateTime
		err := odt.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, odt.IsZero())
	})

	t.Run("invalid", func(t *testing.T) {
		var odt OffsetDateTime
		err := odt.UnmarshalText([]byte("invalid"))
		assert.Error(t, err)
	})
}

func TestOffsetDateTime_MarshalJSON(t *testing.T) {
	offset := MustZoneOffsetOfHours(-5)
	odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)

	data, err := json.Marshal(odt)
	require.NoError(t, err)
	assert.Equal(t, `"2024-03-15T14:30:45.123456789-05:00"`, string(data))
}

func TestOffsetDateTime_UnmarshalJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var odt OffsetDateTime
		err := json.Unmarshal([]byte(`"2024-03-15T14:30:45.123456789+01:00"`), &odt)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
	})

	t.Run("null", func(t *testing.T) {
		var odt OffsetDateTime
		err := json.Unmarshal([]byte("null"), &odt)
		require.NoError(t, err)
		assert.True(t, odt.IsZero())
	})
}

func TestOffsetDateTime_Scan(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		var odt OffsetDateTime
		err := odt.Scan("2024-03-15T14:30:45.123456789+01:00")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, 14, odt.Hour())
		assert.Equal(t, 3600, odt.Offset().TotalSeconds())
	})

	t.Run("from bytes", func(t *testing.T) {
		var odt OffsetDateTime
		err := odt.Scan([]byte("2024-03-15T14:30:45Z"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, 0, odt.Offset().TotalSeconds())
	})

	t.Run("from time.Time", func(t *testing.T) {
		var odt OffsetDateTime
		loc := time.FixedZone("TEST", 3600)
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, loc)
		err := odt.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), odt.Year())
		assert.Equal(t, 14, odt.Hour())
		assert.Equal(t, 3600, odt.Offset().TotalSeconds())
	})

	t.Run("from nil", func(t *testing.T) {
		var odt OffsetDateTime
		err := odt.Scan(nil)
		require.NoError(t, err)
		assert.True(t, odt.IsZero())
	})
}

func TestOffsetDateTime_Value(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		offset := MustZoneOffsetOfHours(1)
		odt := MustOffsetDateTimeOf(2024, March, 15, 14, 30, 45, 123456789, offset)
		val, err := odt.Value()
		require.NoError(t, err)
		assert.Equal(t, "2024-03-15T14:30:45.123456789+01:00", val)
	})

	t.Run("zero", func(t *testing.T) {
		var odt OffsetDateTime
		val, err := odt.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestOffsetDateTime_ValuePostgres(t *testing.T) {
	var pg = GetPG(t)

	t.Run("normal with offset", func(t *testing.T) {
		expected := MustOffsetDateTimeParse("2000-12-29T12:00:00+01:00")
		var actual OffsetDateTime
		var expectedTrue bool
		// PostgreSQL converts TIMESTAMPTZ to the session time zone
		// We'll test round-trip with string representation
		err := pg.QueryRow("SELECT $1::text, $1::text = '2000-12-29T12:00:00+01:00'", expected.String()).Scan(&actual, &expectedTrue)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})

	t.Run("UTC", func(t *testing.T) {
		expected := MustOffsetDateTimeParse("2000-12-29T12:00:00Z")
		var actual OffsetDateTime
		var expectedTrue bool
		err := pg.QueryRow("SELECT $1::text, $1::text = '2000-12-29T12:00:00Z'", expected.String()).Scan(&actual, &expectedTrue)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})

	t.Run("null_value", func(t *testing.T) {
		var actual OffsetDateTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT NULL::text, $1::text is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})

	t.Run("with timestamptz", func(t *testing.T) {
		// Test with actual TIMESTAMPTZ type
		expected := MustOffsetDateTimeParse("2000-12-29T12:00:00+01:00")
		var actual OffsetDateTime
		// When scanning from TIMESTAMPTZ, Go's time.Time will be in the session timezone
		err := pg.QueryRow("SELECT $1::timestamptz", expected.GoTime()).Scan(&actual)
		assert.NoError(t, err)
		// The instant should be the same even if the offset differs
		assert.Equal(t, expected.EpochSecond(), actual.EpochSecond())
	})
}

func TestOffsetDateTime_ValueMySQL(t *testing.T) {
	var mysql = GetMySQL(t)

	t.Run("normal with offset", func(t *testing.T) {
		expected := MustOffsetDateTimeParse("2000-12-29T12:00:00+01:00")
		var actual OffsetDateTime
		var expectedTrue bool
		// MySQL doesn't have native timezone support in DATETIME, so we test with string
		err := mysql.QueryRow("SELECT ?, ? = '2000-12-29T12:00:00+01:00'", expected.String(), expected.String()).Scan(&actual, &expectedTrue)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})

	t.Run("UTC", func(t *testing.T) {
		expected := MustOffsetDateTimeParse("2000-12-29T12:00:00Z")
		var actual OffsetDateTime
		var expectedTrue bool
		err := mysql.QueryRow("SELECT ?, ? = '2000-12-29T12:00:00Z'", expected.String(), expected.String()).Scan(&actual, &expectedTrue)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})

	t.Run("null_value", func(t *testing.T) {
		var actual OffsetDateTime
		var expectedTrue bool
		var e = mysql.QueryRow("SELECT CAST(NULL AS CHAR), ? is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})

	t.Run("round_trip_instant", func(t *testing.T) {
		// Test that we can store and retrieve the instant correctly
		expected := MustOffsetDateTimeParse("2000-12-29T12:00:00+05:30")
		var actual OffsetDateTime
		err := mysql.QueryRow("SELECT ?", expected.String()).Scan(&actual)
		assert.NoError(t, err)
		assert.Equal(t, expected.EpochSecond(), actual.EpochSecond())
		assert.Equal(t, expected, actual)
	})
}
