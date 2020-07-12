package config

import (
	"fmt"
	"net"

	"github.com/zachdeibert/protomux/config/ast"
)

// Connection represents a data type which contains information needed to open a socket
type Connection struct {
	Host string
	IP   net.IP
	Port int
}

// ParseConnection parses an AST Connection parameter into a Connection
func ParseConnection(conn ast.ConnectionParameterData) (*Connection, error) {
	return &Connection{
		Host: conn.Host,
		IP:   conn.IP,
		Port: conn.Port,
	}, nil
}

func (c Connection) String() string {
	if len(c.Host) > 0 {
		return fmt.Sprintf("%s:%d", c.Host, c.Port)
	}
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}
