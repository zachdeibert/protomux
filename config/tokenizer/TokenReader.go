package tokenizer

import (
	"bufio"
	"io"
)

// TokenReader reads Tokens from a stream
type TokenReader struct {
	Scanner    *bufio.Scanner
	FileName   string
	Line       []byte
	LineNo     int
	Type       TokenType
	CharStart  int
	CharLen    int
	Escape     bool
	StringBuf  []byte
	CharNo     int
	WasNewline bool
}

// CreateTokenReader creates a new TokenReader
func CreateTokenReader(stream io.Reader, filename string) *TokenReader {
	return &TokenReader{
		Scanner:    bufio.NewScanner(stream),
		FileName:   filename,
		Line:       []byte{},
		LineNo:     -1,
		CharNo:     1,
		WasNewline: true,
	}
}

// Next reads the next token from the TokenReader
func (t *TokenReader) Next() (*Token, error) {
	for {
		if t.CharNo >= len(t.Line) {
			if !t.WasNewline {
				t.WasNewline = true
				return &Token{
					Type:  LineFeedToken,
					Value: "\n",
					Location: TokenLocation{
						FileName:  t.FileName,
						Line:      t.Line,
						LineNo:    t.LineNo,
						CharStart: len(t.Line),
						CharLen:   1,
					},
				}, nil
			}
			if !t.Scanner.Scan() {
				return nil, t.Scanner.Err()
			}
			t.LineNo++
			t.Line = t.Scanner.Bytes()
			t.CharStart = -1
			t.Escape = false
			t.CharNo = -1
		}
		t.CharNo++
		var c byte
		if t.CharNo == len(t.Line) {
			c = ' '
		} else {
			c = t.Line[t.CharNo]
		}
		if t.CharStart < 0 {
			switch t.Type = TokenLookup[c]; t.Type {
			case CommentToken:
				t.CharNo = len(t.Line)
				break
			case InvalidToken:
				return nil, ErrorInvalidChar(TokenLocation{
					FileName:  t.FileName,
					Line:      t.Line,
					LineNo:    t.LineNo,
					CharStart: t.CharNo,
					CharLen:   1,
				})
			case StringToken:
				t.StringBuf = []byte{}
				fallthrough
			case KeyToken, IntToken:
				t.CharStart = t.CharNo
				t.CharLen = 1
				break
			case OpenBraceSymbol, CloseBraceSymbol, OpenBracketSymbol, CloseBracketSymbol, ColonSymbol, DotSymbol, CommaSymbol:
				t.WasNewline = false
				return &Token{
					Type:  t.Type,
					Value: string([]byte{c}),
					Location: TokenLocation{
						FileName:  t.FileName,
						Line:      t.Line,
						LineNo:    t.LineNo,
						CharStart: t.CharNo,
						CharLen:   1,
					},
				}, nil
			case WhitespaceToken:
				break
			default:
				panic("Missing case")
			}
		} else {
			switch t.Type {
			case KeyToken, IntToken:
				if TokenLookup[c] == t.Type {
					t.CharLen++
				} else {
					start := t.CharStart
					t.CharStart = -1
					t.CharNo--
					t.WasNewline = false
					return &Token{
						Type:  t.Type,
						Value: string(t.Line[start : t.CharNo+1]),
						Location: TokenLocation{
							FileName:  t.FileName,
							Line:      t.Line,
							LineNo:    t.LineNo,
							CharStart: start,
							CharLen:   t.CharNo - start + 1,
						},
					}, nil
				}
				break
			case StringToken:
				t.CharLen++
				if t.Escape {
					var n byte
					switch c {
					case 'n':
						n = '\n'
						break
					case 't':
						n = '\t'
						break
					default:
						n = c
						break
					}
					t.StringBuf = append(t.StringBuf, n)
					t.Escape = false
				} else {
					switch c {
					case '\\':
						t.Escape = true
						break
					case t.Line[t.CharStart]:
						start := t.CharStart
						t.CharStart = -1
						t.WasNewline = false
						return &Token{
							Type:  StringToken,
							Value: string(t.StringBuf),
							Location: TokenLocation{
								FileName:  t.FileName,
								Line:      t.Line,
								LineNo:    t.LineNo,
								CharStart: start,
								CharLen:   t.CharLen,
							},
						}, nil
					default:
						t.StringBuf = append(t.StringBuf, c)
						break
					}
				}
				break
			default:
				panic("Missing case")
			}
		}
	}
}
