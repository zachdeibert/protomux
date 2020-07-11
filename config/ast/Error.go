package ast

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
	"github.com/zachdeibert/protomux/config/lexer"
)

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeUnexpectedStartLexeme represents when a lexeme was received at the beginning of a statement that cannot be the beginning of a statement
	ErrorCodeUnexpectedStartLexeme ErrorCode = iota
	// ErrorCodeUnexpectedBlockEOF represents when the end of the file was reached before a block finished
	ErrorCodeUnexpectedBlockEOF ErrorCode = iota
	// ErrorCodeUnexpectedParameterLexeme represents when an invalid lexeme was reached while parsing a parameter
	ErrorCodeUnexpectedParameterLexeme ErrorCode = iota
	// ErrorCodeSingleLexemeLine represents when only a single lexeme exists on a line being parsed by the AST
	ErrorCodeSingleLexemeLine ErrorCode = iota
	// ErrorCodeParameterTooLong represents when too many lexemes are on a parameter line
	ErrorCodeParameterTooLong ErrorCode = iota
	// ErrorCodeUnexpectedArrayLexeme represents when an invalid lexeme was reached while parsing an array
	ErrorCodeUnexpectedArrayLexeme ErrorCode = iota
	// ErrorCodeParameterArrayType represents when a single array contains multiple parameter types
	ErrorCodeParameterArrayType ErrorCode = iota
	// ErrorCodeUnexpectedArrayEOF represents when the end of file was reached before an array finished
	ErrorCodeUnexpectedArrayEOF ErrorCode = iota
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

// ErrorUnexpectedStartLexeme creates a new ErrorUnexpectedStartLexeme error
func ErrorUnexpectedStartLexeme(lexeme lexer.Lexeme) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected symbol '%s', expected keyword", lexeme.RawString()),
		Code:     ErrorCodeUnexpectedStartLexeme,
		Location: lexeme.Location,
	}
}

// ErrorUnexpectedBlockEOF creates a new ErrorUnexpectedBlockEOF error
func ErrorUnexpectedBlockEOF(block *Block) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected EOF while processing block '%s'", block.Name),
		Code:     ErrorCodeUnexpectedBlockEOF,
		Location: block.Location,
	}
}

// ErrorUnexpectedParameterLexeme creates a new ErrorUnexpectedParameterLexeme error
func ErrorUnexpectedParameterLexeme(lexeme lexer.Lexeme) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected symbol '%s' while parsing parameter value", lexeme.RawString()),
		Code:     ErrorCodeUnexpectedParameterLexeme,
		Location: lexeme.Location,
	}
}

// ErrorSingleLexemeLine creates a new ErrorSingleLexemeLine error
func ErrorSingleLexemeLine(lexeme lexer.Lexeme) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected line containing only one symbol '%s'", lexeme.RawString()),
		Code:     ErrorCodeSingleLexemeLine,
		Location: lexeme.Location,
	}
}

// ErrorParameterTooLong creates a new ErrorParameterTooLong error
func ErrorParameterTooLong(name string, overflow lexer.Lexeme) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected extra symbol '%s' after end of parameter '%s'", overflow.RawString(), name),
		Code:     ErrorCodeParameterTooLong,
		Location: overflow.Location,
	}
}

// ErrorUnexpectedArrayLexeme creates a new ErrorUnexpectedArrayLexeme error
func ErrorUnexpectedArrayLexeme(lexeme lexer.Lexeme) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected symbol '%s' while parsing parameter array value", lexeme.RawString()),
		Code:     ErrorCodeUnexpectedArrayLexeme,
		Location: lexeme.Location,
	}
}

// ErrorParameterArrayType creates a new ErrorParameterArrayType error
func ErrorParameterArrayType(lexeme lexer.Lexeme, typename string) error {
	return &Error{
		Message:  fmt.Sprintf("Expected entire array to be of type %s, but got an element '%s' of a different type", typename, lexeme.RawString()),
		Code:     ErrorCodeParameterArrayType,
		Location: lexeme.Location,
	}
}

// ErrorUnexpectedArrayEOF creates a new ErrorUnexpectedArrayEOF error
func ErrorUnexpectedArrayEOF(param *Parameter) error {
	return &Error{
		Message:  fmt.Sprintf("Unexpected EOF while processing array '%s'", param.Name),
		Code:     ErrorCodeUnexpectedArrayEOF,
		Location: param.Location,
	}
}
