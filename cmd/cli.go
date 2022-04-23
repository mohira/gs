package cmd

import (
	"flag"
	"fmt"
	"io"
)

const (
	ExitCodeOK         = 0
	ExitCodeParseError = 1
	ExitCodeSomeError  = 2
)

type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
}

func NewCLI(inStream io.Reader, outStream io.Writer, errStream io.Writer) *CLI {
	return &CLI{inStream: inStream, outStream: outStream, errStream: errStream}
}

func (c *CLI) Run(args []string) int {
	var version bool

	flags := flag.NewFlagSet("gs", flag.ContinueOnError)
	flags.BoolVar(&version, "version", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		_, err := fmt.Fprint(c.errStream, "flags.Parseに失敗したよ")
		if err != nil {
			return ExitCodeSomeError
		}
		return ExitCodeParseError
	}

	if version {
		_, err := fmt.Fprintln(c.outStream, OutputVersion())
		if err != nil {
			return ExitCodeSomeError
		}

		return ExitCodeOK
	}

	return ExitCodeOK
}
