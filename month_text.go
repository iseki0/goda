package goda

import "time"

// String returns the English name of the month (e.g., "January", "February").
// Returns empty string for zero value.
func (m Month) String() string {
	if m.IsZero() {
		return ""
	}
	return time.Month(m).String()
}
