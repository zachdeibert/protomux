package minecraft

import (
	"bufio"

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

// Handle the protocol
func (p ProtocolInstance) Handle(conn framework.Connection) error {
	stream := bufio.NewReader(conn)
	data, err := stream.Peek(1)
	if err != nil {
		return err
	}
	if data[0] == 0xFE {
		// Before Netty rewrite
		// TODO
		return ErrorProtocol("Before Netty rewrite not supported")
	}
	// After Netty rewrite
	return p.HandleNettyRewrite(conn, stream)
}
