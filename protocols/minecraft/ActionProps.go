package minecraft

import "github.com/zachdeibert/protomux/config"

// ActionProps represents properties that are used for actions to run once a connection is received
type ActionProps struct {
	Remote *config.Connection
	MOTD   *string
	Kick   *string
}

// ParseActionProps parses the ActionProps from Parameters
func ParseActionProps(global config.Parameters, local config.Parameters) (*ActionProps, []string, []string, error) {
	props := &ActionProps{
		Remote: nil,
		MOTD:   nil,
		Kick:   nil,
	}
	globalUsed := []string{}
	localUsed := []string{}
	{
		var val []config.Connection = nil
		if v, ok := global.Connections["remote"]; ok {
			globalUsed = append(globalUsed, "remote")
			val = v
		}
		if v, ok := local.Connections["remote"]; ok {
			localUsed = append(localUsed, "remote")
			val = v
		}
		if val != nil {
			if len(val) > 1 {
				return nil, nil, nil, ErrorMultipleValues("remote")
			}
			props.Remote = &val[0]
		}
	}
	{
		var val []string = nil
		if v, ok := global.Strings["motd"]; ok {
			globalUsed = append(globalUsed, "motd")
			val = v
		}
		if v, ok := local.Strings["motd"]; ok {
			localUsed = append(localUsed, "motd")
			val = v
		}
		if val != nil {
			if len(val) > 1 {
				return nil, nil, nil, ErrorMultipleValues("motd")
			}
			props.MOTD = &val[0]
		}
	}
	{
		var val []string = nil
		if v, ok := global.Strings["kick"]; ok {
			globalUsed = append(globalUsed, "kick")
			val = v
		}
		if v, ok := local.Strings["kick"]; ok {
			localUsed = append(localUsed, "kick")
			val = v
		}
		if val != nil {
			if len(val) > 1 {
				return nil, nil, nil, ErrorMultipleValues("kick")
			}
			props.Kick = &val[0]
		}
	}
	if props.Remote != nil && props.Kick != nil {
		return nil, nil, nil, ErrorParameterRequirement("Both 'remote' and 'kick' may not be specified on the same server")
	}
	if props.Remote == nil && props.MOTD == nil {
		return nil, nil, nil, ErrorParameterRequirement("Either 'remote' or 'motd' must be specified on every server")
	}
	if props.Remote == nil && props.Kick == nil {
		return nil, nil, nil, ErrorParameterRequirement("Either 'remote' or 'kick' must be specified on every server")
	}
	return props, globalUsed, localUsed, nil
}
