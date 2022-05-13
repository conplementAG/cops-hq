package hq

import (
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/conplementag/cops-hq/pkg/logging"
)

// HQ is an easy one-stop setup for typical IaC projects. Don't forget to call the Run() method after you complete
// setting up (e.g. all CLI commands added to HQ.Cli). Consider this object similar to an IoC container, which can be
// used to retrieve main dependencies for other objects, such as the command executor or the CLI.
type HQ interface {
	// Run starts the HQ CLI parsing functionality
	Run() error

	// GetExecutor retrieves the currently configured executor
	GetExecutor() commands.Executor

	// GetCli retrieves the current cli instance
	GetCli() cli.Cli

	// CheckToolingDependencies can be called to check if installed tooling (Azure CLI, Terraform, Helm etc.) is of minimal
	// expected version for all of HQ functionality to work. It is highly recommended to call this method in your code, and fail
	// in case of errors.
	CheckToolingDependencies() error
}

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

func addInbuiltHqCliCommands(cli cli.Cli, container *hqContainer) {
	hqBaseCommand := cli.AddBaseCommand("hq", "in-build HQ command group", "Command predefined by cops-hq.", nil)
	hqBaseCommand.AddCommand("check-dependencies", "Checks the installed tools and their versions", "Use this command "+
		"to check the installed versions of tools such as azure-cli, kubectl etc.", func() {
		container.CheckToolingDependencies()
	})
}
