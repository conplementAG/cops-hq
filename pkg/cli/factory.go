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

	// OnInitialize sets the passed function to be run when each command is called. Consider this like a global initializer
	// hook.
	OnInitialize(initFunction func())

	SetDefaultCommand(command string)
}

// New creates a new Cli instance
func New(programName string, version string) Cli {
	// This method will mark the resolution in viper to use env variables as priority when searching for values.
	// Note: This does not load all the env variables at this point in code, so it is ok to hev this call in the beginning!
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

	rootCmd.PersistentFlags().BoolP("silence-long-running-progress-indicators", "", false, "Set to "+
		"silence long running operation indicator. Useful for CI.")

	viper.BindPFlag("silence-long-running-progress-indicators", rootCmd.PersistentFlags().Lookup("silence-long-running-progress-indicators"))

	return &cli{
		programName: programName,
		version:     version,
		rootCmd:     rootCmd,
	}
}
