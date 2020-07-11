package config

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/zachdeibert/protomux/config/lexer"
	"github.com/zachdeibert/protomux/config/tokenizer"
)

// Config contains the configuration for the protomux instance
type Config struct {
}

// Load the Config object from a stream
func Load(stream io.Reader, filename string) (*Config, error) {
	tokenStream := tokenizer.CreateTokenReader(stream, filename)
	lexemeStream := lexer.CreateLexemeReader(tokenStream)
	var lexeme *lexer.Lexeme
	var err error
	for lexeme, err = lexemeStream.Next(); lexeme != nil; lexeme, err = lexemeStream.Next() {
		fmt.Println(lexeme)
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// LoadFile loads the Config object from a file
func LoadFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f, path.Base(filename))
}
