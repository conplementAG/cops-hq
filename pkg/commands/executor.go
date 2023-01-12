package commands

import (
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Executor interface {
	// Execute will run the given command, returning the stdout output and errors (if any).
	// In chatty mode, Execute shows output in both console and file. In quiet mode, output is piped only to the file.
	// Stdout is collected as return output, Stderr is not collected in the output.
	// Any errors are returned as err return value.
	Execute(command string) (output string, err error)

	// ExecuteCmd is same as Execute, except you can provide the os/exec command directly. Useful to avoid Executor escaping logic,
	// in rare cases where the command does not follow the usual --argument value semantics.
	ExecuteCmd(cmd *exec.Cmd) (output string, err error)

	// ExecuteWithProgressInfo is same as Execute, except an infinite progress bar is shown, signaling an async operation to
	// the user. The progress bar can be overridden with Viper parameter "silence-long-running-progress-indicators" - useful
	// for CI for example.
	ExecuteWithProgressInfo(command string) (output string, err error)

	// ExecuteCmdWithProgressInfo is same as ExecuteWithProgressInfo, except you can provide the os/exec command directly. Useful to avoid Executor escaping logic,
	// in rare cases where the command does not follow the usual --argument value semantics.
	ExecuteCmdWithProgressInfo(cmd *exec.Cmd) (output string, err error)

	// ExecuteSilent will run the given command, returning the stdout output and errors (if any).
	// No command output is shown on the console or logged to the file (irrelevant of the chatty / quiet setting). Can be
	// used for commands that are too verbose and clutter the output.
	// Stdout is collected as return output, Stderr is not collected in the output.
	// Any errors are returned as err return value.
	ExecuteSilent(command string) (output string, err error)

	// ExecuteCmdSilent is same as ExecuteSilent, except you can provide the os/exec command directly. Useful to avoid Executor escaping logic,
	// in rare cases where the command does not follow the usual --argument value semantics.
	ExecuteCmdSilent(cmd *exec.Cmd) (output string, err error)

	// ExecuteTTY is a special executor for cases where the called command needs to detect it runs in a TTY session.
	// One example of such command is Docker. Commands executed via ExecuteTTY have their output shown on the console,
	// but the output is NOT saved to a log file. Chatty / Quiet settings have no effect on this method.
	ExecuteTTY(command string) error

	// ExecuteCmdTTY is same as ExecuteTTY, except you can provide the os/exec command directly. Useful to avoid Executor escaping logic,
	// in rare cases where the command does not follow the usual --argument value semantics.
	ExecuteCmdTTY(cmd *exec.Cmd) error

	// AskUserToConfirm pauses the execution, and awaits for user to confirm (by either typing yes, Y or y).
	// Parameter displayMessage can be used to show a message on the screen.
	AskUserToConfirm(displayMessage string) bool

	// AskUserToConfirmWithKeyword pauses the execution, and awaits for user to confirm with the requested keyword.
	// Parameter displayMessage can be used to show a message on the screen.
	AskUserToConfirmWithKeyword(displayMessage string, keyword string) bool
}

type executor struct {
	logFileName string
	logger      *logrus.Logger
	chatty      bool

	stdin io.Reader
}

func (e *executor) Execute(command string) (output string, err error) {
	return e.executeString(command, false)
}

func (e *executor) ExecuteCmd(cmd *exec.Cmd) (output string, err error) {
	return e.executeCmd(cmd, false)
}

func (e *executor) ExecuteWithProgressInfo(command string) (output string, err error) {
	if !viper.GetBool("silence-long-running-progress-indicators") {
		spinner := createAndStartSpinner()
		defer spinner.Stop()
	}

	return e.executeString(command, false)
}

func (e *executor) ExecuteCmdWithProgressInfo(cmd *exec.Cmd) (output string, err error) {
	if !viper.GetBool("silence-long-running-progress-indicators") {
		spinner := createAndStartSpinner()
		defer spinner.Stop()
	}

	return e.executeCmd(cmd, false)
}

func (e *executor) ExecuteSilent(command string) (output string, err error) {
	return e.executeString(command, true)
}

func (e *executor) ExecuteCmdSilent(cmd *exec.Cmd) (output string, err error) {
	return e.executeCmd(cmd, true)
}

func (e *executor) ExecuteTTY(command string) error {
	e.logger.Info("[Command] " + command)
	cmd := Create(command)

	// only the direct pipe to os.Std* will work for TTY, using io.MultiWriter like in
	// the standard Execute() did not work that executing process recognizes it is in TTY session...
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return internal.ReturnErrorOrPanic(err)
}

func (e *executor) ExecuteCmdTTY(cmd *exec.Cmd) error {
	e.logger.Info("[Command] " + cmd.String())

	// only the direct pipe to os.Std* will work for TTY, using io.MultiWriter like in
	// the standard Execute() did not work that executing process recognizes it is in TTY session...
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return internal.ReturnErrorOrPanic(err)
}

func (e *executor) executeString(command string, silent bool) (output string, err error) {
	if !silent {
		commandStartMessage := "[Command] " + command

		if e.chatty {
			e.logger.Info(commandStartMessage)
		} else {
			logging.NewLogFileAppender(e.logFileName).Write([]byte(commandStartMessage))
		}
	}

	cmd := Create(command)
	return e.execute(cmd, silent)
}

func (e *executor) executeCmd(cmd *exec.Cmd, silent bool) (output string, err error) {
	if !silent {
		commandStartMessage := "[Command (via os/exec)] " + cmd.String()

		if e.chatty {
			e.logger.Info(commandStartMessage)
		} else {
			logging.NewLogFileAppender(e.logFileName).Write([]byte(commandStartMessage))
		}
	}

	return e.execute(cmd, silent)
}

func (e *executor) execute(cmd *exec.Cmd, silent bool) (output string, err error) {
	// Logic of conditional pipe-ing of command outputs here:
	// 1. We "capture" the sources of command output from the command itself, by assigning the pipes to local variables.
	//    These variables are of type io.Reader.
	cmdStdOut, pipeError := cmd.StdoutPipe()
	if pipeError != nil {
		return "", internal.ReturnErrorOrPanic(pipeError)
	}

	cmdStdErr, pipeError := cmd.StderrPipe()
	if pipeError != nil {
		return "", internal.ReturnErrorOrPanic(pipeError)
	}

	// Start command
	commandStartError := cmd.Start()

	if commandStartError != nil {
		return "", internal.ReturnErrorOrPanic(commandStartError)
	}

	// 2. We create a composite io.Writer consisting of multiple sinks. Depending on the configuration, these writers
	//    either write to "nothing" (discard), or they write to a file / console / buffer to collect the output, etc.
	stdoutWriter := ioutil.Discard
	stderrWriter := ioutil.Discard
	logFileWriter := ioutil.Discard
	var stdoutCollector strings.Builder
	var stderrCollector strings.Builder

	if !silent {
		logFileWriter = logging.NewLogFileAppender(e.logFileName)

		if e.chatty || viper.GetBool("verbose") {
			stdoutWriter = os.Stdout
			stderrWriter = os.Stderr
		}
	}

	writerStdout := io.MultiWriter(stdoutWriter, logFileWriter, &stdoutCollector)
	writerStderr := io.MultiWriter(stderrWriter, logFileWriter, &stderrCollector)

	// 3. We connect the reader(s) to writer(s) via io.Copy, executed asynchronously. We wait until both are completed.
	// Note: only after the io.Copy is done will our stdoutCollector be filled, so we have to wait!
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

	commandError := cmd.Wait()
	multiWritingSteps.Wait()

	// some consoles always append a \n at the end, but this is safe to be removed
	cleanedStringOutput := strings.TrimSuffix(stdoutCollector.String(), "\n")

	// composite error will be used to return stderr in case an error occurs, otherwise
	// stderr will be ignored completely (unless verbose mode is used, or chatty executor)
	var compositeError error

	if commandError != nil {
		compositeError = fmt.Errorf("%w; Stderr stream: "+stderrCollector.String(), commandError)
	}

	return cleanedStringOutput, internal.ReturnErrorOrPanic(compositeError)
}

func createAndStartSpinner() *spinner.Spinner {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Prefix = "Please wait "
	spinner.Color("green", "bold")
	spinner.Start()
	return spinner
}

func (e *executor) AskUserToConfirm(displayMessage string) bool {
	// Asks the user for confirmation, returns true if the user inputs yes, otherwise false
	logrus.Info(displayMessage + " [yes|no]")

	confirmation := bufio.NewScanner(e.stdin)
	confirmation.Scan()

	acceptedValues := []string{"yes", "YES", "y", "Y"}

	for okValueIndex := range acceptedValues {
		okText := acceptedValues[okValueIndex]
		text := confirmation.Text()

		if text == okText {
			return true
		}
	}

	return false
}

func (e *executor) AskUserToConfirmWithKeyword(displayMessage string, keyword string) bool {
	logrus.Info(displayMessage)

	confirmation := bufio.NewScanner(e.stdin)
	confirmation.Scan()

	text := confirmation.Text()

	if text == keyword {
		return true
	}

	return false
}

func (e *executor) OverrideStdIn(override io.Reader) {
	e.stdin = override
}
