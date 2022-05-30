package hq

import (
	"github.com/conplementag/cops-hq/internal"
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

// ProjectBasePath simply points to root of the Go project, which should always be two levels above
// the currently executed directory (which per convention, should always be cmd/project_name
var ProjectBasePath = filepath.Join(".", "../", "../")

type hqContainer struct {
	Executor commands.Executor
	Cli      cli.Cli
	Logger   *logrus.Logger
}

func (hq *hqContainer) Run() error {
	err := hq.Cli.Run()
	return internal.ReturnErrorOrPanic(err)
}

func (hq *hqContainer) GetExecutor() commands.Executor {
	return hq.Executor
}

func (hq *hqContainer) GetCli() cli.Cli {
	return hq.Cli
}

func (hq *hqContainer) GetLogrusLogger() *logrus.Logger {
	return hq.Logger
}
