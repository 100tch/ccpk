package main

import (
	"fmt"
	"os"

	"ccpk/internal/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "[error] use -help")
		os.Exit(1)
	}

	os.Exit(cli.Run(args))
}
