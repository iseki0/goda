package goda

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZoneOffsetOfSeconds(t *testing.T) {
	t.Run("valid offsets", func(t *testing.T) {
		tests := []struct {
			seconds int
			hours   int
			minutes int
			secs    int
		}{
			{0, 0, 0, 0},
			{3600, 1, 0, 0},
			{-3600, -1, 0, 0},
			{7200, 2, 0, 0},
			{5400, 1, 30, 0},
			{-5400, -1, -30, 0},
			{19800, 5, 30, 0},   // +05:30 (India)
			{34200, 9, 30, 0},   // +09:30 (Australia)
			{-18000, -5, 0, 0},  // -05:00 (US Eastern)
			{3661, 1, 1, 1},     // +01:01:01
			{-3661, -1, -1, -1}, // -01:01:01
			{64800, 18, 0, 0},   // Max: +18:00
			{-64800, -18, 0, 0}, // Min: -18:00
		}

		for _, tt := range tests {
			z, err := ZoneOffsetOfSeconds(tt.seconds)
			require.NoError(t, err)
			assert.Equal(t, tt.seconds, z.TotalSeconds())
			assert.Equal(t, tt.hours, z.Hours())
			assert.Equal(t, tt.minutes, z.Minutes())
			assert.Equal(t, tt.secs, z.Seconds())
		}
	})

	t.Run("out of range", func(t *testing.T) {
		_, err := ZoneOffsetOfSeconds(64801)
		assert.Error(t, err)

		_, err = ZoneOffsetOfSeconds(-64801)
		assert.Error(t, err)

		_, err = ZoneOffsetOfSeconds(100000)
		assert.Error(t, err)
	})
}

func TestMustZoneOffsetOfSeconds(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			z := MustZoneOffsetOfSeconds(3600)
			assert.Equal(t, 3600, z.TotalSeconds())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustZoneOffsetOfSeconds(100000)
		})
	})
}

func TestZoneOffsetOf(t *testing.T) {
	t.Run("valid offsets", func(t *testing.T) {
		tests := []struct {
			hours   int
			minutes int
			seconds int
			total   int
		}{
			{0, 0, 0, 0},
			{1, 0, 0, 3600},
			{-1, 0, 0, -3600},
			{1, 30, 0, 5400},
			{-1, -30, 0, -5400},
			{5, 30, 0, 19800},
			{-5, -30, 0, -19800},
			{1, 1, 1, 3661},
			{-1, -1, -1, -3661},
			{18, 0, 0, 64800},
			{-18, 0, 0, -64800},
		}

		for _, tt := range tests {
			z, err := ZoneOffsetOf(tt.hours, tt.minutes, tt.seconds)
			require.NoError(t, err)
			assert.Equal(t, tt.total, z.TotalSeconds())
		}
	})

	t.Run("invalid sign combinations", func(t *testing.T) {
		// Positive hours, negative minutes
		_, err := ZoneOffsetOf(1, -30, 0)
		assert.Error(t, err)

		// Negative hours, positive minutes
		_, err = ZoneOffsetOf(-1, 30, 0)
		assert.Error(t, err)

		// Positive minutes, negative seconds
		_, err = ZoneOffsetOf(0, 1, -30)
		assert.Error(t, err)

		// Negative minutes, positive seconds
		_, err = ZoneOffsetOf(0, -1, 30)
		assert.Error(t, err)
	})

	t.Run("out of range", func(t *testing.T) {
		_, err := ZoneOffsetOf(19, 0, 0)
		assert.Error(t, err)

		_, err = ZoneOffsetOf(-19, 0, 0)
		assert.Error(t, err)

		_, err = ZoneOffsetOf(0, 60, 0)
		assert.Error(t, err)

		_, err = ZoneOffsetOf(0, -60, 0)
		assert.Error(t, err)

		_, err = ZoneOffsetOf(0, 0, 60)
		assert.Error(t, err)

		_, err = ZoneOffsetOf(0, 0, -60)
		assert.Error(t, err)
	})
}

