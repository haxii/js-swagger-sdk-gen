package main

import (
	"fmt"
	"os"
)

func main() {
	opt, err := parseOpt()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(opt)
}
