package framework

// ProtocolState represents the state of the protocol handshake
type ProtocolState int

const (
	// ProtocolNeedsMoreData means there is not yet enough information to judge if the protocol is a match or not
	ProtocolNeedsMoreData ProtocolState = iota
	// ProtocolMatched means that the protocol is a match
	ProtocolMatched ProtocolState = iota
	// ProtocolNotMatched means that the protocol is not a match
	ProtocolNotMatched ProtocolState = iota
)
