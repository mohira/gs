package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/mohira/gs/cmd"
	"strings"
	"testing"
)

func TestRun_versionFlag(t *testing.T) {
	t.Parallel()

	inStream := new(bytes.Buffer)
	outStream := new(bytes.Buffer)
	errStream := new(bytes.Buffer)

	cli := cmd.NewCLI(inStream, outStream, errStream)

	command := "gs -version"
	args := strings.Split(command, " ")

	code := cli.Run(args)
	if code != cmd.ExitCodeOK {
		t.Errorf("%q exits %d, want %d\n", command, code, cmd.ExitCodeOK)
	}

	output := outStream.String()
	expected := fmt.Sprintf("gs version v%s\n", cmd.Version)

	if diff := cmp.Diff(output, expected); diff != "" {
		t.Errorf("%s", diff)
	}
}
