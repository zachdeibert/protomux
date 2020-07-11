package lexer

// LexemeType represents the type of a lexeme
type LexemeType int

const (
	// KeywordLexeme is a lexeme that represents a keyword
	KeywordLexeme LexemeType = iota
	// BlockStartLexeme represents the start of a block
	BlockStartLexeme LexemeType = iota
	// BlockEndLexeme represents the end of a block
	BlockEndLexeme LexemeType = iota
	// ArrayStartLexeme represents the start of an array
	ArrayStartLexeme LexemeType = iota
	// ArrayEndLexeme represents the end of an array
	ArrayEndLexeme LexemeType = iota
	// ArraySeparatorLexeme is between every element in an array
	ArraySeparatorLexeme LexemeType = iota
	// ConnectionLexeme contains an IP address or hostname and a port number
	ConnectionLexeme LexemeType = iota
	// StringLexeme is a literal string contained within ""s.
	StringLexeme LexemeType = iota
	// LineFeedLexeme represents the end of a line
	LineFeedLexeme LexemeType = iota
)
