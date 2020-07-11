package config

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/ast"
	"github.com/zachdeibert/protomux/config/common"
)

// Parameters represents a list of parameters to an object
type Parameters struct {
	Strings     map[string][]string
	Connections map[string][]Connection
	Booleans    map[string][]bool
	Locations   map[string]common.Location
}

// ParseParameters parses an array of AST Parameters into a Parameters object
func ParseParameters(params []ast.Parameter) (*Parameters, error) {
	p := &Parameters{
		Strings:     map[string][]string{},
		Connections: map[string][]Connection{},
		Booleans:    map[string][]bool{},
		Locations:   map[string]common.Location{},
	}
	for _, param := range params {
		if first, ok := p.Locations[param.Name]; ok {
			return nil, ErrorDuplicateParam(param, first)
		}
		p.Locations[param.Name] = param.Location
		switch param.Type {
		case ast.StringParameter:
			vals := make([]string, len(param.Values))
			for i, v := range param.Values {
				vals[i] = v.(*ast.StringParameterData).Value
			}
			p.Strings[param.Name] = vals
			break
		case ast.ConnectionParameter:
			vals := make([]Connection, len(param.Values))
			for i, v := range param.Values {
				conn, err := ParseConnection(*v.(*ast.ConnectionParameterData))
				if err != nil {
					return nil, err
				}
				vals[i] = *conn
			}
			p.Connections[param.Name] = vals
			break
		case ast.BooleanParameter:
			vals := make([]bool, len(param.Values))
			for i, v := range param.Values {
				vals[i] = v.(*ast.BooleanParameterData).Value
			}
			p.Booleans[param.Name] = vals
			break
		default:
			panic("Missing case")
		}
	}
	return p, nil
}

func (p Parameters) String() string {
	buf := strings.Builder{}
	buf.WriteString("Parameters")
	i := 0
	for k, v := range p.Strings {
		var indent string
		var tree rune
		if i < len(p.Strings)-1 || len(p.Connections) > 0 || len(p.Booleans) > 0 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, k))
		for j, w := range v {
			if j < len(v)-1 {
				tree = '\u251C'
			} else {
				tree = '\u2514'
			}
			buf.WriteString(fmt.Sprintf("%s %c\u2500%s", indent, tree, w))
		}
		i++
	}
	i = 0
	for k, v := range p.Connections {
		var indent string
		var tree rune
		if i < len(p.Connections)-1 || len(p.Booleans) > 0 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, k))
		for j, w := range v {
			if j < len(v)-1 {
				tree = '\u251C'
			} else {
				tree = '\u2514'
			}
			buf.WriteString(fmt.Sprintf("%s %c\u2500%s", indent, tree, w))
		}
		i++
	}
	i = 0
	for k, v := range p.Booleans {
		var indent string
		var tree rune
		if i < len(p.Booleans)-1 {
			indent = "\n \u2502 "
			tree = '\u251C'
		} else {
			indent = "\n   "
			tree = '\u2514'
		}
		buf.WriteString(fmt.Sprintf("\n %c\u2500%s", tree, k))
		for j, w := range v {
			if j < len(v)-1 {
				tree = '\u251C'
			} else {
				tree = '\u2514'
			}
			buf.WriteString(fmt.Sprintf("%s %c\u2500%t", indent, tree, w))
		}
		i++
	}
	return buf.String()
}
