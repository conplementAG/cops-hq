package cli

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Cli interface {
	// AddBaseCommand adds a command on the root (base) level of the command tree. Use the *Command return value to add
	// additional subcommands. Parameters: use (command mnemonic), shortInfo (short info shown in help),
	// longDescription (long description shown in help), runFunction (function to be run when command is invoked)
	AddBaseCommand(use string, shortInfo string, longDescription string, runFunction func()) Command

	// Run starts the cli, parsing the given os.Args and executing the matching command
	Run() error

	// GetRootCommand returns the root top level command, directly as cobra.Command which is the library used
	// under the hood.
	GetRootCommand() *cobra.Command
}

// New creates a new Cli instance
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
