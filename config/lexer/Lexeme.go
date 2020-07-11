package lexer

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
)

// Lexeme represents a series of tokens with common meaning grouped together
type Lexeme struct {
	Type        LexemeType
	StringValue string
	IntValue    int
	Location    common.Location
}

func (l Lexeme) String() string {
	return fmt.Sprintf("{Lexeme Type=%d, StringValue='%s', IntValue=%d, At=%s}\n%s", l.Type, l.StringValue, l.IntValue, l.Location.ShortString(), l.Location)
}
