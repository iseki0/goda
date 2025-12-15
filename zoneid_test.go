package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZoneIdOf(t *testing.T) {
	t.Run("valid zone IDs", func(t *testing.T) {
		testCases := []string{
			"UTC",
			"America/New_York",
			"Europe/London",
			"Asia/Tokyo",
			"Australia/Sydney",
		}

		for _, tc := range testCases {
			z, err := ZoneIdOf(tc)
			require.NoError(t, err, "zone: %s", tc)
			assert.False(t, z.IsZero())
			assert.Equal(t, tc, z.String())
		}
	})

	t.Run("invalid zone ID", func(t *testing.T) {
		_, err := ZoneIdOf("Invalid/Zone")
		assert.Error(t, err)
	})

	t.Run("empty zone ID", func(t *testing.T) {
		// On Windows, empty string is treated as Local time, so this may not error
		z, err := ZoneIdOf("")
		if err == nil {
			// If no error, it should return Local
			assert.False(t, z.IsZero())
		}
	})
}

func TestMustZoneIdOf(t *testing.T) {
	t.Run("valid zone ID", func(t *testing.T) {
		assert.NotPanics(t, func() {
			z := MustZoneIdOf("America/New_York")
			assert.Equal(t, "America/New_York", z.String())
		})
	})

	t.Run("invalid zone ID panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustZoneIdOf("Invalid/Zone")
		})
	})
}

func TestZoneIdUTC(t *testing.T) {
	z := ZoneIdUTC()
	assert.False(t, z.IsZero())
	assert.Equal(t, "UTC", z.String())
}

func TestZoneIdOfGoLocation(t *testing.T) {
	t.Run("from time.UTC", func(t *testing.T) {
		z := ZoneIdOfGoLocation(time.UTC)
		assert.False(t, z.IsZero())
		assert.Equal(t, "UTC", z.String())
	})

	t.Run("from time.Local", func(t *testing.T) {
		z := ZoneIdOfGoLocation(time.Local)
		assert.False(t, z.IsZero())
		assert.NotEmpty(t, z.String())
	})

	t.Run("from loaded location", func(t *testing.T) {
		loc, err := time.LoadLocation("Asia/Shanghai")
		require.NoError(t, err)
		z := ZoneIdOfGoLocation(loc)
		assert.Equal(t, "Asia/Shanghai", z.String())
	})
}

func TestZoneIdDefault(t *testing.T) {
	z := ZoneIdDefault()
	assert.False(t, z.IsZero())
	assert.NotEmpty(t, z.String())
}

func TestZoneId_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		assert.True(t, z.IsZero())
		assert.Empty(t, z.String())
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := ZoneIdUTC()
		assert.False(t, z.IsZero())
	})
}

func TestZoneId_String(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		assert.Empty(t, z.String())
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := MustZoneIdOf("America/Los_Angeles")
		assert.Equal(t, "America/Los_Angeles", z.String())
	})
}

func TestZoneId_MarshalJSON(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		z := MustZoneIdOf("Asia/Tokyo")
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, `"Asia/Tokyo"`, string(data))
	})

	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, `""`, string(data))
	})

	t.Run("UTC", func(t *testing.T) {
		z := ZoneIdUTC()
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, `"UTC"`, string(data))
	})
}

func TestZoneId_UnmarshalJSON(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		var z ZoneId
		err := json.Unmarshal([]byte(`"America/New_York"`), &z)
		require.NoError(t, err)
		assert.Equal(t, "America/New_York", z.String())
	})

	t.Run("null", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := json.Unmarshal([]byte(`null`), &z)
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("empty string", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := json.Unmarshal([]byte(`""`), &z)
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("invalid zone", func(t *testing.T) {
		var z ZoneId
		err := json.Unmarshal([]byte(`"Invalid/Zone"`), &z)
		assert.Error(t, err)
	})
}

func TestZoneId_MarshalText(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		z := MustZoneIdOf("Europe/Berlin")
		data, err := z.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, "Europe/Berlin", string(data))
	})

	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		data, err := z.MarshalText()
		require.NoError(t, err)
		assert.Empty(t, data)
	})
}

