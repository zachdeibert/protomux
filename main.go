package main

import (
	"fmt"
	"os"

	"github.com/zachdeibert/protomux/config"
)

func main() {
	cfg, _, err := config.LoadCommandLine(os.Args[1:])
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
