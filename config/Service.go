package config

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/ast"
)

// Service represents a set of addresses to listen on with how to multiplex the different protocols on those addresses
type Service struct {
	ListenAddresses []Connection
	Protocols       []Protocol
}

// ParseService parses a Block into a Service
func ParseService(block ast.Block) (*Service, error) {
	srv := &Service{
		Protocols: make([]Protocol, len(block.Children.Blocks)),
	}
	params, err := ParseParameters(block.Children.Parameters)
	if err != nil {
		return nil, err
	}
	var ok bool
	if srv.ListenAddresses, ok = params.Connections["listen"]; !ok {
		return nil, ErrorMissingParam("listen", "Service", block.Location)
	}
	for k := range params.Strings {
		return nil, ErrorUnknownParam(k, "Service", params.Locations[k])
	}
	for k := range params.Connections {
		if k != "listen" {
			return nil, ErrorUnknownParam(k, "Service", params.Locations[k])
		}
	}
	for k := range params.Booleans {
		return nil, ErrorUnknownParam(k, "Service", params.Locations[k])
	}
	protos := map[string]ast.Block{}
	for i, b := range block.Children.Blocks {
		proto, err := ParseProtocol(b)
		if err != nil {
			return nil, err
		}
		if first, ok := protos[proto.Name]; ok {
			return nil, ErrorDuplicateProtocol(first, b)
		}
		protos[proto.Name] = b
		srv.Protocols[i] = *proto
	}
	return srv, nil
}

func (s Service) String() string {
	buf := strings.Builder{}
	buf.WriteString("Service")
	for i, addr := range s.ListenAddresses {
		var indent string
		var tree rune
		if i < len(s.ListenAddresses)-1 || len(s.Protocols) > 0 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(addr.String(), "\n", indent)))
	}
	for i, proto := range s.Protocols {
		var indent string
		var tree rune
		if i < len(s.Protocols)-1 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, strings.ReplaceAll(proto.String(), "\n", indent)))
	}
	return buf.String()
}
