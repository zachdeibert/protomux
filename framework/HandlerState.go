package framework

// HandlerState represents the state the Listener is currently in
type HandlerState int

const (
	// ProbingState means multiple protocols are still possible, so ProtocolInstance.Handle should not open external connections
	ProbingState HandlerState = iota
	// HandoffState means the socket is ready to be handed off, so ProtocolInstance.Handle can open external connections
	HandoffState HandlerState = iota
	// StoppingState means the listener is being stopped
	StoppingState HandlerState = iota
)
