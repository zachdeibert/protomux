package config

import (
	"fmt"
	"io"
	"os"

	"github.com/zachdeibert/protomux/config/tokenizer"
)

// Config contains the configuration for the protomux instance
type Config struct {
}

// Load the Config object from a stream
func Load(stream io.Reader) (*Config, error) {
	tokens, err := tokenizer.Tokenize(stream)
	if err != nil {
		return nil, err
	}
	for _, t := range tokens {
		fmt.Println(t)
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
	return Load(f)
}
