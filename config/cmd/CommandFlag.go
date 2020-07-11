package cmd

import "fmt"

// CommandFlag represents a flag that can be passed in on the command-line arguments
type CommandFlag struct {
	ShortFlag   string
	LongFlag    string
	Description string
	TreePath    string
	NumArgs     int
}

// FlagString returns the string that makes up the first column of the help text
func (f CommandFlag) FlagString() string {
	switch f.NumArgs {
	case 0:
		return fmt.Sprintf("-%s, --%s", f.ShortFlag, f.LongFlag)
	case 1:
		return fmt.Sprintf("-%s, --%s [value]", f.ShortFlag, f.LongFlag)
	case 2:
		return fmt.Sprintf("-%s, --%s [key] [value]", f.ShortFlag, f.LongFlag)
	default:
		panic("Missing case")
	}
}

// DescString returns the string that makes up the second column of the help text
func (f CommandFlag) DescString() string {
	return f.Description
}

func (f CommandFlag) String() string {
	return fmt.Sprintf("%s    %s", f.FlagString(), f.DescString())
}
