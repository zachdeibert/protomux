package minecraft

import "github.com/zachdeibert/protomux/config"

// FilteringProps represents properties that are used for protocol filtering
type FilteringProps struct {
	VersionName   []string
	Version       []Version
	ServerAddress []config.Connection
}

// ParseFilteringProps parses the FilteringProps from Parameters
func ParseFilteringProps(global config.Parameters, local config.Parameters) (*FilteringProps, []string, []string, error) {
	props := &FilteringProps{
		VersionName:   []string{},
		Version:       []Version{},
		ServerAddress: []config.Connection{},
	}
	globalUsed := []string{}
	localUsed := []string{}
	{
		var val []string = nil
		if v, ok := global.Strings["version"]; ok {
			globalUsed = append(globalUsed, "version")
			val = v
		}
		if v, ok := local.Strings["version"]; ok {
			localUsed = append(localUsed, "version")
			val = v
		}
		if val != nil {
			props.VersionName = val
			props.Version = make([]Version, len(val))
			for i, v := range val {
				p, err := ParseVersion(v)
				if err != nil {
					return nil, nil, nil, err
				}
				props.Version[i] = p
			}
		}
	}
	{
		var val []config.Connection = nil
		if v, ok := global.Connections["inbound"]; ok {
			globalUsed = append(globalUsed, "inbound")
			val = v
		}
		if v, ok := local.Connections["inbound"]; ok {
			localUsed = append(localUsed, "inbound")
			val = v
		}
		if val != nil {
			props.ServerAddress = val
		}
	}
	return props, globalUsed, localUsed, nil
}

// Check determines if this FilteringProps matches the filter
func (p FilteringProps) Check(filter FilteringProps) bool {
	if len(filter.Version) != 0 {
		if len(p.Version) == 0 {
			return false
		}
		found := false
		for _, v := range filter.Version {
			if p.Version[0] == v {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(filter.ServerAddress) != 0 {
		if len(p.ServerAddress) == 0 {
			return false
		}
		found := false
		for _, v := range filter.ServerAddress {
			if p.ServerAddress[0].Host == v.Host && (p.ServerAddress[0].IP.Equal(v.IP) || v.Host != "") && p.ServerAddress[0].Port == v.Port {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// IsEmpty determines if there are no filters set
func (p FilteringProps) IsEmpty() bool {
	return len(p.Version) == 0 &&
		len(p.ServerAddress) == 0
}
