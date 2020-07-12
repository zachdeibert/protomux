package engine

import (
	"fmt"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeUnknownProtocol represents when a protocol is specified that does not have an implementaiton
	ErrorCodeUnknownProtocol ErrorCode = iota
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
