package hq

import (
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/conplementag/cops-hq/pkg/logging"
)

// New creates a new HQ instance, configuring internally used modules for usage. Keep the created HQ instance and
// avoid re-instantiation since some setup steps might have global impacts (like logging setup).
// It will create a chatty executor, piping all commands and outputs to both the console and the file (e.g. like
// in shell scripts)
func New(programName string, version string, logFileName string) HQ {
	return create(programName, version, logFileName, false)
}

// NewQuiet creates a new HQ instance, configuring internally used modules for usage. Keep the created HQ instance and
// avoid re-instantiation since some setup steps might have global impacts (like logging setup).
// Quiet HQ will create a quiet executor, piping all commands and outputs to the log file, but the console will be
// kept clean. If needed, like in CI, a viper flag "verbose" can be used to override this behavior.
func NewQuiet(programName string, version string, logFileName string) HQ {
	return create(programName, version, logFileName, true)
}

func create(programName string, version string, logFileName string, quiet bool) HQ {
	logger := logging.Init(logFileName)
	cli := cli.New(programName, version)

	var exec commands.Executor

	if quiet {
		exec = commands.NewQuiet(logFileName, logger)
	} else {
		exec = commands.NewChatty(logFileName, logger)
	}

	container := &hqContainer{
		Executor: exec,
		Cli:      cli,
	}

	addInbuiltHqCliCommands(cli, container)
	return container
}
