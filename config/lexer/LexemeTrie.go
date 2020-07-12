package lexer

import (
	"net"
	"strconv"
	"strings"

	"github.com/zachdeibert/protomux/config/tokenizer"
)

var lexemeTrieKeyTokenDot = &LexemeTrieNode{
	Children: []*LexemeTrieNode{
		nil, // KeyToken (lexemeTrieKeyToken)
		nil, // IntToken
		nil, // StringToken
		nil, // OpenBraceSymbol
		nil, // CloseBraceSymbol
		nil, // OpenBracketSymbol
		nil, // CloseBracketSymbol
		nil, // ColonSymbol
		nil, // DotSymbol
		nil, // CommaSymbol
		nil, // LineFeedToken
	},
}

func init() {
	lexemeTrieKeyTokenDot.Children[tokenizer.KeyToken] = lexemeTrieKeyToken
}

var lexemeTrieKeyToken = &LexemeTrieNode{
	Children: []*LexemeTrieNode{
		nil, // KeyToken
		nil, // IntToken
		nil, // StringToken
		nil, // OpenBraceSymbol
		nil, // CloseBraceSymbol
		nil, // OpenBracketSymbol
		nil, // CloseBracketSymbol
		{ // ColonSymbol
			Children: []*LexemeTrieNode{
				nil, // KeyToken
				{ // IntToken
					Handler: func(t []tokenizer.Token) (*Lexeme, error) {
						port, err := strconv.ParseUint(t[len(t)-1].Value, 10, 16)
						if err != nil {
							// TODO wrap the error
							return nil, err
						}
						domains := make([]string, len(t)/2)
						for i := range domains {
							domains[i] = t[i*2].Value
						}
						return &Lexeme{
							Type:        ConnectionLexeme,
							StringValue: strings.Join(domains, "."),
							IntValue:    int(port),
						}, nil
					},
				},
				nil, // StringToken
				nil, // OpenBraceSymbol
				nil, // CloseBraceSymbol
				nil, // OpenBracketSymbol
				nil, // CloseBracketSymbol
				nil, // ColonSymbol
				nil, // DotSymbol
				nil, // CommaSymbol
				nil, // LineFeedToken
			},
		},
		lexemeTrieKeyTokenDot, // DotSymbol
		nil,                   // CommaSymbol
		nil,                   // LineFeedToken
	},
}

