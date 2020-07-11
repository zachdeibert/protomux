package tokenizer

import "fmt"

// ErrorCode describes a specific error
type ErrorCode int

const (
	// ErrorCodeInvalidChar represents an error caused by an invalid character
	ErrorCodeInvalidChar ErrorCode = iota
)

// Error describes a tokenizing error
type Error struct {
	Message string
	Code    ErrorCode
	Line    []byte
	LineNo  int
	CharNo  int
}

func (e Error) Error() string {
	buf := make([]byte, len(e.Message)+1+len(e.Line)+1+e.CharNo+1)
	copy(buf, []byte(e.Message))
	bufI := len(e.Message)
	buf[bufI] = '\n'
	bufI++
	copy(buf[bufI:], e.Line)
	bufI += len(e.Line)
	buf[bufI] = '\n'
	bufI++
	for i := 0; i < e.CharNo; i++ {
		buf[bufI] = ' '
		bufI++
	}
	buf[bufI] = '^'
	return string(buf)
}

// ErrorInvalidChar creates a new ErrorInvalidChar error
func ErrorInvalidChar(line []byte, lineNo, charNo int) error {
	return &Error{
		Message: fmt.Sprintf("Invalid character '%c'", line[charNo]),
		Code:    ErrorCodeInvalidChar,
		Line:    line,
		LineNo:  lineNo,
		CharNo:  charNo,
	}
}
