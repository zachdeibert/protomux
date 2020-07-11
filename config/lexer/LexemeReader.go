package lexer

import (
	"github.com/zachdeibert/protomux/config/common"
	"github.com/zachdeibert/protomux/config/tokenizer"
)

// LexemeReader reads Lexemes from a stream of Tokens
type LexemeReader struct {
	TokenStream *tokenizer.TokenReader
	Tokens      []tokenizer.Token
}

// CreateLexemeReader creates a new LexemeReader
func CreateLexemeReader(stream *tokenizer.TokenReader) *LexemeReader {
	return &LexemeReader{
		TokenStream: stream,
		Tokens:      []tokenizer.Token{},
	}
}

// Next reads the next lexeme from the LexemeReader
func (l *LexemeReader) Next() (*Lexeme, error) {
	for {
		token, err := l.TokenStream.Next()
		if err != nil {
			return nil, err
		}
		if token != nil {
			l.Tokens = append(l.Tokens, *token)
		}
		var lastValidNode *LexemeTrieNode = nil
		lastValidCount := -1
		node := &LexemeTrie
		broke := false
		for i, t := range l.Tokens {
			if node.Children != nil && node.Children[t.Type] != nil {
				node = node.Children[t.Type]
				if node.Handler != nil {
					lastValidNode = node
					lastValidCount = i
				}
			} else {
				broke = true
				if lastValidCount < 0 {
					return nil, ErrorInvalidToken(l.Tokens[0])
				}
			}
		}
		if lastValidCount >= 0 && (broke || (lastValidCount+1 == len(l.Tokens) && token == nil)) && (lastValidCount < len(l.Tokens) || token == nil) {
			lexeme, err := lastValidNode.Handler(l.Tokens[0 : lastValidCount+1])
			if err != nil {
				return nil, err
			}
			locations := make([]common.Location, lastValidCount+1)
			for i, t := range l.Tokens[0 : lastValidCount+1] {
				locations[i] = t.Location
			}
			l.Tokens = l.Tokens[lastValidCount+1:]
			if lexeme != nil {
				lexeme.Location = common.Merge(locations)
				return lexeme, nil
			}
		} else if token == nil {
			if len(l.Tokens) == 0 {
				return nil, nil
			}
			return nil, ErrorInvalidToken(l.Tokens[0])
		}
	}
}
