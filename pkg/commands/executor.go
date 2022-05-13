package commands

import (
	"bufio"
	"github.com/briandowns/spinner"
	"github.com/conplementag/cops-hq/internal/commands"
	"github.com/conplementag/cops-hq/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type executor struct {
	logFileName string
	logger      *logrus.Logger
	chatty      bool

	stdin io.Reader
}

func (e *executor) Execute(command string) (output string, err error) {
	return e.execute(command, false)
}

func (e *executor) ExecuteWithProgressInfo(command string) (output string, err error) {
	if !viper.GetBool("silence-long-running-progress-indicators") {
		spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		spinner.Prefix = "Please wait "
		spinner.Color("green", "bold")
		spinner.Start()
		defer spinner.Stop()
	}

	return e.execute(command, false)
}

func (e *executor) ExecuteSilent(command string) (output string, err error) {
	return e.execute(command, true)
}

func (e *executor) ExecuteTTY(command string) error {
	e.logger.Info("[Command] " + command)
	cmd := commands.Create(command)

	// only the direct pipe to os.Std* will work for TTY, using io.MultiWriter like in
	// the standard Execute() did not work that executing process recognizes it is in TTY session...
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *executor) execute(command string, silent bool) (output string, err error) {
	if !silent {
		commandStartMessage := "[Command] " + command

		if e.chatty {
			e.logger.Info(commandStartMessage)
		} else {
			logging.NewLogFileAppender(e.logFileName).Write([]byte(commandStartMessage))
		}
	}

	cmd := commands.Create(command)

	// Logic of conditional pipe-ing of command outputs here:
	// 1. We "capture" the sources of command output from the command itself, by assigning the pipes to local variables.
	//    These variables are of type io.Reader.
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

	// 2. We create a composite io.Writer consisting of multiple sinks. Depending on the configuration, these writers
	//    either write to "nothing" (discard), or they write to a file / console / buffer to collect the output, etc.
	stdoutWriter := ioutil.Discard
	stderrWriter := ioutil.Discard
	logFileWriter := ioutil.Discard
	var stdoutCollector strings.Builder

	if !silent {
		logFileWriter = logging.NewLogFileAppender(e.logFileName)

		if e.chatty || viper.GetBool("verbose") {
			stdoutWriter = os.Stdout
			stderrWriter = os.Stderr
		}
	}

	writerStdout := io.MultiWriter(stdoutWriter, logFileWriter, &stdoutCollector)
	writerStderr := io.MultiWriter(stderrWriter, logFileWriter)

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

	err = cmd.Wait()
	multiWritingSteps.Wait()

	// some consoles always append a \n at the end, but this is safe to be removed
	cleanedStringOutput := strings.TrimSuffix(stdoutCollector.String(), "\n")
	return cleanedStringOutput, err
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

func (e *executor) OverrideStdIn(override io.Reader) {
	e.stdin = override
}
