package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalTime(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		lt, err := LocalTimeOf(0, 0, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, lt.Hour())
		assert.Equal(t, 0, lt.Minute())
		assert.Equal(t, 0, lt.Second())
		assert.Equal(t, 0, lt.Nano())

		lt, err = LocalTimeOf(23, 59, 59, 999999999)
		require.NoError(t, err)
		assert.Equal(t, 23, lt.Hour())
		assert.Equal(t, 59, lt.Minute())
		assert.Equal(t, 59, lt.Second())
		assert.Equal(t, 999999999, lt.Nano())

		lt, err = LocalTimeOf(12, 30, 45, 123456789)
		require.NoError(t, err)
		assert.Equal(t, 12, lt.Hour())
		assert.Equal(t, 30, lt.Minute())
		assert.Equal(t, 45, lt.Second())
		assert.Equal(t, 123456789, lt.Nano())
	})

	t.Run("invalid hour", func(t *testing.T) {
		_, err := LocalTimeOf(24, 0, 0, 0)
		assert.Error(t, err)

		_, err = LocalTimeOf(-1, 0, 0, 0)
		assert.Error(t, err)

		_, err = LocalTimeOf(25, 0, 0, 0)
		assert.Error(t, err)
	})

	t.Run("invalid minute", func(t *testing.T) {
		_, err := LocalTimeOf(0, 60, 0, 0)
		assert.Error(t, err)

		_, err = LocalTimeOf(0, -1, 0, 0)
		assert.Error(t, err)

		_, err = LocalTimeOf(0, 61, 0, 0)
		assert.Error(t, err)
	})

	t.Run("invalid second", func(t *testing.T) {
		_, err := LocalTimeOf(0, 0, 60, 0)
		assert.Error(t, err)

		_, err = LocalTimeOf(0, 0, -1, 0)
		assert.Error(t, err)

		_, err = LocalTimeOf(0, 0, 61, 0)
		assert.Error(t, err)
	})

	t.Run("invalid nanosecond", func(t *testing.T) {
		_, err := LocalTimeOf(0, 0, 0, 1000000000)
		assert.Error(t, err)

		_, err = LocalTimeOf(0, 0, 0, -1)
		assert.Error(t, err)
	})
}

func TestMustNewLocalTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		assert.NotPanics(t, func() {
			lt := MustLocalTimeOf(14, 30, 45, 123456789)
			assert.Equal(t, 14, lt.Hour())
			assert.Equal(t, 30, lt.Minute())
			assert.Equal(t, 45, lt.Second())
			assert.Equal(t, 123456789, lt.Nano())
		})
	})

	t.Run("invalid time panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalTimeOf(24, 0, 0, 0)
		})

		assert.Panics(t, func() {
			MustLocalTimeOf(0, 60, 0, 0)
		})

		assert.Panics(t, func() {
			MustLocalTimeOf(0, 0, 60, 0)
		})

		assert.Panics(t, func() {
			MustLocalTimeOf(0, 0, 0, 1000000000)
		})
	})
}

func TestLocalTime_IsZero(t *testing.T) {
	var zero LocalTime
	assert.True(t, zero.IsZero())

	lt := MustLocalTimeOf(0, 0, 0, 0)
	assert.False(t, lt.IsZero())

	lt = MustLocalTimeOf(12, 30, 45, 0)
	assert.False(t, lt.IsZero())
}

func TestLocalTime_Components(t *testing.T) {
	tests := []struct {
		name        string
		hour        int
		minute      int
		second      int
		nanosecond  int
		millisecond int
	}{
		{"midnight", 0, 0, 0, 0, 0},
		{"noon", 12, 0, 0, 0, 0},
		{"end of day", 23, 59, 59, 999999999, 999},
		{"with milliseconds", 14, 30, 45, 123000000, 123},
		{"with nanoseconds", 9, 15, 30, 123456789, 123},
		{"1 second before midnight", 23, 59, 59, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := MustLocalTimeOf(tt.hour, tt.minute, tt.second, tt.nanosecond)
			assert.Equal(t, tt.hour, lt.Hour(), "Hour")
			assert.Equal(t, tt.minute, lt.Minute(), "Minute")
			assert.Equal(t, tt.second, lt.Second(), "Second")
			assert.Equal(t, tt.nanosecond, lt.Nano(), "Nano")
			assert.Equal(t, tt.millisecond, lt.Millisecond(), "Millisecond")
		})
	}

	var zero LocalTime
	assert.Equal(t, 0, zero.Hour())
	assert.Equal(t, 0, zero.Minute())
	assert.Equal(t, 0, zero.Second())
	assert.Equal(t, 0, zero.Nano())
	assert.Equal(t, 0, zero.Millisecond())
}

