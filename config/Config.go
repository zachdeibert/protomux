package config

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/ast"
)

// Config contains the configuration for the protomux instance
type Config struct {
	Services []Service
}

// ParseConfig parses an AST into a Config
func ParseConfig(tree ast.AST) (*Config, error) {
	params, err := ParseParameters(tree.Parameters)
	if err != nil {
		return nil, err
	}
	for k := range params.Strings {
		return nil, ErrorUnknownParam(k, "Config", params.Locations[k])
	}
	for k := range params.Connections {
		if k != "listen" {
			return nil, ErrorUnknownParam(k, "Config", params.Locations[k])
		}
	}
	cfg := &Config{
		Services: make([]Service, len(tree.Blocks)),
	}
	for i, b := range tree.Blocks {
		if b.Name != "service" {
			return nil, ErrorUnknownBlock(b.Name, "Config", b.Location)
		}
		srv, err := ParseService(b)
		if err != nil {
			return nil, err
		}
		cfg.Services[i] = *srv
	}
	return cfg, nil
}

func (c Config) String() string {
	buf := strings.Builder{}
	buf.WriteString("Config")
	for i, srv := range c.Services {
		var indent string
		var tree rune
		if i < len(c.Services)-1 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(srv.String(), "\n", indent)))
	}
	return buf.String()
}
