package cmd

import (
	"strings"

	"github.com/zachdeibert/protomux/config/common"
)

// Locations of the command-line arguments
func Locations(args []string) []common.Location {
	charStart := 0
	line := []byte(strings.Join(args, " "))
	locs := make([]common.Location, len(args))
	for i, arg := range args {
		locs[i] = common.Location{
			FileName:  "cmdline",
			Line:      line,
			LineNo:    i,
			CharStart: charStart,
			CharLen:   len(arg),
		}
		charStart += len(arg) + 1
	}
	return locs
}
