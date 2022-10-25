package cli

import (
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type cli struct {
	programName    string
	version        string
	rootCmd        *cobra.Command
	defaultCommand string
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
	if cli.defaultCommand != "" {
		rootCommand := cli.GetRootCommand()
		registeredCommands := rootCommand.Commands()

		var isCommandSet = false
		for _, a := range registeredCommands {
			for _, b := range os.Args[1:] {
				if a.Name() == b {
					isCommandSet = true
					break
				}
			}
		}

		// if no command set on the command line, use the default command by extending the existing command line args
		if !isCommandSet {
			args := append([]string{cli.defaultCommand}, os.Args[1:]...)
			rootCommand.SetArgs(args)
		}
	}

	err := cli.rootCmd.Execute()
	return internal.ReturnErrorOrPanic(err)
}

func (cli *cli) GetRootCommand() *cobra.Command {
	return cli.rootCmd
}

func (cli *cli) OnInitialize(initFunction func()) {
	cobra.OnInitialize(initFunction)
}

func (cli *cli) SetDefaultCommand(command string) {
	cli.defaultCommand = command
}