func TestLocalTime_Compare(t *testing.T) {
	t1 := MustLocalTimeOf(12, 30, 45, 123456789)
	t2 := MustLocalTimeOf(12, 30, 45, 123456789)
	t3 := MustLocalTimeOf(12, 30, 46, 0)
	t4 := MustLocalTimeOf(12, 31, 0, 0)
	t5 := MustLocalTimeOf(13, 0, 0, 0)
	t6 := MustLocalTimeOf(12, 30, 45, 123456788)

	assert.Equal(t, 0, t1.Compare(t2), "same time")
	assert.Equal(t, -1, t1.Compare(t3), "earlier by second")
	assert.Equal(t, 1, t3.Compare(t1), "later by second")
	assert.Equal(t, -1, t1.Compare(t4), "earlier by minute")
	assert.Equal(t, -1, t1.Compare(t5), "earlier by hour")
	assert.Equal(t, 1, t1.Compare(t6), "later by nanosecond")

	var zero LocalTime
	assert.Equal(t, 0, zero.Compare(LocalTime{}), "zero equals zero")
	assert.Equal(t, -1, zero.Compare(t1), "zero is before non-zero")
	assert.Equal(t, 1, t1.Compare(zero), "non-zero is after zero")
}

func TestLocalTime_IsBefore(t *testing.T) {
	t1 := MustLocalTimeOf(10, 0, 0, 0)
	t2 := MustLocalTimeOf(11, 0, 0, 0)
	t3 := MustLocalTimeOf(10, 0, 0, 0)

	assert.True(t, t1.IsBefore(t2))
	assert.False(t, t2.IsBefore(t1))
	assert.False(t, t1.IsBefore(t3))
}

func TestLocalTime_IsAfter(t *testing.T) {
	t1 := MustLocalTimeOf(10, 0, 0, 0)
	t2 := MustLocalTimeOf(11, 0, 0, 0)
	t3 := MustLocalTimeOf(10, 0, 0, 0)

	assert.False(t, t1.IsAfter(t2))
	assert.True(t, t2.IsAfter(t1))
	assert.False(t, t1.IsAfter(t3))
}

func TestLocalTime_GoTime(t *testing.T) {
	lt := MustLocalTimeOf(14, 30, 45, 123456789)
	goTime := lt.GoTime()

	assert.Equal(t, 14, goTime.Hour())
	assert.Equal(t, 30, goTime.Minute())
	assert.Equal(t, 45, goTime.Second())
	assert.Equal(t, 123456789, goTime.Nanosecond())
	assert.Equal(t, time.UTC, goTime.Location())

	// check that date is set to epoch
	assert.Equal(t, 1970, goTime.Year())
	assert.Equal(t, time.January, goTime.Month())
	assert.Equal(t, 1, goTime.Day())

	var zero LocalTime
	assert.True(t, zero.GoTime().IsZero())
}

func TestNewLocalTimeByGoTime(t *testing.T) {
	goTime := time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
	lt := LocalTimeOfGoTime(goTime)

	assert.Equal(t, 14, lt.Hour())
	assert.Equal(t, 30, lt.Minute())
	assert.Equal(t, 45, lt.Second())
	assert.Equal(t, 123456789, lt.Nano())

	// Test with different time zone (should ignore timezone)
	loc, _ := time.LoadLocation("America/New_York")
	goTime = time.Date(2024, time.March, 15, 14, 30, 45, 123456789, loc)
	lt = LocalTimeOfGoTime(goTime)

	assert.Equal(t, 14, lt.Hour())
	assert.Equal(t, 30, lt.Minute())
	assert.Equal(t, 45, lt.Second())
	assert.Equal(t, 123456789, lt.Nano())

	// Test with zero time
	lt = LocalTimeOfGoTime(time.Time{})
	assert.True(t, lt.IsZero())
}

