package commands

import (
	"github.com/denisbiondic/cops-hq/internal/commands"
	"github.com/denisbiondic/cops-hq/internal/logging"
	"io"
	"os"
	"strings"
	"sync"
)

type Executor struct {
	logFileName string
}

// Execute will run the given command, returning the stdout output and errors (if any).
// Stdout is collected as return output, Stderr is shown on console and logged, but not collected in the output.
// Any errors are returned as err return value.
func (e *Executor) Execute(command string) (output string, err error) {
	cmd := commands.Create(command)

	// Create cmdStdOut, cmdStdErr streams of type io.Reader
	cmdStdOut, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	cmdStdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	// Start command
	err = cmd.Start()

	if err != nil {
		return "", err
	}

	var stdoutCollector strings.Builder

	writerStdout := io.MultiWriter(os.Stdout, logging.NewLogFileAppender(e.logFileName), &stdoutCollector)
	writerStderr := io.MultiWriter(os.Stderr, logging.NewLogFileAppender(e.logFileName))

	// Pipe command output to all the writers, and also wait for all the writing to be done
	// so that we can parse the results (only after the io.Copy is done will our stdoutCollector be filled!
	var multiWritingSteps sync.WaitGroup
	multiWritingSteps.Add(2)

	go func() {
		io.Copy(writerStdout, cmdStdOut)
		multiWritingSteps.Done()
	}()

	go func() {
		io.Copy(writerStderr, cmdStdErr)
		multiWritingSteps.Done()
	}()

	err = cmd.Wait()
	multiWritingSteps.Wait()

	// some consoles always append a \n at the end, but this is safe to be removed
	cleanedStringOutput := strings.TrimSuffix(stdoutCollector.String(), "\n")
	return cleanedStringOutput, err
}

// ExecuteTTY is a special executor for cases where the called command needs to detect it runs in a TTY session.
// One example of such command is Docker. Commands executed via ExecuteTTY have their output shown on the console,
// but the output is not saved to a log file.
func (e *Executor) ExecuteTTY(command string) error {
	cmd := commands.Create(command)

	// only the direct pipe to os.Std* will work for TTY, using io.MultiWriter like in
	// the standard Execute() did not work that executing process recognizes it is in TTY session...
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
