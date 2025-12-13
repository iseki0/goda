package goda

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
)

// ZoneOffset represents a time-zone offset from UTC, such as +02:00.
//
// A time-zone offset is the amount of time that a time-zone differs from Greenwich/UTC.
// This is usually a fixed number of hours and minutes.
//
// Different parts of the world have different time-zone offsets.
// The rules for how offsets vary by place and time of year are captured in the ZoneId class.
//
// For example, Paris is one hour ahead of Greenwich/UTC in winter and two hours ahead in summer.
// The ZoneId instance for Paris will reference two ZoneOffset instances - a +01:00 instance for winter,
// and a +02:00 instance for summer.
//
// ZoneOffset is comparable and can be used as a map key.
// The zero value represents UTC (+00:00).
//
// ZoneOffset implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: ±HH:mm, ±HH:mm:ss, or Z for UTC (e.g., "+02:00", "-05:30", "Z").
// The seconds field is output only if non-zero.
type ZoneOffset struct {
	totalSeconds int32
}

func (z ZoneOffset) IsSupportedField(field Field) bool {
	return field == FieldOffsetSeconds
}

func (z ZoneOffset) GetField(field Field) TemporalValue {
	if field == FieldOffsetSeconds {
		return TemporalValue{
			v: int64(z.totalSeconds),
		}
	}
	return TemporalValue{unsupported: true}
}

// TotalSeconds returns the total zone offset in seconds.
func (z ZoneOffset) TotalSeconds() int {
	return int(z.totalSeconds)
}

// IsZero returns true if this is the zero value (UTC).
// Note: This returns true for UTC (+00:00), not for invalid offsets.
func (z ZoneOffset) IsZero() bool {
	return z.totalSeconds == 0
}

// Hours returns the hours component of the offset.
// For negative offsets, this returns a negative value.
func (z ZoneOffset) Hours() int {
	return int(z.totalSeconds / 3600)
}

// Minutes returns the minutes component of the offset (0-59).
// For negative offsets, this returns a negative value if the offset is not a whole number of hours.
func (z ZoneOffset) Minutes() int {
	remainder := z.totalSeconds % 3600
	return int(remainder / 60)
}

// Seconds returns the seconds component of the offset (0-59).
// For negative offsets, this returns a negative value if the offset is not a whole number of minutes.
func (z ZoneOffset) Seconds() int {
	return int(z.totalSeconds % 60)
}

// Compare compares this zone offset with another.
// Returns -1 if this offset is less (more westward), 0 if equal, 1 if greater (more eastward).
func (z ZoneOffset) Compare(other ZoneOffset) int {
	if z.totalSeconds < other.totalSeconds {
		return -1
	} else if z.totalSeconds > other.totalSeconds {
		return 1
	}
	return 0
}

// String returns the string representation of the zone offset.
// Returns "Z" for UTC, otherwise returns the format ±HH:MM or ±HH:MM:SS.
func (z ZoneOffset) String() string {
	return stringImpl(z)
}

