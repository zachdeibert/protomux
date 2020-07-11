package config

import "github.com/zachdeibert/protomux/config/cmd"

// CommandFlags is a list of all supported command-line flags
var CommandFlags = []cmd.CommandFlag{
	{
		ShortFlag:   "h",
		LongFlag:    "help",
		Description: "Prints this help text to the console and exits",
		TreePath:    "/showHelp",
		NumArgs:     0,
	},
	{
		ShortFlag:   "o",
		LongFlag:    "option",
		Description: "Sets an arbitrary option in the configuration",
		TreePath:    "",
		NumArgs:     2,
	},
	{
		ShortFlag:   "c",
		LongFlag:    "config-file",
		Description: "Loads a configuration file",
		TreePath:    "/include",
		NumArgs:     1,
	},
}
