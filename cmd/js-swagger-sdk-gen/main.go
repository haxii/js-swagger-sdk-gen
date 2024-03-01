package main

import (
	"fmt"
	"os"
)

var (
	// Build of git, got by LDFLAGS on build
	Build = "dev"
	// Version of git, got by LDFLAGS on build
	Version = "dev"
)

func main() {
	var err error
	if _, err = parseOpt(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
