package cmd

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"sync"
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

	switch flags.Args()[0] {
	case "th":
		thCmd := flag.NewFlagSet("th", flag.ExitOnError)
		thUserToken := thCmd.String("token", os.Getenv("USER_TOKEN"), "token")
		thChannelId := thCmd.String("channel", os.Getenv("CHANNEL_ID"), "channel")

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
	case "dl":
		dlCmd := flag.NewFlagSet("dl", flag.ExitOnError)
		dlUserToken := dlCmd.String("token", os.Getenv("USER_TOKEN"), "token")
		dlChannelId := dlCmd.String("channel", os.Getenv("CHANNEL_ID"), "channel")
		dlTs := dlCmd.String("ts", "", "ts")

		if err := dlCmd.Parse(args[2:]); err != nil {
			return ExitCodeParseError
		}

		slackFiles, err := FetchSlackFiles(*dlUserToken, *dlChannelId, *dlTs)
		if err != nil {
			fmt.Fprintf(c.errStream, "err: %s", err.Error())
			return ExitCodeSomeError
		}

		var wg sync.WaitGroup

		// TODO: エラー処理が謎(存在しないディレクトリを使えばエラー起こせる)
		for _, file := range slackFiles {
			wg.Add(1)
			fmt.Fprintf(c.outStream, "%s\n", file.Name)

			go func(name, url string) error {
				defer wg.Done()

				savePath := "out/" + name
				f, err := os.Create(savePath)
				if err != nil {
					fmt.Fprintf(c.errStream, "os.Create error: %s\n", err)
					return err
				}

				response, err := http.Get(url)
				if err != nil {
					fmt.Fprintf(c.errStream, "http.Get error: %s\n", err)
					return err
				}

				if err := response.Write(f); err != nil {
					fmt.Fprintf(c.errStream, "response.Write error: %s\n", err)
					return err
				}

				return nil
			}(file.Name, file.UrlPrivateDownload)
		}

		wg.Wait()

	}

	return ExitCodeOK
}