func TestLocalTime_String(t *testing.T) {
	tests := []struct {
		time     LocalTime
		expected string
	}{
		{MustLocalTimeOf(0, 0, 0, 0), "00:00:00"},
		{MustLocalTimeOf(12, 30, 45, 0), "12:30:45"},
		{MustLocalTimeOf(23, 59, 59, 0), "23:59:59"},
		{MustLocalTimeOf(9, 5, 7, 0), "09:05:07"},
		{MustLocalTimeOf(14, 30, 45, 123000000), "14:30:45.123"},
		{MustLocalTimeOf(14, 30, 45, 123456000), "14:30:45.123456"},
		{MustLocalTimeOf(14, 30, 45, 123456789), "14:30:45.123456789"},
		{MustLocalTimeOf(14, 30, 45, 100000000), "14:30:45.100"},
		{MustLocalTimeOf(14, 30, 45, 120000000), "14:30:45.120"},
		{MustLocalTimeOf(14, 30, 45, 1), "14:30:45.000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.time.String())
		})
	}
}

func TestLocalTime_MarshalText(t *testing.T) {
	lt := MustLocalTimeOf(14, 30, 45, 123456789)
	text, err := lt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "14:30:45.123456789", string(text))

	lt = MustLocalTimeOf(9, 5, 7, 0)
	text, err = lt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "09:05:07", string(text))

	var zero LocalTime
	text, err = zero.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(text))
}

func TestLocalTime_UnmarshalText(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte("14:30:45"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 0), lt)

		err = lt.UnmarshalText([]byte("09:05:07"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(9, 5, 7, 0), lt)

		err = lt.UnmarshalText([]byte("00:00:00"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(0, 0, 0, 0), lt)

		err = lt.UnmarshalText([]byte("23:59:59"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(23, 59, 59, 0), lt)
	})

	t.Run("with fractional seconds", func(t *testing.T) {
		var lt LocalTime

		err := lt.UnmarshalText([]byte("14:30:45.123"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123000000), lt)

		err = lt.UnmarshalText([]byte("14:30:45.123456"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456000), lt)

		err = lt.UnmarshalText([]byte("14:30:45.123456789"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), lt)

		err = lt.UnmarshalText([]byte("14:30:45.1"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 100000000), lt)

		err = lt.UnmarshalText([]byte("14:30:45.000000001"))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 1), lt)
	})

	t.Run("empty string", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("invalid format", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte("14-30-45"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("14:30"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("not-a-time"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("1:2:3"))
		assert.Error(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte("24:00:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("23:60:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("23:59:60"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("25:00:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("12:61:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("12:30:61"))
		assert.Error(t, err)
	})
}

func TestLocalTime_MarshalJSON(t *testing.T) {
	lt := MustLocalTimeOf(14, 30, 45, 123456789)
	data, err := json.Marshal(lt)
	require.NoError(t, err)
	assert.Equal(t, `"14:30:45.123456789"`, string(data))

	lt = MustLocalTimeOf(9, 5, 7, 0)
	data, err = json.Marshal(lt)
	require.NoError(t, err)
	assert.Equal(t, `"09:05:07"`, string(data))

	var zero LocalTime
	data, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `""`, string(data))
}

func TestLocalTime_UnmarshalJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`"14:30:45"`), &lt)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 0), lt)

		err = json.Unmarshal([]byte(`"14:30:45.123456789"`), &lt)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), lt)
	})

	t.Run("empty string", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`""`), &lt)
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("null", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`null`), &lt)
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`"invalid-time"`), &lt)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`"24:00:00"`), &lt)
		assert.Error(t, err)
	})
}

func TestLocalTime_Scan(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		var lt LocalTime
		err := lt.Scan(nil)
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("string value", func(t *testing.T) {
		var lt LocalTime
		err := lt.Scan("14:30:45")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 0), lt)

		err = lt.Scan("14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), lt)
	})

	t.Run("time.LocalTime value", func(t *testing.T) {
		var lt LocalTime
		goTime := time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
		err := lt.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), lt)
	})

	t.Run("unsupported type", func(t *testing.T) {
		var lt LocalTime
		err := lt.Scan(12345)
		assert.Error(t, err)
	})
}

