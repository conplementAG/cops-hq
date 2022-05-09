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
func New(programName string, version string, logFileName string) *HQ {
	logger := logging.Init(logFileName)
	exec := commands.NewChatty(logFileName, logger)
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
