package main

import (
	"os"

	"github.com/will2469/charites/internal/cli"
)

func main() {
	exitCode := cli.Execute(os.Args[1:])
	os.Exit(exitCode)
}