func TestLocalTime_Value(t *testing.T) {
	lt := MustLocalTimeOf(14, 30, 45, 123456789)
	val, err := lt.Value()
	require.NoError(t, err)
	assert.Equal(t, "14:30:45.123456789", val)

	lt = MustLocalTimeOf(9, 5, 7, 0)
	val, err = lt.Value()
	require.NoError(t, err)
	assert.Equal(t, "09:05:07", val)

	var zero LocalTime
	val, err = zero.Value()
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestLocalTime_AppendText(t *testing.T) {
	lt := MustLocalTimeOf(14, 30, 45, 123456789)
	buf := []byte("LocalTime: ")
	buf, err := lt.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalTime: 14:30:45.123456789", string(buf))

	lt = MustLocalTimeOf(9, 5, 7, 0)
	buf = []byte("LocalTime: ")
	buf, err = lt.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalTime: 09:05:07", string(buf))

	var zero LocalTime
	buf = []byte("LocalTime: ")
	buf, err = zero.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalTime: ", string(buf))
}

func TestLocalTime_SpecialCases(t *testing.T) {
	t.Run("midnight", func(t *testing.T) {
		lt := MustLocalTimeOf(0, 0, 0, 0)
		assert.Equal(t, "00:00:00", lt.String())
		assert.Equal(t, 0, lt.Hour())
		assert.Equal(t, 0, lt.Minute())
		assert.Equal(t, 0, lt.Second())
		assert.Equal(t, 0, lt.Nano())
	})

	t.Run("one nanosecond before midnight", func(t *testing.T) {
		lt := MustLocalTimeOf(23, 59, 59, 999999999)
		assert.Equal(t, 23, lt.Hour())
		assert.Equal(t, 59, lt.Minute())
		assert.Equal(t, 59, lt.Second())
		assert.Equal(t, 999999999, lt.Nano())
		assert.Equal(t, 999, lt.Millisecond())
	})

	t.Run("noon", func(t *testing.T) {
		lt := MustLocalTimeOf(12, 0, 0, 0)
		assert.Equal(t, "12:00:00", lt.String())
		assert.Equal(t, 12, lt.Hour())
	})

	t.Run("fractional seconds precision", func(t *testing.T) {
		// Millisecond precision
		lt := MustLocalTimeOf(12, 0, 0, 123000000)
		assert.Equal(t, "12:00:00.123", lt.String())
		assert.Equal(t, 123, lt.Millisecond())

		// Microsecond precision
		lt = MustLocalTimeOf(12, 0, 0, 123456000)
		assert.Equal(t, "12:00:00.123456", lt.String())

		// Nano precision
		lt = MustLocalTimeOf(12, 0, 0, 123456789)
		assert.Equal(t, "12:00:00.123456789", lt.String())

		// Single digit fractional second (aligned to milliseconds)
		lt = MustLocalTimeOf(12, 0, 0, 100000000)
		assert.Equal(t, "12:00:00.100", lt.String())

		// Trailing zeros aligned to 3-digit boundaries
		lt = MustLocalTimeOf(12, 0, 0, 120000000)
		assert.Equal(t, "12:00:00.120", lt.String())

		// More trailing zero examples (aligned to microseconds)
		lt = MustLocalTimeOf(12, 0, 0, 123400000)
		assert.Equal(t, "12:00:00.123400", lt.String())
	})
}

func TestLocalTime_BoundaryValues(t *testing.T) {
	tests := []struct {
		name       string
		hour       int
		minute     int
		second     int
		nanosecond int
	}{
		{"min values", 0, 0, 0, 0},
		{"max values", 23, 59, 59, 999999999},
		{"max hour", 23, 0, 0, 0},
		{"max minute", 0, 59, 0, 0},
		{"max second", 0, 0, 59, 0},
		{"max nanosecond", 0, 0, 0, 999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := MustLocalTimeOf(tt.hour, tt.minute, tt.second, tt.nanosecond)
			assert.Equal(t, tt.hour, lt.Hour())
			assert.Equal(t, tt.minute, lt.Minute())
			assert.Equal(t, tt.second, lt.Second())
			assert.Equal(t, tt.nanosecond, lt.Nano())

			// Round-trip through string
			str := lt.String()
			var lt2 LocalTime
			err := lt2.UnmarshalText([]byte(str))
			require.NoError(t, err)
			assert.Equal(t, lt, lt2)
		})
	}
}

