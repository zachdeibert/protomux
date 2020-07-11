package cmd

import (
	"fmt"
	"io"
)

// Help prints the help text to a given stream
func Help(stream io.Writer, flags []CommandFlag, progName string) {
	fmt.Fprintf(stream, "Usage: %s [options...]\n", progName)
	fmt.Fprintf(stream, "\n")
	fmt.Fprintf(stream, "Options:\n")
	firstCol := make([]string, len(flags))
	firstColSize := 0
	for i, opt := range flags {
		firstCol[i] = opt.FlagString()
		if l := len(firstCol[i]); l > firstColSize {
			firstColSize = l
		}
	}
	format := fmt.Sprintf("    %%-%ds    %%s\n", firstColSize)
	for i, opt := range flags {
		fmt.Fprintf(stream, format, firstCol[i], opt.DescString())
	}
}
