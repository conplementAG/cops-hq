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
	return create(programName, version, &HqOptions{
		Quiet:       false,
		LogFileName: logFileName,
	})
}

// NewQuiet creates a new HQ instance, configuring internally used modules for usage. Keep the created HQ instance and
// avoid re-instantiation since some setup steps might have global impacts (like logging setup).
// Quiet HQ will create a quiet executor, piping all commands and outputs to the log file, but the console will be
// kept clean. If needed, like in CI, a viper flag "verbose" can be used to override this behavior.
func NewQuiet(programName string, version string, logFileName string) HQ {
	return create(programName, version, &HqOptions{
		Quiet:       true,
		LogFileName: logFileName,
	})
}

// NewCustom creates a new HQ instance, configuring internally used modules for usage. Keep the created HQ instance and
// avoid re-instantiation since some setup steps might have global impacts (like logging setup).
// Custom HQ lets you override most of the behaviours and HQ setup options
func NewCustom(programName string, version string, options *HqOptions) HQ {
	return create(programName, version, options)
}

func create(programName string, version string, options *HqOptions) HQ {
	logger := logging.Init(options.LogFileName)
	cli := cli.New(programName, version)

	var exec commands.Executor

	if options.Quiet {
		exec = commands.NewQuiet(options.LogFileName, logger)
	} else {
		exec = commands.NewChatty(options.LogFileName, logger)
	}

	container := &hqContainer{
		Executor: exec,
		Cli:      cli,
		Logger:   logger,
	}

	addInbuiltHqCliCommands(cli, container)
	return container
}
