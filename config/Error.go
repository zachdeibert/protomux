package config

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/ast"
	"github.com/zachdeibert/protomux/config/common"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeRemoteBlocks represents when a Remote AST contains block children
	ErrorCodeRemoteBlocks ErrorCode = iota
	// ErrorCodeUnknownParam represents when an unknown parameter is given to a Service or Config
	ErrorCodeUnknownParam ErrorCode = iota
	// ErrorCodeDuplicateProtocol represents when the same Protocol is specified twice for the same Service
	ErrorCodeDuplicateProtocol ErrorCode = iota
	// ErrorCodeMissingParam represents when a required parameter is missing on a Service or Config
	ErrorCodeMissingParam ErrorCode = iota
	// ErrorCodeDuplicateParam represents when the same parameter is specified twice for the same object
	ErrorCodeDuplicateParam ErrorCode = iota
	// ErrorCodeUnknownBlock represents when an unknonw block is given to a Config
	ErrorCodeUnknownBlock ErrorCode = iota
)

// Error describes a parsing error
type Error struct {
	Message   string
	Code      ErrorCode
	Locations []common.Location
}

func (e Error) Error() string {
	if len(e.Locations) == 0 {
		return e.Message
	}
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("%s (at %s)", e.Message, e.Locations[0].ShortString()))
	for _, loc := range e.Locations {
		buf.WriteString(fmt.Sprintf("\n%s", loc))
	}
	return buf.String()
}

// ErrorRemoteBlocks creates a new ErrorRemoteBlocks error
func ErrorRemoteBlocks(blocks []ast.Block) error {
	locs := make([]common.Location, len(blocks))
	for i, b := range blocks {
		locs[i] = b.Location
	}
	return &Error{
		Message:   "Remote block should not contain any children blocks",
		Code:      ErrorCodeRemoteBlocks,
		Locations: locs,
	}
}

// ErrorUnknownParam creates a new ErrorUnknownParam error
func ErrorUnknownParam(name, objType string, location common.Location) error {
	return &Error{
		Message:   fmt.Sprintf("Unrecognized parameter '%s' on %s object", name, objType),
		Code:      ErrorCodeUnknownParam,
		Locations: []common.Location{location},
	}
}

// ErrorDuplicateProtocol creates a new ErrorDuplicateProtocol error
func ErrorDuplicateProtocol(first ast.Block, second ast.Block) error {
	return &Error{
		Message: fmt.Sprintf("Duplicate protocol '%s' specified", first.Name),
		Code:    ErrorCodeDuplicateProtocol,
		Locations: []common.Location{
			first.Location,
			second.Location,
		},
	}
}

// ErrorMissingParam creates a new ErrorMissingParam error
func ErrorMissingParam(name, objType string, location common.Location) error {
	return &Error{
		Message:   fmt.Sprintf("Required parameter '%s' was not given on %s object", name, objType),
		Code:      ErrorCodeMissingParam,
		Locations: []common.Location{location},
	}
}

// ErrorDuplicateParam creates a new ErrorDuplicateParam error
func ErrorDuplicateParam(param ast.Parameter, first common.Location) error {
	return &Error{
		Message: fmt.Sprintf("Duplicate value without array for parameter '%s'", param.Name),
		Code:    ErrorCodeDuplicateParam,
		Locations: []common.Location{
			first,
			param.Location,
		},
	}
}

// ErrorUnknownBlock creates a new ErrorUnknownBlock error
func ErrorUnknownBlock(name, objType string, location common.Location) error {
	return &Error{
		Message:   fmt.Sprintf("Unrecognized block '%s' on %s object", name, objType),
		Code:      ErrorCodeUnknownBlock,
		Locations: []common.Location{location},
	}
}