// LexemeTrie is a trie that is used for converting tokens into lexemes
var LexemeTrie = LexemeTrieNode{
	Children: []*LexemeTrieNode{
		{ // KeyToken
			Children: lexemeTrieKeyToken.Children,
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type:        KeywordLexeme,
					StringValue: t[0].Value,
				}, nil
			},
		},
		{ // IntToken
			Children: []*LexemeTrieNode{
				nil, // KeyToken
				nil, // IntToken
				nil, // StringToken
				nil, // OpenBraceSymbol
				nil, // CloseBraceSymbol
				nil, // OpenBracketSymbol
				nil, // CloseBracketSymbol
				nil, // ColonSymbol
				{ // DotSymbol
					Children: []*LexemeTrieNode{
						nil, // KeyToken
						{ // IntToken
							Children: []*LexemeTrieNode{
								nil, // KeyToken
								nil, // IntToken
								nil, // StringToken
								nil, // OpenBraceSymbol
								nil, // CloseBraceSymbol
								nil, // OpenBracketSymbol
								nil, // CloseBracketSymbol
								nil, // ColonSymbol
								{ // DotSymbol
									Children: []*LexemeTrieNode{
										nil, // KeyToken
										{ // IntToken
											Children: []*LexemeTrieNode{
												nil, // KeyToken
												nil, // IntToken
												nil, // StringToken
												nil, // OpenBraceSymbol
												nil, // CloseBraceSymbol
												nil, // OpenBracketSymbol
												nil, // CloseBracketSymbol
												nil, // ColonSymbol
												{ // DotSymbol
													Children: []*LexemeTrieNode{
														nil, // KeyToken
														{ // IntToken
															Children: []*LexemeTrieNode{
																nil, // KeyToken
																nil, // IntToken
																nil, // StringToken
																nil, // OpenBraceSymbol
																nil, // CloseBraceSymbol
																nil, // OpenBracketSymbol
																nil, // CloseBracketSymbol
																{ // ColonSymbol
																	Children: []*LexemeTrieNode{
																		nil, // KeyToken
																		{ // IntToken
																			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
																				a, err := strconv.ParseUint(t[0].Value, 10, 8)
																				if err != nil {
																					return nil, ErrorIntParse(t[0].Location, t[0].Value, err)
																				}
																				b, err := strconv.ParseUint(t[2].Value, 10, 8)
																				if err != nil {
																					return nil, ErrorIntParse(t[2].Location, t[2].Value, err)
																				}
																				c, err := strconv.ParseUint(t[4].Value, 10, 8)
																				if err != nil {
																					return nil, ErrorIntParse(t[4].Location, t[4].Value, err)
																				}
																				d, err := strconv.ParseUint(t[6].Value, 10, 8)
																				if err != nil {
																					return nil, ErrorIntParse(t[6].Location, t[6].Value, err)
																				}
																				port, err := strconv.ParseUint(t[8].Value, 10, 16)
																				if err != nil {
																					return nil, ErrorIntParse(t[8].Location, t[8].Value, err)
																				}
																				return &Lexeme{
																					Type:     ConnectionLexeme,
																					IntValue: int(port),
																					IPValue:  net.IPv4(byte(a), byte(b), byte(c), byte(d)),
																				}, nil
																			},
																		},
																		nil, // StringToken
																		nil, // OpenBraceSymbol
																		nil, // CloseBraceSymbol
																		nil, // OpenBracketSymbol
																		nil, // CloseBracketSymbol
																		nil, // ColonSymbol
																		nil, // DotSymbol
																		nil, // CommaSymbol
																		nil, // LineFeedToken
																	},
																},
																nil, // DotSymbol
																nil, // CommaSymbol
																nil, // LineFeedToken
															},
														},
														nil, // StringToken
														nil, // OpenBraceSymbol
														nil, // CloseBraceSymbol
														nil, // OpenBracketSymbol
														nil, // CloseBracketSymbol
														nil, // ColonSymbol
														nil, // DotSymbol
														nil, // CommaSymbol
														nil, // LineFeedToken
													},
												},
												nil, // CommaSymbol
												nil, // LineFeedToken
											},
										},
										nil, // StringToken
										nil, // OpenBraceSymbol
										nil, // CloseBraceSymbol
										nil, // OpenBracketSymbol
										nil, // CloseBracketSymbol
										nil, // ColonSymbol
										nil, // DotSymbol
										nil, // CommaSymbol
										nil, // LineFeedToken
									},
								},
								nil, // CommaSymbol
								nil, // LineFeedToken
							},
						},
						nil, // StringToken
						nil, // OpenBraceSymbol
						nil, // CloseBraceSymbol
						nil, // OpenBracketSymbol
						nil, // CloseBracketSymbol
						nil, // ColonSymbol
						nil, // DotSymbol
						nil, // CommaSymbol
						nil, // LineFeedToken
					},
				},
				nil, // CommaSymbol
				nil, // LineFeedToken
			},
		},
		{ // StringToken
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type:        StringLexeme,
					StringValue: t[0].Value,
				}, nil
			},
		},
		{ // OpenBraceSymbol
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type: BlockStartLexeme,
				}, nil
			},
		},
		{ // CloseBraceSymbol
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type: BlockEndLexeme,
				}, nil
			},
		},
		{ // OpenBracketSymbol
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type: ArrayStartLexeme,
				}, nil
			},
		},
		{ // CloseBracketSymbol
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type: ArrayEndLexeme,
				}, nil
			},
		},
		{ // ColonSymbol
			Children: []*LexemeTrieNode{
				nil, // KeyToken (lexemeTrieKeyToken)
				{ // IntToken
					Handler: func(t []tokenizer.Token) (*Lexeme, error) {
						port, err := strconv.ParseUint(t[1].Value, 10, 16)
						if err != nil {
							// TODO wrap the error
							return nil, err
						}
						return &Lexeme{
							Type:        ConnectionLexeme,
							StringValue: "0.0.0.0",
							IntValue:    int(port),
						}, nil
					},
				},
				nil, // StringToken
				nil, // OpenBraceSymbol
				nil, // CloseBraceSymbol
				nil, // OpenBracketSymbol
				nil, // CloseBracketSymbol
				nil, // ColonSymbol
				nil, // DotSymbol
				nil, // CommaSymbol
				nil, // LineFeedToken
			},
		},
		nil, // DotSymbol
		{ // CommaSymbol
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type: ArraySeparatorLexeme,
				}, nil
			},
		},
		{ // LineFeedToken
			Handler: func(t []tokenizer.Token) (*Lexeme, error) {
				return &Lexeme{
					Type: LineFeedLexeme,
				}, nil
			},
		},
	},
}
