package minecraft

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/zachdeibert/protomux/framework"
)

// ProtocolInstance implementation for the Minecraft protocol
type ProtocolInstance struct {
	Action ActionProps
	Filter FilteringProps
}

// CreateProtocolInstance creates a new ProtocolInstance
func CreateProtocolInstance(action ActionProps, filter FilteringProps) *ProtocolInstance {
	return &ProtocolInstance{
		Action: action,
		Filter: filter,
	}
}

func e(err error) framework.ProtocolState {
	if err == io.EOF {
		return framework.ProtocolNeedsMoreData
	}
	return framework.ProtocolNotMatched
}

// Handle the protocol
func (p ProtocolInstance) Handle(conn net.Conn, state framework.HandlerState, waitGroup *sync.WaitGroup) framework.ProtocolState {
	stream := bufio.NewReader(conn)
	data, err := stream.Peek(1)
	if err != nil {
		return e(err)
	}
	if data[0] == 0xFE {
		// Before Netty rewrite
		// TODO
		return framework.ProtocolNotMatched
	}
	// After Netty rewrite
	return p.HandleNettyRewrite(conn, stream, state, waitGroup)
}
