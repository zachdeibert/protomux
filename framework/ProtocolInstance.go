package framework

import (
	"net"
	"sync"
)

// ProtocolInstance represents a single, configurabled instance of a Protocol
type ProtocolInstance interface {
	Handle(conn net.Conn, state HandlerState, waitGroup *sync.WaitGroup) ProtocolState
}