func TestZoneOffsetOfHours(t *testing.T) {
	tests := []struct {
		hours int
		total int
	}{
		{0, 0},
		{1, 3600},
		{-1, -3600},
		{5, 18000},
		{-5, -18000},
		{18, 64800},
		{-18, -64800},
	}

	for _, tt := range tests {
		z, err := ZoneOffsetOfHours(tt.hours)
		require.NoError(t, err)
		assert.Equal(t, tt.total, z.TotalSeconds())
		assert.Equal(t, tt.hours, z.Hours())
		assert.Equal(t, 0, z.Minutes())
		assert.Equal(t, 0, z.Seconds())
	}
}

func TestZoneOffsetOfHoursMinutes(t *testing.T) {
	tests := []struct {
		hours   int
		minutes int
		total   int
	}{
		{0, 0, 0},
		{1, 30, 5400},
		{-1, -30, -5400},
		{5, 30, 19800},
		{-5, -30, -19800},
	}

	for _, tt := range tests {
		z, err := ZoneOffsetOfHoursMinutes(tt.hours, tt.minutes)
		require.NoError(t, err)
		assert.Equal(t, tt.total, z.TotalSeconds())
		assert.Equal(t, tt.hours, z.Hours())
		assert.Equal(t, tt.minutes, z.Minutes())
		assert.Equal(t, 0, z.Seconds())
	}
}

func TestZoneOffset_IsZero(t *testing.T) {
	assert.True(t, ZoneOffsetUTC().IsZero())
	assert.True(t, MustZoneOffsetOfSeconds(0).IsZero())

	z := MustZoneOffsetOfSeconds(3600)
	assert.False(t, z.IsZero())

	z = MustZoneOffsetOfSeconds(-3600)
	assert.False(t, z.IsZero())
}

func TestZoneOffset_Compare(t *testing.T) {
	utc := ZoneOffsetUTC()
	z1 := MustZoneOffsetOfSeconds(0)
	z2 := MustZoneOffsetOfSeconds(3600)
	z3 := MustZoneOffsetOfSeconds(-3600)

	assert.Equal(t, 0, utc.Compare(z1))
	assert.Equal(t, 0, z1.Compare(utc))
	assert.Equal(t, -1, z1.Compare(z2))
	assert.Equal(t, 1, z1.Compare(z3))
	assert.Equal(t, 1, z2.Compare(z1))
	assert.Equal(t, -1, z3.Compare(z1))
	assert.Equal(t, -1, z3.Compare(z2))
}

