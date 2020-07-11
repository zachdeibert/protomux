package main

import (
	"fmt"

	"github.com/zachdeibert/protomux/config"
)

func main() {
	cfg, err := config.LoadFile("example.conf")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
