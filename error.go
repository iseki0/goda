package goda

import (
	"fmt"
)

// newError creates a new Error with the given format and arguments.
// All error messages are prefixed with "goda: ".
func newError(format string, a ...any) error {
	return &Error{message: "goda: " + fmt.Sprintf(format, a...)}
}

// unmarshalError creates an error for invalid unmarshaling input.
func unmarshalError(userInput []byte) error {
	return newError("unable to unmarshal user input: %q", string(userInput))
}

// sqlScannerDefaultBranch creates an error for unsupported SQL scan types.
func sqlScannerDefaultBranch(value any) error {
	return newError("cannot scan value of type %T", value)
}

// Error is the error type used by this package.
// It wraps error messages with the "goda: " prefix.
type Error struct {
	unsupportedField Field
	outOfRangeField  Field
	outOfRangeValue  int64
	message          string
}

// Error implements the error interface.
func (e Error) Error() string {
	if e.unsupportedField.Valid() {
		return fmt.Sprintf("goda: unsupported field %s", e.unsupportedField)
	}
	if e.outOfRangeField.Valid() {
		return fmt.Sprintf("goda: invalid value of %s (valid range %d - %d): %d", e.outOfRangeField, e.outOfRangeField.fieldRange().Min, e.outOfRangeField.fieldRange().Max, e.outOfRangeValue)
	}
	return e.message
}

func overflowError() error {
	return newError("overflow")
}

func fieldOutOfRangeError(field Field, value int64) error {
	return &Error{outOfRangeField: field, outOfRangeValue: value}
}

func unsupportedField(field Field) error {
	return &Error{unsupportedField: field}
}
