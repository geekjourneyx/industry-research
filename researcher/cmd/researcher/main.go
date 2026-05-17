package main

import (
	"os"

	"github.com/geekjourneyx/researcher/internal/cli"
)

var Version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], Version, os.Stdout, os.Stderr))
}
