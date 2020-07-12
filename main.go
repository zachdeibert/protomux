package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/config/cmd"
	"github.com/zachdeibert/protomux/framework/engine"

	_ "github.com/zachdeibert/protomux/protocols/minecraft"
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
	eng, err := engine.CreateEngine(*cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	eng.Start()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	fmt.Println("ProtoMux started.")
	<-c
	eng.Stop()
	fmt.Println("ProtoMux stopped.")
}
