package cli

import (
	"github.com/conplementag/cops-hq/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cli struct {
	programName string
	version     string
	rootCmd     *cobra.Command
}

func (cli *cli) AddBaseCommand(use string, shortInfo string, longDescription string, runFunction func()) Command {
	cw := &commandWrapper{}

	command := &cobra.Command{
		Use:   use,
		Short: shortInfo,
		Long:  longDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if runFunction != nil {
				runFunction()
			} else {
				cmd.Help()
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			// we have to map the viper parameters on runtime, when the command is executing, to prevent
			// overwriting of viper mappings in case multiple commands have the same named parameters
			for _, p := range cw.parameters {
				viper.BindPFlag(p, cmd.Flags().Lookup(p))
			}

			// there is no need to traverse the parent commands here, since there is no parent so no one else could have
			// set any persistent flags
			for _, p := range cw.persistentParameters {
				viper.BindPFlag(p, cmd.PersistentFlags().Lookup(p))
			}
		},
	}

	cw.cobraCommand = command
	cli.rootCmd.AddCommand(command)

	return cw
}

func (cli *cli) Run() error {
	err := cli.rootCmd.Execute()
	return internal.ReturnErrorOrPanic(err)
}

func (cli *cli) GetRootCommand() *cobra.Command {
	return cli.rootCmd
}
