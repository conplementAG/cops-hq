package hq

import (
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/conplementag/cops-hq/pkg/commands"
)

type hqContainer struct {
	Executor commands.Executor
	Cli      cli.Cli
}

func (hq *hqContainer) Run() error {
	return hq.Cli.Run()
}

func (hq *hqContainer) GetExecutor() commands.Executor {
	return hq.Executor
}

func (hq *hqContainer) GetCli() cli.Cli {
	return hq.Cli
}