func TestLocalTime_Serialization(t *testing.T) {
	tests := []struct {
		name       string
		time       LocalTime
		textFormat string
	}{
		{"midnight", MustLocalTimeOf(0, 0, 0, 0), "00:00:00"},
		{"noon", MustLocalTimeOf(12, 0, 0, 0), "12:00:00"},
		{"with milliseconds", MustLocalTimeOf(14, 30, 45, 123000000), "14:30:45.123"},
		{"with microseconds", MustLocalTimeOf(14, 30, 45, 123456000), "14:30:45.123456"},
		{"with nanoseconds", MustLocalTimeOf(14, 30, 45, 123456789), "14:30:45.123456789"},
		{"end of day", MustLocalTimeOf(23, 59, 59, 999999999), "23:59:59.999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test String
			assert.Equal(t, tt.textFormat, tt.time.String())

			// Test MarshalText
			text, err := tt.time.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.textFormat, string(text))

			// Test UnmarshalText
			var lt LocalTime
			err = lt.UnmarshalText([]byte(tt.textFormat))
			require.NoError(t, err)
			assert.Equal(t, tt.time, lt)

			// Test MarshalJSON
			jsonData, err := json.Marshal(tt.time)
			require.NoError(t, err)
			assert.Equal(t, `"`+tt.textFormat+`"`, string(jsonData))

			// Test UnmarshalJSON
			var lt2 LocalTime
			err = json.Unmarshal(jsonData, &lt2)
			require.NoError(t, err)
			assert.Equal(t, tt.time, lt2)
		})
	}
}

func TestLocalTime_CompareConsistency(t *testing.T) {
	times := []LocalTime{
		MustLocalTimeOf(0, 0, 0, 0),
		MustLocalTimeOf(6, 0, 0, 0),
		MustLocalTimeOf(12, 0, 0, 0),
		MustLocalTimeOf(12, 30, 0, 0),
		MustLocalTimeOf(12, 30, 30, 0),
		MustLocalTimeOf(12, 30, 30, 500000000),
		MustLocalTimeOf(18, 0, 0, 0),
		MustLocalTimeOf(23, 59, 59, 999999999),
	}

	// Test that times are ordered correctly
	for i := 0; i < len(times)-1; i++ {
		assert.True(t, times[i].IsBefore(times[i+1]), "times[%d] should be before times[%d]", i, i+1)
		assert.False(t, times[i].IsAfter(times[i+1]), "times[%d] should not be after times[%d]", i, i+1)
		assert.Equal(t, -1, times[i].Compare(times[i+1]), "times[%d] should compare as -1 to times[%d]", i, i+1)
	}

	// Test equality
	for i, lt := range times {
		assert.Equal(t, 0, lt.Compare(lt), "time should equal itself")
		assert.False(t, lt.IsBefore(lt), "time should not be before itself")
		assert.False(t, lt.IsAfter(lt), "time should not be after itself")

		// Create copy and test
		copy := MustLocalTimeOf(lt.Hour(), lt.Minute(), lt.Second(), lt.Nano())
		assert.Equal(t, 0, lt.Compare(copy), "times[%d] should equal its copy", i)
	}
}

func TestLocalTimeNow(t *testing.T) {
	// Test that LocalTimeNow() returns a valid time
	now := LocalTimeNow()
	assert.False(t, now.IsZero(), "LocalTimeNow should not be zero")

	// Test that it's reasonable (between midnight and end of day)
	assert.True(t, now.Hour() >= 0 && now.Hour() < 24, "Hour should be valid")
	assert.True(t, now.Minute() >= 0 && now.Minute() < 60, "Minute should be valid")
	assert.True(t, now.Second() >= 0 && now.Second() < 60, "Second should be valid")
}

func TestLocalTimeNowUTC(t *testing.T) {
	nowUTC := LocalTimeNowUTC()
	assert.False(t, nowUTC.IsZero(), "LocalTimeNowUTC should not be zero")

	// Test that it's reasonable
	assert.True(t, nowUTC.Hour() >= 0 && nowUTC.Hour() < 24, "Hour should be valid")
	assert.True(t, nowUTC.Minute() >= 0 && nowUTC.Minute() < 60, "Minute should be valid")
	assert.True(t, nowUTC.Second() >= 0 && nowUTC.Second() < 60, "Second should be valid")
}

func TestParseLocalTime(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		time, err := LocalTimeParse("14:30:45")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 0), time)

		time, err = LocalTimeParse("09:05:07")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(9, 5, 7, 0), time)

		time, err = LocalTimeParse("00:00:00")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(0, 0, 0, 0), time)

		time, err = LocalTimeParse("23:59:59")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(23, 59, 59, 0), time)
	})

	t.Run("with fractional seconds", func(t *testing.T) {
		time, err := LocalTimeParse("14:30:45.123")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123000000), time)

		time, err = LocalTimeParse("14:30:45.123456")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456000), time)

		time, err = LocalTimeParse("14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), time)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := LocalTimeParse("14-30-45")
		assert.Error(t, err)

		_, err = LocalTimeParse("14:30")
		assert.Error(t, err)

		_, err = LocalTimeParse("not-a-time")
		assert.Error(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		_, err := LocalTimeParse("24:00:00")
		assert.Error(t, err)

		_, err = LocalTimeParse("23:60:00")
		assert.Error(t, err)

		_, err = LocalTimeParse("23:59:60")
		assert.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		time, err := LocalTimeParse("")
		require.NoError(t, err)
		assert.True(t, time.IsZero())
	})
}

func TestMustParseLocalTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		assert.NotPanics(t, func() {
			time := MustLocalTimeParse("14:30:45.123456789")
			assert.Equal(t, MustLocalTimeOf(14, 30, 45, 123456789), time)
		})
	})

	t.Run("invalid time panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalTimeParse("24:00:00")
		})

		assert.Panics(t, func() {
			MustLocalTimeParse("invalid")
		})
	})
}

func TestLocalTimeNowIn(t *testing.T) {
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
			nowIn := LocalTimeNowIn(tt.loc)
			assert.False(t, nowIn.IsZero(), "LocalTimeNowIn should not be zero")

			// Test that it's reasonable
			assert.True(t, nowIn.Hour() >= 0 && nowIn.Hour() < 24, "Hour should be valid")
			assert.True(t, nowIn.Minute() >= 0 && nowIn.Minute() < 60, "Minute should be valid")
			assert.True(t, nowIn.Second() >= 0 && nowIn.Second() < 60, "Second should be valid")
		})
	}
}

func TestLocalTime_ValuePostgres(t *testing.T) {
	var pg = GetPG(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustLocalTimeParse("12:00:00")
		var actual LocalTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT $1::time without time zone, $1::time without time zone = '12:00:00'", expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT NULL::time without time zone, $1::time without time zone is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalTime_ValueMySQL(t *testing.T) {
	var mysql = GetMySQL(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustLocalTimeParse("08:00:00")
		var actual LocalTime
		var expectedTrue bool
		// I don't want to understand why MySQL has this bug
		var e = mysql.QueryRow("SELECT cast(cast(? as char) as time), cast(cast(? as char) as time) = '08:00:00'", expected, expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		t.Log(expected.Value())
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalTime
		var expectedTrue bool
		// I don't want to understand why MySQL has this bug
		var e = mysql.QueryRow("SELECT cast(cast(? as char) as time), cast(cast(? as char) as time) is null", actual, actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalTimeOfNanoOfDay(t *testing.T) {
	t.Run("valid nano of day", func(t *testing.T) {
		// Midnight (0 nanoseconds)
		lt, err := LocalTimeOfNanoOfDay(0)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(0, 0, 0, 0), lt)

		// 1 hour = 3,600,000,000,000 nanoseconds
		lt, err = LocalTimeOfNanoOfDay(int64(time.Hour))
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(1, 0, 0, 0), lt)

		// 12:30:45.123456789
		nanos := int64(12)*int64(time.Hour) +
			int64(30)*int64(time.Minute) +
			int64(45)*int64(time.Second) +
			123456789
		lt, err = LocalTimeOfNanoOfDay(nanos)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(12, 30, 45, 123456789), lt)

		// Last nanosecond of the day (23:59:59.999999999)
		maxNanos := int64(24)*int64(time.Hour) - 1
		lt, err = LocalTimeOfNanoOfDay(maxNanos)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(23, 59, 59, 999999999), lt)
	})

	t.Run("invalid nano of day", func(t *testing.T) {
		// Negative value
		_, err := LocalTimeOfNanoOfDay(-1)
		assert.Error(t, err)

		// value >= 24 hours in nanoseconds
		_, err = LocalTimeOfNanoOfDay(24 * int64(time.Hour))
		assert.Error(t, err)

		// Large positive value
		_, err = LocalTimeOfNanoOfDay(100 * int64(time.Hour))
		assert.Error(t, err)
	})

	t.Run("round trip with GetField", func(t *testing.T) {
		// Create time from components
		original := MustLocalTimeOf(15, 45, 30, 987654321)

		// Get nano of day
		nanoOfDay := original.GetField(FieldNanoOfDay).Int64()

		// Create new time from nano of day
		reconstructed, err := LocalTimeOfNanoOfDay(nanoOfDay)
		require.NoError(t, err)

		// Should be equal
		assert.Equal(t, original, reconstructed)
	})
}

func TestMustLocalTimeOfNanoOfDay(t *testing.T) {
	t.Run("valid value", func(t *testing.T) {
		assert.NotPanics(t, func() {
			lt := MustLocalTimeOfNanoOfDay(int64(12 * time.Hour))
			assert.Equal(t, MustLocalTimeOf(12, 0, 0, 0), lt)
		})
	})

	t.Run("invalid value panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalTimeOfNanoOfDay(-1)
		})

		assert.Panics(t, func() {
			MustLocalTimeOfNanoOfDay(24 * int64(time.Hour))
		})
	})
}

