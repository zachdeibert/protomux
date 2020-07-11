package tokenizer

import (
	"bufio"
	"io"
)

// Tokenize takes in a stream and returns a list of tokens contained in the string.
func Tokenize(stream io.Reader) ([]Token, error) {
	scan := bufio.NewScanner(stream)
	tokens := []Token{}
	for lineNo := 0; scan.Scan(); lineNo++ {
		line := scan.Bytes()
		next := Token{
			Line:      line,
			LineNo:    lineNo,
			CharStart: -1,
		}
		escape := false
		var stringBuf []byte = nil
		for charNo, c := range append(line, ' ') {
			isComment := false
			for consumed := false; !consumed; {
				if next.CharStart < 0 {
					consumed = true
					switch next.Type = TokenLookup[c]; next.Type {
					case CommentToken:
						isComment = true
						break
					case InvalidToken:
						return nil, ErrorInvalidChar(line, lineNo, charNo)
					case StringToken:
						stringBuf = []byte{}
						fallthrough
					case KeyToken, IntToken:
						next.CharStart = charNo
						next.CharLen = 1
						break
					case OpenBraceSymbol, CloseBraceSymbol, OpenBracketSymbol, CloseBracketSymbol, ColonSymbol, DotSymbol, CommaSymbol:
						tokens = append(tokens, Token{
							Type:      next.Type,
							Value:     string([]byte{c}),
							Line:      line,
							LineNo:    lineNo,
							CharStart: charNo,
							CharLen:   1,
						})
						break
					case WhitespaceToken:
						break
					default:
						panic("Missing case")
					}
				} else {
					switch next.Type {
					case KeyToken, IntToken:
						if TokenLookup[c] == next.Type {
							consumed = true
							next.CharLen++
						} else {
							next.Value = string(line[next.CharStart : next.CharStart+next.CharLen])
							tokens = append(tokens, next)
							next.CharStart = -1
						}
						break
					case StringToken:
						consumed = true
						next.CharLen++
						if escape {
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
							stringBuf = append(stringBuf, n)
							escape = false
						} else {
							switch c {
							case '\\':
								escape = true
								break
							case next.Line[next.CharStart]:
								next.Value = string(stringBuf)
								tokens = append(tokens, next)
								next.CharLen = -2
								break
							default:
								stringBuf = append(stringBuf, c)
								break
							}
						}
						break
					default:
						panic("Missing case")
					}
				}
			}
			if isComment {
				break
			}
		}
		if len(tokens) == 0 || tokens[len(tokens)-1].Type != LineFeedToken {
			tokens = append(tokens, Token{
				Type:      LineFeedToken,
				Value:     "\n",
				Line:      line,
				LineNo:    lineNo,
				CharStart: len(line),
				CharLen:   1,
			})
		}
	}
	return tokens, scan.Err()
}
