package main

import (
	"github.com/zachdeibert/protomux/config"
)

func main() {
	if _, err := config.LoadFile("example.conf"); err != nil {
		panic(err)
	}
}