// AppendText implements encoding.TextAppender.
func (z ZoneOffset) AppendText(b []byte) ([]byte, error) {
	if z.totalSeconds == 0 {
		return append(b, 'Z'), nil
	}

	absSeconds := z.totalSeconds
	if absSeconds < 0 {
		b = append(b, '-')
		absSeconds = -absSeconds
	} else {
		b = append(b, '+')
	}

	hours := absSeconds / 3600
	minutes := (absSeconds % 3600) / 60
	seconds := absSeconds % 60

	// Format hours (always 2 digits)
	b = append(b, byte('0'+hours/10), byte('0'+hours%10))
	b = append(b, ':')
	// Format minutes (always 2 digits)
	b = append(b, byte('0'+minutes/10), byte('0'+minutes%10))

	// Only append seconds if non-zero
	if seconds != 0 {
		b = append(b, ':')
		b = append(b, byte('0'+seconds/10), byte('0'+seconds%10))
	}

	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (z ZoneOffset) MarshalText() ([]byte, error) {
	return marshalTextImpl(z)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (z *ZoneOffset) UnmarshalText(text []byte) error {
	s := string(text)
	if len(s) == 0 {
		return newError("zone offset cannot be empty")
	}

	// Handle UTC
	if s == "Z" || s == "z" {
		*z = ZoneOffsetUTC()
		return nil
	}

	// Must start with + or -
	if s[0] != '+' && s[0] != '-' {
		return newError("zone offset must start with + or -, got %q", s)
	}

	negative := s[0] == '-'
	s = s[1:] // Remove sign

	var hours, minutes, seconds int
	var err error

	// Determine format based on length and colons
	if len(s) == 0 {
		return newError("zone offset has no digits after sign")
	}

	// check for colon-separated format
	hasColon := false
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			hasColon = true
			break
		}
	}

	if hasColon {
		// Colon-separated format: HH:MM or HH:MM:SS or H:MM
		parts := []string{}
		start := 0
		for i := 0; i <= len(s); i++ {
			if i == len(s) || s[i] == ':' {
				if i > start {
					parts = append(parts, s[start:i])
				}
				start = i + 1
			}
		}

		if len(parts) < 2 || len(parts) > 3 {
			return newError("invalid zone offset format %q", string(text))
		}

		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return newError("invalid zone offset hours: %v", err)
		}

		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return newError("invalid zone offset minutes: %v", err)
		}

		if len(parts) == 3 {
			seconds, err = strconv.Atoi(parts[2])
			if err != nil {
				return newError("invalid zone offset seconds: %v", err)
			}
		}
	} else {
		// Compact format: H, HH, HHMM, or HHMMSS
		switch len(s) {
		case 1, 2: // H or HH
			hours, err = strconv.Atoi(s)
			if err != nil {
				return newError("invalid zone offset hours: %v", err)
			}
		case 4: // HHMM
			hours, err = strconv.Atoi(s[0:2])
			if err != nil {
				return newError("invalid zone offset hours: %v", err)
			}
			minutes, err = strconv.Atoi(s[2:4])
			if err != nil {
				return newError("invalid zone offset minutes: %v", err)
			}
		case 6: // HHMMSS
			hours, err = strconv.Atoi(s[0:2])
			if err != nil {
				return newError("invalid zone offset hours: %v", err)
			}
			minutes, err = strconv.Atoi(s[2:4])
			if err != nil {
				return newError("invalid zone offset minutes: %v", err)
			}
			seconds, err = strconv.Atoi(s[4:6])
			if err != nil {
				return newError("invalid zone offset seconds: %v", err)
			}
		default:
			return newError("invalid zone offset format %q", string(text))
		}
	}

	// Apply sign
	if negative {
		hours = -hours
		minutes = -minutes
		seconds = -seconds
	}

	offset, err := ZoneOffsetOf(hours, minutes, seconds)
	if err != nil {
		return err
	}

	*z = offset
	return nil
}

// MarshalJSON implements json.Marshaler.
func (z ZoneOffset) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(z)
}

// UnmarshalJSON implements json.Unmarshaler.
func (z *ZoneOffset) UnmarshalJSON(data []byte) error {
	return unmarshalJsonImpl(z, data)
}

// ZoneOffsetUTC returns UTC offset (+00:00).
func ZoneOffsetUTC() ZoneOffset {
	return ZoneOffset{totalSeconds: 0}
}

// ZoneOffsetMin returns the minimum supported offset (-18:00).
func ZoneOffsetMin() ZoneOffset {
	return ZoneOffset{totalSeconds: -18 * 3600}
}

// ZoneOffsetMax returns the maximum supported offset (+18:00).
func ZoneOffsetMax() ZoneOffset {
	return ZoneOffset{totalSeconds: 18 * 3600}
}

// ZoneOffsetOfSeconds creates a ZoneOffset from the total offset in seconds.
// The offset must be in the range -18:00 to +18:00, which corresponds to -64800 to +64800 seconds.
//
// Returns an error if the offset is outside the valid range.
func ZoneOffsetOfSeconds(seconds int) (ZoneOffset, error) {
	if seconds > 18*3600 || seconds < -18*3600 {
		return ZoneOffset{}, newError("zone offset seconds must be in range -64800 to +64800, got %d", seconds)
	}
	return ZoneOffset{totalSeconds: int32(seconds)}, nil
}

// MustZoneOffsetOfSeconds creates a ZoneOffset from the total offset in seconds.
// Panics if the offset is outside the valid range.
func MustZoneOffsetOfSeconds(seconds int) ZoneOffset {
	return mustValue(ZoneOffsetOfSeconds(seconds))
}

