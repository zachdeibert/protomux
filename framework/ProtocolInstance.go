package framework

// ProtocolInstance represents a single, configurabled instance of a Protocol
type ProtocolInstance interface {
	Handle(conn Connection) error
}