func TestParseZoneOffset(t *testing.T) {
	t.Run("UTC formats", func(t *testing.T) {
		tests := []string{"Z", "z"}
		for _, s := range tests {
			z, err := ZoneOffsetParse(s)
			require.NoError(t, err)
			assert.Equal(t, 0, z.TotalSeconds())
			assert.True(t, z.IsZero())
		}
	})

	t.Run("hour only formats", func(t *testing.T) {
		tests := []struct {
			input   string
			seconds int
		}{
			{"+0", 0},
			{"+1", 3600},
			{"-1", -3600},
			{"+9", 32400},
			{"-9", -32400},
			{"+00", 0},
			{"+01", 3600},
			{"-01", -3600},
			{"+09", 32400},
			{"-09", -32400},
		}

		for _, tt := range tests {
			z, err := ZoneOffsetParse(tt.input)
			require.NoError(t, err, "input: %s", tt.input)
			assert.Equal(t, tt.seconds, z.TotalSeconds(), "input: %s", tt.input)
		}
	})

	t.Run("hours and minutes colon format", func(t *testing.T) {
		tests := []struct {
			input   string
			seconds int
		}{
			{"+00:00", 0},
			{"+01:00", 3600},
			{"-01:00", -3600},
			{"+01:30", 5400},
			{"-01:30", -5400},
			{"+05:30", 19800},  // India
			{"+09:30", 34200},  // Australia
			{"-05:00", -18000}, // US Eastern
			{"+18:00", 64800},
			{"-18:00", -64800},
		}

		for _, tt := range tests {
			z, err := ZoneOffsetParse(tt.input)
			require.NoError(t, err, "input: %s", tt.input)
			assert.Equal(t, tt.seconds, z.TotalSeconds(), "input: %s", tt.input)
		}
	})

	t.Run("hours and minutes compact format", func(t *testing.T) {
		tests := []struct {
			input   string
			seconds int
		}{
			{"+0000", 0},
			{"+0100", 3600},
			{"-0100", -3600},
			{"+0130", 5400},
			{"-0130", -5400},
			{"+0530", 19800},
			{"+0930", 34200},
			{"-0500", -18000},
		}

		for _, tt := range tests {
			z, err := ZoneOffsetParse(tt.input)
			require.NoError(t, err, "input: %s", tt.input)
			assert.Equal(t, tt.seconds, z.TotalSeconds(), "input: %s", tt.input)
		}
	})

	t.Run("with seconds colon format", func(t *testing.T) {
		tests := []struct {
			input   string
			seconds int
		}{
			{"+00:00:00", 0},
			{"+01:00:00", 3600},
			{"+01:01:01", 3661},
			{"-01:01:01", -3661},
			{"+01:30:45", 5445},
		}

		for _, tt := range tests {
			z, err := ZoneOffsetParse(tt.input)
			require.NoError(t, err, "input: %s", tt.input)
			assert.Equal(t, tt.seconds, z.TotalSeconds(), "input: %s", tt.input)
		}
	})

	t.Run("with seconds compact format", func(t *testing.T) {
		tests := []struct {
			input   string
			seconds int
		}{
			{"+000000", 0},
			{"+010000", 3600},
			{"+010101", 3661},
			{"-010101", -3661},
		}

		for _, tt := range tests {
			z, err := ZoneOffsetParse(tt.input)
			require.NoError(t, err, "input: %s", tt.input)
			assert.Equal(t, tt.seconds, z.TotalSeconds(), "input: %s", tt.input)
		}
	})

	t.Run("invalid formats", func(t *testing.T) {
		invalid := []string{
			"",
			"invalid",
			"1",
			"01:00",
			"1:00",
			"+",
			"-",
			"+25:00",
			"+01:60",
			"+01:00:60",
			"+19:00",
		}

		for _, s := range invalid {
			_, err := ZoneOffsetParse(s)
			assert.Error(t, err, "should fail for: %s", s)
		}
	})
}

func TestMustParseZoneOffset(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			z := MustZoneOffsetParse("+01:00")
			assert.Equal(t, 3600, z.TotalSeconds())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustZoneOffsetParse("invalid")
		})
	})
}

func TestZoneOffset_String(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "Z"},
		{3600, "+01:00"},
		{-3600, "-01:00"},
		{5400, "+01:30"},
		{-5400, "-01:30"},
		{19800, "+05:30"},
		{34200, "+09:30"},
		{-18000, "-05:00"},
		{3661, "+01:01:01"},
		{-3661, "-01:01:01"},
		{64800, "+18:00"},
		{-64800, "-18:00"},
	}

	for _, tt := range tests {
		z := MustZoneOffsetOfSeconds(tt.seconds)
		assert.Equal(t, tt.expected, z.String())
	}
}

func TestZoneOffset_AppendText(t *testing.T) {
	z := MustZoneOffsetOfSeconds(5400) // +01:30
	buf := []byte("Offset: ")
	buf, err := z.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "Offset: +01:30", string(buf))

	// UTC
	z = ZoneOffsetUTC()
	buf = []byte("Offset: ")
	buf, err = z.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "Offset: Z", string(buf))
}

func TestZoneOffset_MarshalText(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "Z"},
		{3600, "+01:00"},
		{-3600, "-01:00"},
		{5400, "+01:30"},
	}

	for _, tt := range tests {
		z := MustZoneOffsetOfSeconds(tt.seconds)
		text, err := z.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, tt.expected, string(text))
	}
}

