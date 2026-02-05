package main

import (
	"os"

	"github.com/tasuku43/gionx/internal/cli"
)

var version = "dev"

func main() {
	c := cli.New(os.Stdout, os.Stderr)
	c.Version = version
	os.Exit(c.Run(os.Args[1:]))
}
