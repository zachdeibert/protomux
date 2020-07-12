package minecraft

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeInvalidVersion represents when an unknown version is specified
	ErrorCodeInvalidVersion ErrorCode = iota
	// ErrorCodeMultipleValues represents when a property that should have only had one value has multiple
	ErrorCodeMultipleValues ErrorCode = iota
	// ErrorCodeParameterRequirement represents whan a requirement for a parameter is not met
	ErrorCodeParameterRequirement ErrorCode = iota
	// ErrorCodeUnrecognizedParameter represents when a parameter name is not recognized
	ErrorCodeUnrecognizedParameter ErrorCode = iota
	// ErrorCodeUnknownRemoteType represents when an unknown report type is specified
	ErrorCodeUnknownRemoteType ErrorCode = iota
)

// Error describes an error with the Minecraft protocol implementation
type Error struct {
	Message string
	Code    ErrorCode
}

func (e Error) Error() string {
	return e.Message
}

// ErrorInvalidVersion creates a new ErrorInvalidVersion error
func ErrorInvalidVersion(version string) error {
	return &Error{
		Message: fmt.Sprintf("Unrecognized version moniker '%s'", version),
		Code:    ErrorCodeInvalidVersion,
	}
}

// ErrorMultipleValues creates a new ErrorMultipleValues error
func ErrorMultipleValues(param string) error {
	return &Error{
		Message: fmt.Sprintf("Parameter '%s' can only have one value, but has an array", param),
		Code:    ErrorCodeMultipleValues,
	}
}

// ErrorParameterRequirement creates a new ErrorParameterRequirement error
func ErrorParameterRequirement(message string) error {
	return &Error{
		Message: message,
		Code:    ErrorCodeParameterRequirement,
	}
}

// ErrorUnrecognizedParameter creates a new ErrorUnrecognizedParameter error
func ErrorUnrecognizedParameter(name string, location common.Location) error {
	return &Error{
		Message: fmt.Sprintf("Unrecognized parameter '%s' (at %s)\n%s", name, location.ShortString(), location),
		Code:    ErrorCodeUnrecognizedParameter,
	}
}

// ErrorUnknownRemoteType creates a new ErrorUnknownRemoteType error
func ErrorUnknownRemoteType(name string) error {
	return &Error{
		Message: fmt.Sprintf("Unrecognized remote type '%s'", name),
		Code:    ErrorCodeUnknownRemoteType,
	}
}
