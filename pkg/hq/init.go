package hq

import (
	"github.com/denisbiondic/cops-hq/pkg/cli"
	"github.com/denisbiondic/cops-hq/pkg/commands"
	"github.com/denisbiondic/cops-hq/pkg/logging"
)

type HQ struct {
	Executor *commands.Executor
	Cli      *cli.Cli
}

func Init(programName string, version string, logFileName string) *HQ {
	logging.Init(logFileName)
	exec := commands.NewExecutor(logFileName)
	cli := cli.Init(programName, version)

	return &HQ{
		Executor: exec,
		Cli:      cli,
	}
}