func TestZoneOffset_UnmarshalText(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var z ZoneOffset
		err := z.UnmarshalText([]byte("+01:30"))
		require.NoError(t, err)
		assert.Equal(t, 5400, z.TotalSeconds())
	})

	t.Run("UTC", func(t *testing.T) {
		var z ZoneOffset
		err := z.UnmarshalText([]byte("Z"))
		require.NoError(t, err)
		assert.Equal(t, 0, z.TotalSeconds())
	})

	t.Run("empty", func(t *testing.T) {
		var z ZoneOffset
		err := z.UnmarshalText([]byte(""))
		assert.Error(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		var z ZoneOffset
		err := z.UnmarshalText([]byte("invalid"))
		assert.Error(t, err)
	})
}

func TestZoneOffset_MarshalJSON(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, `"Z"`},
		{3600, `"+01:00"`},
		{-3600, `"-01:00"`},
		{5400, `"+01:30"`},
	}

	for _, tt := range tests {
		z := MustZoneOffsetOfSeconds(tt.seconds)
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, tt.expected, string(data))
	}
}

func TestZoneOffset_UnmarshalJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var z ZoneOffset
		err := json.Unmarshal([]byte(`"+01:30"`), &z)
		require.NoError(t, err)
		assert.Equal(t, 5400, z.TotalSeconds())
	})

	t.Run("UTC", func(t *testing.T) {
		var z ZoneOffset
		err := json.Unmarshal([]byte(`"Z"`), &z)
		require.NoError(t, err)
		assert.Equal(t, 0, z.TotalSeconds())
	})

	t.Run("invalid", func(t *testing.T) {
		var z ZoneOffset
		err := json.Unmarshal([]byte(`"invalid"`), &z)
		assert.Error(t, err)
	})
}

func TestZoneOffset_Constants(t *testing.T) {
	assert.Equal(t, 0, ZoneOffsetUTC().TotalSeconds())
	assert.Equal(t, -64800, ZoneOffsetMin().TotalSeconds())
	assert.Equal(t, 64800, ZoneOffsetMax().TotalSeconds())
}

func TestZoneOffset_RoundTrip(t *testing.T) {
	tests := []int{
		0, 3600, -3600, 5400, -5400,
		19800, 34200, -18000, 3661, -3661,
	}

	for _, seconds := range tests {
		z1 := MustZoneOffsetOfSeconds(seconds)

		// String round-trip
		s := z1.String()
		z2, err := ZoneOffsetParse(s)
		require.NoError(t, err)
		assert.Equal(t, z1.TotalSeconds(), z2.TotalSeconds())

		// JSON round-trip
		data, err := json.Marshal(z1)
		require.NoError(t, err)
		var z3 ZoneOffset
		err = json.Unmarshal(data, &z3)
		require.NoError(t, err)
		assert.Equal(t, z1.TotalSeconds(), z3.TotalSeconds())
	}
}

func TestZoneOffset_CommonOffsets(t *testing.T) {
	// Test some common real-world offsets
	tests := []struct {
		name    string
		offset  string
		seconds int
	}{
		{"UTC", "Z", 0},
		{"New York (EST)", "-05:00", -18000},
		{"Los Angeles (PST)", "-08:00", -28800},
		{"Paris (CET)", "+01:00", 3600},
		{"Tokyo (JST)", "+09:00", 32400},
		{"Sydney (AEDT)", "+11:00", 39600},
		{"India (IST)", "+05:30", 19800},
		{"Nepal", "+05:45", 20700},
		{"Adelaide", "+09:30", 34200},
		{"Newfoundland", "-03:30", -12600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := ZoneOffsetParse(tt.offset)
			require.NoError(t, err)
			assert.Equal(t, tt.seconds, z.TotalSeconds())

			// Round-trip
			assert.Equal(t, tt.offset, z.String())
		})
	}
}
