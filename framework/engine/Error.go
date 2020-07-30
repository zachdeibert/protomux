package engine

import (
	"fmt"

	"github.com/zachdeibert/protomux/config"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeUnknownProtocol represents when a protocol is specified that does not have an implementaiton
	ErrorCodeUnknownProtocol ErrorCode = iota
	// ErrorCodeHostLookup represents when an error occurs while looking up a host name
	ErrorCodeHostLookup ErrorCode = iota
	// ErrorCodeNoHostRecords represents when a host was not able to be resolved
	ErrorCodeNoHostRecords ErrorCode = iota
	// ErrorCodeListenStart represents when a listener could not be started
	ErrorCodeListenStart ErrorCode = iota
	// ErrorCodeClosed represents when the socket has been closed
	ErrorCodeClosed ErrorCode = iota
	// ErrorCodeDeadlinesNotSupported is returned when a deadline is attempted to be set
	ErrorCodeDeadlinesNotSupported ErrorCode = iota
)

// Error describes an engine error
type Error struct {
	Message string
	Code    ErrorCode
}

func (e Error) Error() string {
	return e.Message
}

// ErrorUnknownProtocol creates a new ErrorUnknownProtocol error
func ErrorUnknownProtocol(name string) error {
	return &Error{
		Message: fmt.Sprintf("Unknown protocol name '%s'", name),
	}
}

// ErrorHostLookup creates a new ErrorHostLookup error
func ErrorHostLookup(name string, err error) error {
	return &Error{
		Message: fmt.Sprintf("Could not lookup host '%s': %s", name, err),
		Code:    ErrorCodeHostLookup,
	}
}

// ErrorNoHostRecords creates a new ErrorNoHostRecords error
func ErrorNoHostRecords(name string) error {
	return &Error{
		Message: fmt.Sprintf("Unable to resolve hostname '%s'", name),
		Code:    ErrorCodeNoHostRecords,
	}
}

// ErrorListenStart creates a new ErrorListenStart error
func ErrorListenStart(addr config.Connection, err error) error {
	return &Error{
		Message: fmt.Sprintf("Unable to listen on %s: %s", addr, err),
		Code:    ErrorCodeListenStart,
	}
}

// ErrorClosed error
var ErrorClosed error = &Error{
	Message: "Socket closed",
	Code:    ErrorCodeClosed,
}

// ErrorDeadlinesNotSupported error
var ErrorDeadlinesNotSupported error = &Error{
	Message: "Setting deadlines is not supported",
	Code:    ErrorCodeDeadlinesNotSupported,
}