// ZoneOffsetOf creates a ZoneOffset from hours, minutes, and seconds.
// The offset must be in the range -18:00 to +18:00.
//
// The sign of all components must be the same. If any component is negative,
// all non-zero components must be negative or zero.
//
// Returns an error if the offset is invalid.
func ZoneOffsetOf(hours, minutes, seconds int) (ZoneOffset, error) {
	// Validate that signs are consistent
	if hours < 0 {
		if minutes > 0 || seconds > 0 {
			return ZoneOffset{}, newError("zone offset minutes and seconds must not be positive when hours is negative")
		}
	} else if hours > 0 {
		if minutes < 0 || seconds < 0 {
			return ZoneOffset{}, newError("zone offset minutes and seconds must not be negative when hours is positive")
		}
	} else if minutes < 0 {
		if seconds > 0 {
			return ZoneOffset{}, newError("zone offset seconds must not be positive when minutes is negative")
		}
	} else if minutes > 0 {
		if seconds < 0 {
			return ZoneOffset{}, newError("zone offset seconds must not be negative when minutes is positive")
		}
	}

	// Validate ranges
	if hours < -18 || hours > 18 {
		return ZoneOffset{}, newError("zone offset hours must be in range -18 to +18, got %d", hours)
	}
	if minutes < -59 || minutes > 59 {
		return ZoneOffset{}, newError("zone offset minutes must be in range -59 to +59, got %d", minutes)
	}
	if seconds < -59 || seconds > 59 {
		return ZoneOffset{}, newError("zone offset seconds must be in range -59 to +59, got %d", seconds)
	}

	totalSeconds := hours*3600 + minutes*60 + seconds
	return ZoneOffsetOfSeconds(totalSeconds)
}

// MustZoneOffsetOf creates a ZoneOffset from hours, minutes, and seconds.
// Panics if the offset is invalid.
func MustZoneOffsetOf(hours, minutes, seconds int) ZoneOffset {
	return mustValue(ZoneOffsetOf(hours, minutes, seconds))
}

// ZoneOffsetOfHours creates a ZoneOffset from hours only.
// The offset must be in the range -18 to +18.
func ZoneOffsetOfHours(hours int) (ZoneOffset, error) {
	return ZoneOffsetOf(hours, 0, 0)
}

// MustZoneOffsetOfHours creates a ZoneOffset from hours only.
// Panics if the offset is invalid.
func MustZoneOffsetOfHours(hours int) ZoneOffset {
	return mustValue(ZoneOffsetOfHours(hours))
}

// ZoneOffsetOfHoursMinutes creates a ZoneOffset from hours and minutes.
// The offset must be in the range -18:00 to +18:00.
func ZoneOffsetOfHoursMinutes(hours, minutes int) (ZoneOffset, error) {
	return ZoneOffsetOf(hours, minutes, 0)
}

// MustZoneOffsetOfHoursMinutes creates a ZoneOffset from hours and minutes.
// Panics if the offset is invalid.
func MustZoneOffsetOfHoursMinutes(hours, minutes int) ZoneOffset {
	return mustValue(ZoneOffsetOfHoursMinutes(hours, minutes))
}

// ZoneOffsetParse parses a zone offset string.
// Accepts the following formats:
//   - Z: UTC
//   - ±H: hour offset (-9, +9)
//   - ±HH: hour offset (-09, +09)
//   - ±HH:MM: hour and minute offset (-09:30, +09:30)
//   - ±HHMM: hour and minute offset (-0930, +0930)
//   - ±HH:MM:SS: hour, minute, and second offset (-09:30:45, +09:30:45)
//   - ±HHMMSS: hour, minute, and second offset (-093045, +093045)
//
// Returns an error if the string is invalid.
func ZoneOffsetParse(s string) (ZoneOffset, error) {
	var z ZoneOffset
	err := z.UnmarshalText([]byte(s))
	return z, err
}

// MustZoneOffsetParse parses a zone offset string.
// Panics if the string is invalid.
func MustZoneOffsetParse(s string) ZoneOffset {
	return mustValue(ZoneOffsetParse(s))
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*ZoneOffset)(nil)
	_ fmt.Stringer             = (*ZoneOffset)(nil)
	_ encoding.TextMarshaler   = (*ZoneOffset)(nil)
	_ encoding.TextUnmarshaler = (*ZoneOffset)(nil)
	_ json.Marshaler           = (*ZoneOffset)(nil)
	_ json.Unmarshaler         = (*ZoneOffset)(nil)
	_ TemporalAccessor         = (*ZoneOffset)(nil)
)

// Compile-time check that ZoneOffset is comparable
func _assertZoneOffsetIsComparable[T comparable](t T) {}

var _ = _assertZoneOffsetIsComparable[ZoneOffset]
