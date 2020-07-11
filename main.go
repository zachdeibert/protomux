package main

import (
	"fmt"
	"os"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/config/cmd"
)

func main() {
	cfg, _, err := config.LoadCommandLine(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, ok := err.(*cmd.Error); ok {
			fmt.Fprintln(os.Stderr)
			cmd.Help(os.Stderr, config.CommandFlags, os.Args[0])
		}
		os.Exit(1)
	}
	if cfg.ShowHelp {
		cmd.Help(os.Stdout, config.CommandFlags, os.Args[0])
		os.Exit(0)
	}
	fmt.Println(cfg)
}
