package ast

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/common"
	"github.com/zachdeibert/protomux/config/lexer"
)

// Block represents a block in the AST
type Block struct {
	Name     string
	Location common.Location
	Children AST
}

func parseBlock(stream *lexer.LexemeReader, first lexer.Lexeme, second lexer.Lexeme) (*Block, *lexer.Lexeme, error) {
	block := &Block{
		Name:     first.StringValue,
		Location: common.Merge([]common.Location{first.Location, second.Location}),
	}
	ast, lexeme, err := parseAST(stream)
	if err != nil {
		return nil, nil, err
	}
	if lexeme == nil {
		return nil, nil, ErrorUnexpectedBlockEOF(block)
	}
	block.Children = *ast
	return block, nil, nil
}

func (b Block) String() string {
	return fmt.Sprintf("Block '%s':\n%s", b.Name, b.Children.String()[len("AST\n"):])
}
