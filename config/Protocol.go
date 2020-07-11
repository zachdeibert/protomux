package config

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/ast"
)

// Protocol represents a single protocol configuration for a Service
type Protocol struct {
	Name       string
	Parameters Parameters
	Remotes    []Remote
}

// ParseProtocol parses a Block into a Protocol
func ParseProtocol(block ast.Block) (*Protocol, error) {
	proto := &Protocol{
		Name:    block.Name,
		Remotes: make([]Remote, len(block.Children.Blocks)),
	}
	params, err := ParseParameters(block.Children.Parameters)
	if err != nil {
		return nil, err
	}
	proto.Parameters = *params
	for i, b := range block.Children.Blocks {
		remote, err := ParseRemote(b)
		if err != nil {
			return nil, err
		}
		proto.Remotes[i] = *remote
	}
	return proto, nil
}

func (p Protocol) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("Protocol %s", p.Name))
	var indent string
	var tree rune
	if len(p.Remotes) > 0 {
		indent = "\n \u2502 "
		tree = '\u251C'
	} else {
		indent = "\n   "
		tree = '\u2514'
	}
	buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(p.Parameters.String(), "\n", indent)))
	for i, remote := range p.Remotes {
		if i < len(p.Remotes)-1 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(remote.String(), "\n", indent)))
	}
	return buf.String()
}
