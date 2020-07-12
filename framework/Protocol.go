package framework

import "github.com/zachdeibert/protomux/config"

// Protocol describes a loaded protocol
type Protocol interface {
	Configure(globals config.Parameters, remoteName string, remoteParams config.Parameters) (ProtocolInstance, error)
}
