package commands

import "github.com/sirupsen/logrus"

type Executor interface {
	// Execute will run the given command, returning the stdout output and errors (if any).
	// In chatty mode, Execute shows output in both console and file. In quiet mode, output is piped only to the file.
	// Stdout is collected as return output, Stderr is not collected in the output.
	// Any errors are returned as err return value.
	Execute(command string) (output string, err error)

	// ExecuteWithProgressInfo is same as Execute, except an infinite progress bar is shown, signaling an async operation to
	// the user. The progress bar can be overridden with Viper parameter "silence-long-running-progress-indicators" - useful
	// for CI for example.
	ExecuteWithProgressInfo(command string) (output string, err error)

	// ExecuteSilent will run the given command, returning the stdout output and errors (if any).
	// No command output is shown on the console or logged to the file (irrelevant of the chatty / quiet setting). Can be
	// used for commands that are too verbose and clutter the output.
	// Stdout is collected as return output, Stderr is not collected in the output.
	// Any errors are returned as err return value.
	ExecuteSilent(command string) (output string, err error)

	// ExecuteTTY is a special executor for cases where the called command needs to detect it runs in a TTY session.
	// One example of such command is Docker. Commands executed via ExecuteTTY have their output shown on the console,
	// but the output is NOT saved to a log file. Chatty / Quiet settings have no effect on this method.
	ExecuteTTY(command string) error
}

// NewChatty creates a new Executor instance. Chatty executor outputs the command output to both file and console at
// the same time. Best suited for application IaC projects.
// Required dependencies are the log file name for command output, and the logging subsystem instance.
// Log file is required, because command output requires to be directly outputted to the file, without logging system
// interfering with formatting. Logging system is an explicit dependency, so that it is clear that the logging system
// need to bo be initialized first, before creating an Executor.
func NewChatty(logFileName string, logger *logrus.Logger) Executor {
	return &executor{
		logFileName: logFileName,
		logger:      logger,
		chatty:      true,
	}
}

// NewQuiet creates a new Executor instance. Quiet executor outputs the command output only to a file, console output is
// suppressed. Best suited for applications where the console should be kept clean and simple because of large amount of
// commands being executed.
// Required dependencies are the log file name for command output, and the logging subsystem instance.
// Log file is required, because command output requires to be directly outputted to the file, without logging system
// interfering with formatting. Logging system is an explicit dependency, so that it is clear that the logging system
// need to bo be initialized first, before creating an Executor.
func NewQuiet(logFileName string, logger *logrus.Logger) Executor {
	return &executor{
		logFileName: logFileName,
		logger:      logger,
		chatty:      false,
	}
}
