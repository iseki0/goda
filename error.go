package goda

import (
	"fmt"
)

// newError creates a new Error with the given format and arguments.
// All error messages are prefixed with "goda: ".
func newError(format string, a ...any) error {
	return &Error{message: fmt.Sprintf(format, a...)}
}

// unmarshalError creates an error for invalid unmarshaling input.
func unmarshalError(userInput []byte) error {
	return newError("unable to unmarshal user input: %q", string(userInput))
}

// sqlScannerDefaultBranch creates an error for unsupported SQL scan types.
func sqlScannerDefaultBranch(value any) error {
	return newError("cannot scan value of type %T", value)
}

const (
	errReasonInvalidField int = iota + 1
	errReasonUnsupportedField
	errReasonOutOfRange
	errReasonArithmeticOverflow
)

// Error is the error type used by this package.
// It wraps error messages with the "goda: " prefix.
type Error struct {
	reason     int
	field      Field
	int64Value int64
	message    string
	cause      error
	typeName   string
	funcName   string
}

// Error implements the error interface.
func (e Error) Error() string {
	var text string
	switch e.reason {
	case errReasonInvalidField:
		text = fmt.Sprintf("goda: invalid field (value=%d)", int64(e.field))
	case errReasonUnsupportedField:
		text = fmt.Sprintf("goda: unsupported field %s", e.field)
	case errReasonOutOfRange:
		fr := e.field.fieldRange()
		text = fmt.Sprintf("goda: invalid value of %s (valid range %d - %d): %d", e.field, fr.Min, fr.Max, e.int64Value)
	case errReasonArithmeticOverflow:
		text = "goda: arithmetic overflow"
	default:
		text = "goda: " + e.message
	}
	if e.typeName != "" {
		text += " at " + e.typeName + "/" + e.funcName
	}
	if e.cause != nil {
		text += ", caused by: " + e.cause.Error()
	}
	return text
}

func (e Error) Unwrap() error {
	return e.cause
}

func overflowError() error {
	return &Error{reason: errReasonArithmeticOverflow}
}

func fieldOutOfRangeError(field Field, value int64) error {
	return &Error{reason: errReasonOutOfRange, field: field, int64Value: value}
}

func unsupportedField(field Field) error {
	return &Error{reason: errReasonUnsupportedField, field: field}
}

func invalidFieldError(field Field) error {
	return &Error{reason: errReasonInvalidField, field: field}
}
