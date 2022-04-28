package cli

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Cli struct {
	programName string
	version     string
	rootCmd     *cobra.Command
}

func Init(programName string, version string) *Cli {
	viper.AutomaticEnv()

	myFigure := figure.NewFigure(programName, "", true)

	var rootCmd = &cobra.Command{
		Use:     programName,
		Short:   programName,
		Long:    myFigure.String(),
		Version: version,
	}

	return &Cli{
		programName: programName,
		version:     version,
		rootCmd:     rootCmd,
	}
}

func (cli *Cli) Run() error {
	return cli.rootCmd.Execute()
}
