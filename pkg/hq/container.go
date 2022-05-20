package hq

import (
	"fmt"
	"github.com/conplementag/cops-hq/internal"
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
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

func (hq *hqContainer) LoadEnvironmentConfigFile() error {
	var err error

	configFilePath := filepath.Join(ProjectBasePath, "config", viper.GetString("environment-tag")+".yaml")

	// this should be kept as ExecuteSilent for security reasons, not to leak the whole config file in plaintext
	// to the log file!
	configFile, err := hq.Executor.ExecuteSilent("sops -d " + configFilePath)

	if err != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("error recieved while reading the config file: %w", err))
	}

	viper.SetConfigType("yaml")
	err = viper.MergeConfig(strings.NewReader(configFile))

	if err != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("error recieved while reading the config file: %w", err))
	}

	return nil
}
