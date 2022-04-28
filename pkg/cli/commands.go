package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Command struct {
	cobraCommand *cobra.Command
}

func (cli *Cli) AddRootCommand(use string, short string, long string, runFunction func()) *Command {
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
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

func (command *Command) AddCommand(use string, short string, long string, runFunction func()) *Command {
	newCommand := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			runFunction()
		},
	}

	command.cobraCommand.AddCommand(newCommand)

	return &Command{
		cobraCommand: newCommand,
	}
}

func (command *Command) AddParameterString(name string, defaultValue string, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().StringP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.Flags().Lookup(name))
}

func (command *Command) AddParameterBool(name string, defaultValue bool, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().BoolP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.Flags().Lookup(name))
}

func (command *Command) AddParameterInt(name string, defaultValue int, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().IntP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.Flags().Lookup(name))
}

func (command *Command) AddPersistentParameterString(name string, defaultValue string, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().StringP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.PersistentFlags().Lookup(name))
}

func (command *Command) AddPersistentParameterBool(name string, defaultValue bool, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().BoolP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.PersistentFlags().Lookup(name))
}

func (command *Command) AddPersistentParameterInt(name string, defaultValue int, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().IntP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.PersistentFlags().Lookup(name))
}
