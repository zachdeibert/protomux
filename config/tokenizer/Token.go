package tokenizer

import "fmt"

// Token represents a single token from the config file
type Token struct {
	Type     TokenType
	Value    string
	Location TokenLocation
}

func (t Token) String() string {
	return fmt.Sprintf("{Token Type=%d, Value='%s', At=%s}\n%s", t.Type, t.Value, t.Location.ShortString(), t.Location)
}
