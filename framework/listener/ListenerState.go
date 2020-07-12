package listener

// State represents the state of the Listener
type State int

const (
	// LoadedState is the state when the Listener first loads
	LoadedState State = iota
	// RunningState is when the listener is running
	RunningState State = iota
	// StoppingState is when the listener is stopping
	StoppingState State = iota
)
