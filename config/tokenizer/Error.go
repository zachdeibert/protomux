package tokenizer

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeInvalidChar represents an error caused by an invalid character
	ErrorCodeInvalidChar ErrorCode = iota
)

// Error describes a tokenizing error
type Error struct {
	Message  string
	Code     ErrorCode
	Location common.Location
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (at %s)\n%s", e.Message, e.Location.ShortString(), e.Location)
}

// ErrorInvalidChar creates a new ErrorInvalidChar error
func ErrorInvalidChar(location common.Location) error {
	return &Error{
		Message:  fmt.Sprintf("Invalid character '%c'", location.Line[location.CharStart]),
		Code:     ErrorCodeInvalidChar,
		Location: location,
	}
}
