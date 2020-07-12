package listener

import (
	"fmt"

	"github.com/zachdeibert/protomux/config"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeHostLookup represents when an error occurs while looking up a host name
	ErrorCodeHostLookup ErrorCode = iota
	// ErrorCodeNoHostRecords represents when a host was not able to be resolved
	ErrorCodeNoHostRecords ErrorCode = iota
	// ErrorCodeListenStart represents when a listener could not be started
	ErrorCodeListenStart ErrorCode = iota
)

// Error describes a listening error
type Error struct {
	Message string
	Code    ErrorCode
}

func (e Error) Error() string {
	return e.Message
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