func TestZoneId_UnmarshalText(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		var z ZoneId
		err := z.UnmarshalText([]byte("Asia/Seoul"))
		require.NoError(t, err)
		assert.Equal(t, "Asia/Seoul", z.String())
	})

	t.Run("empty", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := z.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("invalid zone", func(t *testing.T) {
		var z ZoneId
		err := z.UnmarshalText([]byte("Not/A/Zone"))
		assert.Error(t, err)
	})
}

func TestZoneId_AppendText(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		buf := []byte("prefix:")
		result, err := z.AppendText(buf)
		require.NoError(t, err)
		assert.Equal(t, "prefix:", string(result))
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := MustZoneIdOf("Pacific/Auckland")
		buf := []byte("zone=")
		result, err := z.AppendText(buf)
		require.NoError(t, err)
		assert.Equal(t, "zone=Pacific/Auckland", string(result))
	})

	t.Run("empty buffer", func(t *testing.T) {
		z := MustZoneIdOf("UTC")
		result, err := z.AppendText(nil)
		require.NoError(t, err)
		assert.Equal(t, "UTC", string(result))
	})
}

func TestZoneId_Scan(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := z.Scan(nil)
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("from string", func(t *testing.T) {
		var z ZoneId
		err := z.Scan("Europe/Paris")
		require.NoError(t, err)
		assert.Equal(t, "Europe/Paris", z.String())
	})

	t.Run("from []byte", func(t *testing.T) {
		var z ZoneId
		err := z.Scan([]byte("Asia/Dubai"))
		require.NoError(t, err)
		assert.Equal(t, "Asia/Dubai", z.String())
	})

	t.Run("from invalid type", func(t *testing.T) {
		var z ZoneId
		err := z.Scan(12345)
		assert.Error(t, err)
	})

	t.Run("from invalid zone", func(t *testing.T) {
		var z ZoneId
		err := z.Scan("Bad/Zone")
		assert.Error(t, err)
	})
}

func TestZoneId_Value(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		val, err := z.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := MustZoneIdOf("America/Chicago")
		val, err := z.Value()
		require.NoError(t, err)
		assert.Equal(t, "America/Chicago", val)
	})
}

func TestZoneId_RoundTrip(t *testing.T) {
	testZones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
		"Pacific/Honolulu",
	}

	for _, zoneName := range testZones {
		t.Run(zoneName, func(t *testing.T) {
			// Create original
			original := MustZoneIdOf(zoneName)

			// JSON round-trip
			jsonData, err := json.Marshal(original)
			require.NoError(t, err)

			var fromJSON ZoneId
			err = json.Unmarshal(jsonData, &fromJSON)
			require.NoError(t, err)
			assert.Equal(t, original.String(), fromJSON.String())

			// Text round-trip
			textData, err := original.MarshalText()
			require.NoError(t, err)

			var fromText ZoneId
			err = fromText.UnmarshalText(textData)
			require.NoError(t, err)
			assert.Equal(t, original.String(), fromText.String())

			// SQL round-trip
			val, err := original.Value()
			require.NoError(t, err)

			var fromSQL ZoneId
			err = fromSQL.Scan(val)
			require.NoError(t, err)
			assert.Equal(t, original.String(), fromSQL.String())
		})
	}
}

func TestZoneId_Comparable(t *testing.T) {
	z1 := MustZoneIdOf("America/New_York")
	z2 := MustZoneIdOf("America/New_York")
	z3 := MustZoneIdOf("Europe/London")
	var zero ZoneId

	// Test equality
	assert.True(t, z1 == z2)
	assert.False(t, z1 == z3)
	assert.True(t, zero == ZoneId{})

	// Can be used as map key
	m := make(map[ZoneId]string)
	m[z1] = "New York"
	m[z3] = "London"
	assert.Equal(t, "New York", m[z2]) // z2 equals z1
	assert.Equal(t, "London", m[z3])
}

func TestZoneId_ConcurrentLoadLocation(t *testing.T) {
	// Test that concurrent zone ID creation is safe
	const numGoroutines = 100
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			z, err := ZoneIdOf("America/New_York")
			assert.NoError(t, err)
			assert.Equal(t, "America/New_York", z.String())
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