func TestLocalTimeOfSecondOfDay(t *testing.T) {
	t.Run("valid second of day", func(t *testing.T) {
		// Midnight (0 seconds)
		lt, err := LocalTimeOfSecondOfDay(0)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(0, 0, 0, 0), lt)

		// 1 hour = 3,600 seconds
		lt, err = LocalTimeOfSecondOfDay(3600)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(1, 0, 0, 0), lt)

		// 12:30:45 = 45,045 seconds
		seconds := 12*3600 + 30*60 + 45
		lt, err = LocalTimeOfSecondOfDay(seconds)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(12, 30, 45, 0), lt)

		// Last second of the day (23:59:59 = 86,399 seconds)
		lt, err = LocalTimeOfSecondOfDay(86399)
		require.NoError(t, err)
		assert.Equal(t, MustLocalTimeOf(23, 59, 59, 0), lt)
	})

	t.Run("invalid second of day", func(t *testing.T) {
		// Negative value
		_, err := LocalTimeOfSecondOfDay(-1)
		assert.Error(t, err)

		// value >= 24 hours in seconds (86,400)
		_, err = LocalTimeOfSecondOfDay(86400)
		assert.Error(t, err)

		// Large positive value
		_, err = LocalTimeOfSecondOfDay(100000)
		assert.Error(t, err)
	})

	t.Run("round trip with GetField", func(t *testing.T) {
		// Create time from components (no nanoseconds)
		original := MustLocalTimeOf(15, 45, 30, 0)

		// Get second of day
		secondOfDay := int(original.GetField(FieldSecondOfDay).Int64())

		// Create new time from second of day
		reconstructed, err := LocalTimeOfSecondOfDay(secondOfDay)
		require.NoError(t, err)

		// Should be equal
		assert.Equal(t, original, reconstructed)
	})

	t.Run("nanoseconds are zero", func(t *testing.T) {
		// Create time from second of day
		lt, err := LocalTimeOfSecondOfDay(12345)
		require.NoError(t, err)

		// Nanoseconds should be 0
		assert.Equal(t, 0, lt.Nano())
	})
}

func TestMustLocalTimeOfSecondOfDay(t *testing.T) {
	t.Run("valid value", func(t *testing.T) {
		assert.NotPanics(t, func() {
			lt := MustLocalTimeOfSecondOfDay(43200) // 12:00:00
			assert.Equal(t, MustLocalTimeOf(12, 0, 0, 0), lt)
		})
	})

	t.Run("invalid value panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLocalTimeOfSecondOfDay(-1)
		})

		assert.Panics(t, func() {
			MustLocalTimeOfSecondOfDay(86400)
		})
	})
}

func TestLocalTimeOfDay_Consistency(t *testing.T) {
	t.Run("second of day and nano of day consistency", func(t *testing.T) {
		// For same second value, both constructors should create same time (ignoring nanos)
		secondOfDay := 45045 // 12:30:45

		lt1, err := LocalTimeOfSecondOfDay(secondOfDay)
		require.NoError(t, err)

		lt2, err := LocalTimeOfNanoOfDay(int64(secondOfDay) * int64(time.Second))
		require.NoError(t, err)

		assert.Equal(t, lt1, lt2)
	})

	t.Run("all valid times", func(t *testing.T) {
		// Test several times throughout the day
		testCases := []struct {
			hour   int
			minute int
			second int
			nano   int
		}{
			{0, 0, 0, 0},            // midnight
			{6, 30, 15, 500000000},  // morning
			{12, 0, 0, 0},           // noon
			{18, 45, 30, 123456789}, // evening
			{23, 59, 59, 999999999}, // last nanosecond
		}

		for _, tc := range testCases {
			original := MustLocalTimeOf(tc.hour, tc.minute, tc.second, tc.nano)

			// Test nano of day round trip
			nanoOfDay := original.GetField(FieldNanoOfDay).Int64()
			fromNano := MustLocalTimeOfNanoOfDay(nanoOfDay)
			assert.Equal(t, original, fromNano, "nano of day round trip failed for %02d:%02d:%02d.%09d", tc.hour, tc.minute, tc.second, tc.nano)

			// Test second of day round trip (only if no nanoseconds)
			if tc.nano == 0 {
				secondOfDay := int(original.GetField(FieldSecondOfDay).Int64())
				fromSecond := MustLocalTimeOfSecondOfDay(secondOfDay)
				assert.Equal(t, original, fromSecond, "second of day round trip failed for %02d:%02d:%02d", tc.hour, tc.minute, tc.second)
			}
		}
	})
}

