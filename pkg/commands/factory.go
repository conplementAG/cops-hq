package commands

import (
	"github.com/sirupsen/logrus"
	"os"
)

// NewChatty creates a new Executor instance. Chatty executor outputs the command output to both file and console at
// the same time. Best suited for application IaC projects.
// Required dependencies are the log file name for command output, and the logging subsystem instance.
// Log file is required, because command output requires to be directly outputted to the file, without logging system
// interfering with formatting. Logging system is an explicit dependency, so that it is clear that the logging system
// need to bo be initialized first, before creating an Executor.
func NewChatty(logFileName string, logger *logrus.Logger) Executor {
	return create(logFileName, logger, true)
}

// NewQuiet creates a new Executor instance. Quiet executor outputs the command output only to a file, console output is
// suppressed. Best suited for applications where the console should be kept clean and simple because of large amount of
// commands being executed.
// Required dependencies are the log file name for command output, and the logging subsystem instance.
// Quiet Executor supports a viper flag "verbose", so that the command output can be piped to the console if needed (useful
// in CI scenarios).
// Log file is required, because command output requires to be directly outputted to the file, without logging system
// interfering with formatting. Logging system is an explicit dependency, so that it is clear that the logging system
// need to bo be initialized first, before creating an Executor.
func NewQuiet(logFileName string, logger *logrus.Logger) Executor {
	return create(logFileName, logger, false)
}

func create(logFileName string, logger *logrus.Logger, chatty bool) Executor {
	e := &executor{
		logFileName: logFileName,
		logger:      logger,
		chatty:      chatty,
	}

	e.stdin = os.Stdin

	return e
}
