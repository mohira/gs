package main

import (
	"github.com/mohira/gs/cmd"
	"os"
)

func main() {
	cli := cmd.NewCLI(os.Stdin, os.Stdout, os.Stderr)

	os.Exit(cli.Run(os.Args))
}
