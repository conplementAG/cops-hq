package cli

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Cli interface {
	// AddBaseCommand adds a top level command on the root level.
	AddBaseCommand(use string, shortInfo string, longDescription string, runFunction func()) *Command

	// Run runs the Cli itself, starting the parsing of the provided os.Args
	Run() error

	// GetRootCommand returns the top-level Cobra command
	GetRootCommand() *cobra.Command
}

func New(programName string, version string) Cli {
	viper.AutomaticEnv()

	myFigure := figure.NewFigure(programName, "", true)

	var rootCmd = &cobra.Command{
		Use:     programName,
		Short:   programName,
		Long:    myFigure.String(),
		Version: version,
	}

	return &cli{
		programName: programName,
		version:     version,
		rootCmd:     rootCmd,
	}
}
