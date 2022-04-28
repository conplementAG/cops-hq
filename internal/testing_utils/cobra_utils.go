package testing_utils

import (
	"bytes"
	"github.com/spf13/cobra"
)

func PrepareCommandForTesting(command *cobra.Command, args ...string) *bytes.Buffer {
	if args == nil {
		command.SetArgs(make([]string, 0)) // we have to explicitly set an empty array, because otherwise the os.Args will still be used
	} else {
		command.SetArgs(args)
	}

	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	return b
}
