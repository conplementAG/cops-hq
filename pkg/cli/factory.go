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

	// these flags should be available on every command
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Set to override the executor to always "+
		"show the output on the console. Useful in CI scenarios, if using quiet mode.")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("silence-long-running-progress-indicators", "s", false, "Set to "+
		"silence long running operation indicator. Useful for CI.")

	viper.BindPFlag("silence-long-running-progress-indicators", rootCmd.PersistentFlags().Lookup("silence-long-running-progress-indicators"))

	return &cli{
		programName: programName,
		version:     version,
		rootCmd:     rootCmd,
	}
}
