package hq

import (
	"github.com/denisbiondic/cops-hq/pkg/cli"
	"github.com/denisbiondic/cops-hq/pkg/commands"
	"github.com/denisbiondic/cops-hq/pkg/logging"
)

// HQ is an easy one-stop setup for typical IaC projects. Don't forget to call the Run() method after you complete
// setting up (e.g. all CLI commands added to HQ.Cli)
type HQ struct {
	Executor commands.Executor
	Cli      cli.Cli
}

// New creates a new HQ instance, configuring internally used modules for usage. Keep the created HQ instance and
// avoid re-instantiation since some setup steps might have global impacts (like logging setup).
// It will create a chatty executor, piping all commands and outputs to both the console and the file (e.g. like
// in shell scripts)
func New(programName string, version string, logFileName string) *HQ {
	logger := logging.Init(logFileName)
	exec := commands.NewChatty(logFileName, logger)
	cli := cli.New(programName, version)

	return &HQ{
		Executor: exec,
		Cli:      cli,
	}
}

// NewQuiet creates a new HQ instance, configuring internally used modules for usage. Keep the created HQ instance and
// avoid re-instantiation since some setup steps might have global impacts (like logging setup).
// Quiet HQ will create a quiet executor, piping all commands and outputs to the log file, but the console will be
// kept clean. If needed, like in CI, a viper flag "verbose" can be used to override this behavior.
func NewQuiet(programName string, version string, logFileName string) *HQ {
	logger := logging.Init(logFileName)
	exec := commands.NewQuiet(logFileName, logger)
	cli := cli.New(programName, version)

	return &HQ{
		Executor: exec,
		Cli:      cli,
	}
}

// Run starts the HQ CLI parsing functionality
func (hq *HQ) Run() error {
	return hq.Cli.Run()
}