func TestLocalTime_WithTemporal(t *testing.T) {
	nanoOfDay := int64(1*time.Hour + 2*time.Minute + 3*time.Second + 4)
	base := MustLocalTimeOf(15, 20, 30, 123456789)

	tests := []struct {
		name  string
		base  LocalTime
		field Field
		value int64
		want  LocalTime
	}{
		{
			name:  "nano_of_second",
			base:  base,
			field: FieldNanoOfSecond,
			value: 987654321,
			want:  MustLocalTimeOf(15, 20, 30, 987654321),
		},
		{
			name:  "nano_of_day",
			base:  base,
			field: FieldNanoOfDay,
			value: nanoOfDay,
			want:  MustLocalTimeOf(1, 2, 3, 4),
		},
		{
			name:  "micro_of_second",
			base:  base,
			field: FieldMicroOfSecond,
			value: 456789,
			want:  MustLocalTimeOf(15, 20, 30, 456789000),
		},
		{
			name:  "micro_of_day",
			base:  base,
			field: FieldMicroOfDay,
			value: int64((2*time.Hour + 3*time.Minute + 4*time.Second + 500*time.Microsecond) / time.Microsecond),
			want:  MustLocalTimeOf(2, 3, 4, int(500*time.Microsecond)),
		},
		{
			name:  "milli_of_second",
			base:  base,
			field: FieldMilliOfSecond,
			value: 777,
			want:  MustLocalTimeOf(15, 20, 30, 777000000),
		},
		{
			name:  "milli_of_day",
			base:  base,
			field: FieldMilliOfDay,
			value: int64((5*time.Hour + 6*time.Minute + 7*time.Second + 8*time.Millisecond) / time.Millisecond),
			want:  MustLocalTimeOf(5, 6, 7, 8000000),
		},
		{
			name:  "second_of_minute",
			base:  base,
			field: FieldSecondOfMinute,
			value: 45,
			want:  MustLocalTimeOf(15, 20, 45, 123456789),
		},
		{
			name:  "second_of_day",
			base:  base,
			field: FieldSecondOfDay,
			value: 3661,
			want:  MustLocalTimeOf(1, 1, 1, 123456789),
		},
		{
			name:  "minute_of_hour",
			base:  base,
			field: FieldMinuteOfHour,
			value: 5,
			want:  MustLocalTimeOf(15, 5, 30, 123456789),
		},
		{
			name:  "minute_of_day",
			base:  base,
			field: FieldMinuteOfDay,
			value: 70,
			want:  MustLocalTimeOf(1, 10, 30, 123456789),
		},
		{
			name:  "hour_of_ampm",
			base:  MustLocalTimeOf(19, 10, 20, 123000000),
			field: FieldHourOfAmPm,
			value: 3,
			want:  MustLocalTimeOf(15, 10, 20, 123000000),
		},
		{
			name:  "clock_hour_of_ampm",
			base:  MustLocalTimeOf(3, 10, 20, 0),
			field: FieldClockHourOfAmPm,
			value: 12,
			want:  MustLocalTimeOf(0, 10, 20, 0),
		},
		{
			name:  "hour_of_day",
			base:  base,
			field: FieldHourOfDay,
			value: 22,
			want:  MustLocalTimeOf(22, 20, 30, 123456789),
		},
		{
			name:  "clock_hour_of_day",
			base:  MustLocalTimeOf(10, 5, 6, 700000000),
			field: FieldClockHourOfDay,
			value: 24,
			want:  MustLocalTimeOf(0, 5, 6, 700000000),
		},
		{
			name:  "ampm_of_day",
			base:  MustLocalTimeOf(9, 15, 30, 111000000),
			field: FieldAmPmOfDay,
			value: 1,
			want:  MustLocalTimeOf(21, 15, 30, 111000000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.base.Chain().WithField(tt.field, TemporalValue{v: tt.value})
			require.NoError(t, r.eError)
			assert.Equal(t, tt.want, r.MustGet(), "want %s got %s", tt.want, r.MustGet())
		})
	}
}
