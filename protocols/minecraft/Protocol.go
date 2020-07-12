package minecraft

import (
	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
)

// Protocol implementation for the Minecraft protocol
type Protocol struct {
}

// Configure the protocol
func (p Protocol) Configure(globals config.Parameters, remoteName string, remoteParams config.Parameters) (framework.ProtocolInstance, error) {
	action, actionGlobals, actionLocals, err := ParseActionProps(globals, remoteParams)
	if err != nil {
		return nil, err
	}
	filter, filterGlobals, filterLocals, err := ParseFilteringProps(globals, remoteParams)
	if err != nil {
		return nil, err
	}
	usedGlobals := map[string]interface{}{}
	for _, v := range actionGlobals {
		usedGlobals[v] = nil
	}
	for _, v := range filterGlobals {
		usedGlobals[v] = nil
	}
	for k, v := range globals.Locations {
		if _, ok := usedGlobals[k]; !ok {
			return nil, ErrorUnrecognizedParameter(k, v)
		}
	}
	usedLocals := map[string]interface{}{}
	for _, v := range actionLocals {
		usedLocals[v] = nil
	}
	for _, v := range filterLocals {
		usedLocals[v] = nil
	}
	for k, v := range remoteParams.Locations {
		if _, ok := usedLocals[k]; !ok {
			return nil, ErrorUnrecognizedParameter(k, v)
		}
	}
	switch remoteName {
	case "server":
		if filter.IsEmpty() {
			return nil, ErrorParameterRequirement("There must be at least one filter requirement set")
		}
		break
	case "default":
		if !filter.IsEmpty() {
			return nil, ErrorParameterRequirement("The default server cannot have any filter requirement set")
		}
		break
	default:
		return nil, ErrorUnknownRemoteType(remoteName)
	}
	return CreateProtocolInstance(*action, *filter), nil
}

func init() {
	framework.RegisterProtocol("minecraft", &Protocol{})
}
