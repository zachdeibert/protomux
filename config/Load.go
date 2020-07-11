package config

import (
	"io"
	"os"
	"path"

	"github.com/zachdeibert/protomux/config/ast"
	"github.com/zachdeibert/protomux/config/cmd"
	"github.com/zachdeibert/protomux/config/lexer"
	"github.com/zachdeibert/protomux/config/tokenizer"
)

// Load the Config object from a stream
func Load(stream io.Reader, filename string) (*Config, *ast.AST, error) {
	tokenStream := tokenizer.CreateTokenReader(stream, filename)
	lexemeStream := lexer.CreateLexemeReader(tokenStream)
	tree, err := ast.ParseAST(lexemeStream)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := ParseConfig(*tree)
	if err != nil {
		return nil, nil, err
	}
	for _, file := range cfg.Includes {
		_, overwrite, err := LoadFile(file)
		if err != nil {
			return nil, nil, err
		}
		tree.Merge(*overwrite)
	}
	cfg.Includes = []string{}
	if cfg, err = ParseConfig(*tree); err != nil {
		return nil, nil, err
	}
	return cfg, tree, nil
}

// LoadFile loads the Config object from a file
func LoadFile(filename string) (*Config, *ast.AST, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return Load(f, path.Base(filename))
}

// LoadCommandLine loads the Config object from command-line arguments
func LoadCommandLine(args []string) (*Config, *ast.AST, error) {
	tree, err := cmd.Parse(args)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := ParseConfig(*tree)
	if err != nil {
		return nil, nil, err
	}
	for _, file := range cfg.Includes {
		_, overwrite, err := LoadFile(file)
		if err != nil {
			return nil, nil, err
		}
		overwrite.Merge(*tree)
		tree = overwrite
	}
	cfg.Includes = []string{}
	if cfg, err = ParseConfig(*tree); err != nil {
		return nil, nil, err
	}
	return cfg, tree, nil
}
