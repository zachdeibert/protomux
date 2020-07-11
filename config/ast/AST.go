package ast

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/lexer"
)

// AST represents a section of the config file in an Abstract Syntax Tree
type AST struct {
	Blocks     []Block
	Parameters []Parameter
}

func parseAST(stream *lexer.LexemeReader) (*AST, *lexer.Lexeme, error) {
	var lexeme *lexer.Lexeme = nil
	var first *lexer.Lexeme = nil
	var block *Block = nil
	var param *Parameter = nil
	var err error
	ast := &AST{
		Blocks:     []Block{},
		Parameters: []Parameter{},
	}
	for {
		if lexeme, err = stream.Next(); err != nil {
			return nil, nil, err
		}
		if lexeme == nil {
			return ast, nil, nil
		}
		if first == nil {
			switch lexeme.Type {
			case lexer.LineFeedLexeme:
				break
			case lexer.BlockEndLexeme:
				return ast, lexeme, nil
			case lexer.KeywordLexeme:
				first = lexeme
				break
			default:
				return nil, nil, ErrorUnexpectedStartLexeme(*lexeme)
			}
		} else if lexeme.Type == lexer.BlockStartLexeme {
			if block, first, err = parseBlock(stream, *first, *lexeme); err != nil {
				return nil, nil, err
			}
			ast.Blocks = append(ast.Blocks, *block)
		} else {
			if param, first, err = parseParameter(stream, *first, *lexeme); err != nil {
				return nil, nil, err
			}
			ast.Parameters = append(ast.Parameters, *param)
		}
	}
}

// ParseAST parses an AST
func ParseAST(stream *lexer.LexemeReader) (*AST, error) {
	ast, lexeme, err := parseAST(stream)
	if err != nil {
		return nil, err
	}
	if lexeme != nil {
		return nil, ErrorUnexpectedStartLexeme(*lexeme)
	}
	return ast, nil
}

func (a AST) String() string {
	buf := strings.Builder{}
	buf.WriteString("AST")
	for i, block := range a.Blocks {
		var indent string
		var tree rune
		if i < len(a.Blocks)-1 || len(a.Parameters) > 0 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(block.String(), "\n", indent)))
	}
	for i, param := range a.Parameters {
		var indent string
		var tree rune
		if i < len(a.Parameters)-1 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(param.String(), "\n", indent)))
	}
	return buf.String()
}
