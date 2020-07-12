package framework

// Protocols is a map of all supported protocols
var Protocols = map[string]Protocol{}

// RegisterProtocol registers a new protocol type with the framework
func RegisterProtocol(name string, val Protocol) {
	Protocols[name] = val
}
