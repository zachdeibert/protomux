package lexer

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
	"github.com/zachdeibert/protomux/config/tokenizer"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeIntParse represents an error caused by not being able to parse an integer
	ErrorCodeIntParse ErrorCode = iota
	// ErrorCodeInvalidToken represents when a token is received that should not be where it is
	ErrorCodeInvalidToken ErrorCode = iota
)

// Error describes a lexing error
type Error struct {
	Message  string
	Code     ErrorCode
	Location common.Location
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (at %s)\n%s", e.Message, e.Location.ShortString(), e.Location)
}

// ErrorIntParse creates a new ErrorIntParse error
func ErrorIntParse(location common.Location, val string, err error) error {
	return &Error{
		Message:  fmt.Sprintf("Unable to parse '%s' as an int: %s", val, err.Error()),
		Code:     ErrorCodeIntParse,
		Location: location,
	}
}

// ErrorInvalidToken creates a new ErrorInvalidToken error
func ErrorInvalidToken(token tokenizer.Token) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected token '%s'", token.Value),
		Code:     ErrorCodeInvalidToken,
		Location: token.Location,
	}
}
