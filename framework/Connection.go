package framework

import "net"

// Connection represents a socket stream that is given to a protocol
type Connection interface {
	net.Conn
	RequireExclusive(priority int) error
}
