package cmd

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"os"
)

const (
	ExitCodeOK         = 0
	ExitCodeParseError = 1
	ExitCodeSomeError  = 2
)

var helpText = `Usage: myxargs [OPTION]... COMMAND [INITIAL-ARGS]...
A simple tool, written in Go.

  -help       display this help and exit
  -version    output version information and exit
`

type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
}

func NewCLI(inStream io.Reader, outStream io.Writer, errStream io.Writer) *CLI {
	return &CLI{inStream: inStream, outStream: outStream, errStream: errStream}
}

func (c *CLI) Run(args []string) int {
	err := godotenv.Load()
	if err != nil {
		return ExitCodeSomeError
	}

	var version bool

	flags := flag.NewFlagSet("gs", flag.ContinueOnError)
	flags.BoolVar(&version, "version", false, "")

	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}

	if len(args) == 1 {
		flags.Usage()
		return ExitCodeSomeError
	}
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseError
	}

	if version {
		_, err := fmt.Fprintln(c.outStream, OutputVersion())
		if err != nil {
			return ExitCodeSomeError
		}

		return ExitCodeOK
	}

	thCmd := flag.NewFlagSet("th", flag.ExitOnError)
	thUserToken := thCmd.String("token", os.Getenv("USER_TOKEN"), "token")
	thChannelId := thCmd.String("channel", os.Getenv("CHANNEL_ID"), "channel")

	switch flags.Args()[0] {
	case "th":
		if err := thCmd.Parse(args[2:]); err != nil {
			return ExitCodeParseError
		}

		slackThreads, err := FetchSlackThreads(*thUserToken, *thChannelId)

		if err != nil {
			if _, err := fmt.Fprint(c.errStream, err); err != nil {
				return ExitCodeSomeError
			}
			return ExitCodeSomeError
		}

		f, err := os.Create("threads.csv")
		if err != nil {
			return ExitCodeSomeError
		}
		defer f.Close()

		w := csv.NewWriter(f)
		for _, thread := range slackThreads {
			record := []string{thread.Ts, thread.FirstLine()}

			if err := w.Write(record); err != nil {
				return ExitCodeSomeError
			}
		}
		w.Flush()

	}

	return ExitCodeOK
}
