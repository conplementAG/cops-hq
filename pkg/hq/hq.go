package hq

import (
	"github.com/conplementag/cops-hq/v2/pkg/cli"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/sirupsen/logrus"
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

	// GetLogrusLogger retrieves the currently initialized logrus logger
	GetLogrusLogger() *logrus.Logger

	// LoadEnvironmentConfigFile is a special version of LoadConfigFile in which the config file is expected in
	// the location 'config/<<environment_tag>>.yaml'. This command relies on the defined variable
	// 'environment-tag', available through Viper. Most common way to provide the 'environment-tag' is through cli
	// parameters, which are automatically bound to Viper.
	LoadEnvironmentConfigFile() error

	// LoadConfigFile loads the specified config file, which is expected to be saved encrypted (with sops) on disk.
	// You can generate the path using filepath.Join and the hq.ProjectBasePath variable. Sops is expected to be
	// available in PATH, and in correct version (use the CheckToolingDependencies method or call hq 'check-dependencies')
	LoadConfigFile(filePath string) error

	// GetRawConfigurationFile returns the currently loaded configuration decrypted as string. Loading configuration is
	// a procondition due to get raw configuration string
	GetRawConfigurationFile() (string, error)

	// CheckToolingDependencies can be called to check if installed tooling (Azure CLI, Terraform, Helm etc.) is of minimal
	// expected version for all of HQ functionality to work. It is highly recommended to call this method in your code, and fail
	// in case of errors.
	CheckToolingDependencies() error
}
