package config

import (
	"fmt"

	"github.com/zachdeibert/protomux/config/ast"
)

// Remote represents a remote server that a Protocol handles
type Remote struct {
	Name       string
	Parameters Parameters
}

// ParseRemote parses a Block into a Remote
func ParseRemote(block ast.Block) (*Remote, error) {
	if len(block.Children.Blocks) > 0 {
		return nil, ErrorRemoteBlocks(block.Children.Blocks)
	}
	params, err := ParseParameters(block.Children.Parameters)
	if err != nil {
		return nil, err
	}
	return &Remote{
		Name:       block.Name,
		Parameters: *params,
	}, nil
}

func (r Remote) String() string {
	return fmt.Sprintf("Remote %s\n%s", r.Name, r.Parameters.String()[len("Parameters\n"):])
}
