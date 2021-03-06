package lexer

import (
	"fmt"
	"net"

	"github.com/zachdeibert/protomux/config/common"
)

// Lexeme represents a series of tokens with common meaning grouped together
type Lexeme struct {
	Type        LexemeType
	StringValue string
	IntValue    int
	IPValue     net.IP
	Location    common.Location
}

// RawString gets the raw value of the lexeme, as it appears in the file
func (l Lexeme) RawString() string {
	if l.Type == LineFeedLexeme {
		return "\\n"
	}
	return string(l.Location.Line[l.Location.CharStart : l.Location.CharStart+l.Location.CharLen])
}

func (l Lexeme) String() string {
	return fmt.Sprintf("{Lexeme Type=%d, StringValue='%s', IntValue=%d, IPValue=%s, At=%s}\n%s", l.Type, l.StringValue, l.IntValue, l.IPValue, l.Location.ShortString(), l.Location)
}
