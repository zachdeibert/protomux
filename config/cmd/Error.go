package cmd

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeUnrecognizedFlag represents when a flag was received that was not recognized
	ErrorCodeUnrecognizedFlag ErrorCode = iota
	// ErrorCodeMissingArgs represents when a flag did not have enough arguments
	ErrorCodeMissingArgs ErrorCode = iota
	// ErrorCodeTooManyArgs represents when too many arguments were given to a flag
	ErrorCodeTooManyArgs ErrorCode = iota
	// ErrorCodeTreePathFormat when a tree path was in an invalid format
	ErrorCodeTreePathFormat ErrorCode = iota
)

// Error describes an AST parsing error
type Error struct {
	Message  string
	Code     ErrorCode
	Location common.Location
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (at %s)\n%s", e.Message, e.Location.ShortString(), e.Location)
}

// ErrorUnrecognizedFlag creates a new ErrorUnrecognizedFlag error
func ErrorUnrecognizedFlag(arg string, location common.Location) error {
	return &Error{
		Message:  fmt.Sprintf("Unknown flag '%s'", arg),
		Code:     ErrorCodeUnrecognizedFlag,
		Location: location,
	}
}

// ErrorMissingArgs creates new ErrorMissingArgs error
func ErrorMissingArgs(arg string, location common.Location) error {
	return &Error{
		Message:  fmt.Sprintf("Flag %s missing arguments", arg),
		Code:     ErrorCodeMissingArgs,
		Location: location,
	}
}

// ErrorTooManyArgs creates a new ErrorTooManyArgs error
func ErrorTooManyArgs(arg string, location common.Location) error {
	return &Error{
		Message:  fmt.Sprintf("Flag %s has too many arguments", arg),
		Code:     ErrorCodeTooManyArgs,
		Location: location,
	}
}

// ErrorTreePathFormat creates a new ErrorTreePathFormat error
func ErrorTreePathFormat(desc string, location common.Location) error {
	return &Error{
		Message:  desc,
		Code:     ErrorCodeTreePathFormat,
		Location: location,
	}
}
