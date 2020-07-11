package lexer

import "github.com/zachdeibert/protomux/config/tokenizer"

// LexemeTrieNode represents a node in the trie used for converting tokens into Lexemes
type LexemeTrieNode struct {
	Children []*LexemeTrieNode
	Handler  func([]tokenizer.Token) (*Lexeme, error)
}
