package goda

import (
	"fmt"
	"strconv"
)

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
func (z *ZoneOffset) UnmarshalText(text []byte) (e error) {
	defer deferOpInParse(text, &e)
	s := string(text)
	if len(s) == 0 {
		return fmt.Errorf("zone offset cannot be empty")
	}

	// Handle UTC
	if s == "Z" || s == "z" {
		*z = ZoneOffsetUTC()
		return nil
	}

	// Must start with + or -
	if s[0] != '+' && s[0] != '-' {
		return fmt.Errorf("zone offset must start with + or -, got %q", s)
	}

	negative := s[0] == '-'
	s = s[1:] // Remove sign

	var hours, minutes, seconds int
	var err error

	// Determine format based on length and colons
	if len(s) == 0 {
		return fmt.Errorf("zone offset has no digits after sign")
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
		var parts []string
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
			return fmt.Errorf("invalid zone offset format %q", string(text))
		}

		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid zone offset hours: %v", err)
		}

		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid zone offset minutes: %v", err)
		}

		if len(parts) == 3 {
			seconds, err = strconv.Atoi(parts[2])
			if err != nil {
				return fmt.Errorf("invalid zone offset seconds: %v", err)
			}
		}
	} else {
		// Compact format: H, HH, HHMM, or HHMMSS
		switch len(s) {
		case 1, 2: // H or HH
			hours, err = strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("invalid zone offset hours: %v", err)
			}
		case 4: // HHMM
			hours, err = strconv.Atoi(s[0:2])
			if err != nil {
				return fmt.Errorf("invalid zone offset hours: %v", err)
			}
			minutes, err = strconv.Atoi(s[2:4])
			if err != nil {
				return fmt.Errorf("invalid zone offset minutes: %v", err)
			}
		case 6: // HHMMSS
			hours, err = strconv.Atoi(s[0:2])
			if err != nil {
				return fmt.Errorf("invalid zone offset hours: %v", err)
			}
			minutes, err = strconv.Atoi(s[2:4])
			if err != nil {
				return fmt.Errorf("invalid zone offset minutes: %v", err)
			}
			seconds, err = strconv.Atoi(s[4:6])
			if err != nil {
				return fmt.Errorf("invalid zone offset seconds: %v", err)
			}
		default:
			return fmt.Errorf("invalid zone offset format %q", string(text))
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
