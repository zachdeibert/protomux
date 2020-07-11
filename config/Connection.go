package config

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/ast"
)

// Connection represents a data type which contains information needed to open a socket
type Connection struct {
	Host string
	Port int
}

// ParseConnection parses an AST Connection parameter into a Connection
func ParseConnection(conn ast.ConnectionParameterData) (*Connection, error) {
	return &Connection{
		Host: conn.Host,
		Port: conn.Port,
	}, nil
}

func (c Connection) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
