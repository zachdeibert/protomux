package cmd

import (
	"fmt"
	"strings"

	"github.com/zachdeibert/protomux/config/common"
)

// ArgumentGroup represents a group of command-line arguments interpreted together
type ArgumentGroup struct {
	Arguments []string
	Locations []common.Location
}

func (g ArgumentGroup) String() string {
	buf := strings.Builder{}
	buf.WriteString("ArgumentGroup:")
	for i, arg := range g.Arguments {
		buf.WriteString(fmt.Sprintf("\n  Argument '%s'\n  ", arg))
		buf.WriteString(strings.ReplaceAll(g.Locations[i].String(), "\n", "\n  "))
	}
	return buf.String()
}
