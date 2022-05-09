package cli

import "github.com/spf13/cobra"

type cli struct {
	programName string
	version     string
	rootCmd     *cobra.Command
}

func (cli *cli) AddBaseCommand(use string, shortInfo string, longDescription string, runFunction func()) *Command {
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
	}

	cli.rootCmd.AddCommand(command)

	return &Command{
		cobraCommand: command,
	}
}

func (cli *cli) Run() error {
	return cli.rootCmd.Execute()
}

func (cli *cli) GetRootCommand() *cobra.Command {
	return cli.rootCmd
}
