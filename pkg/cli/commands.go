package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command represents any command instantiated through Cli. Keep these instances (assign them to variables), so
// that you can add additional sub-commands or parameters.
type Command interface {
	// AddCommand add a subcommand to this command. Parameters: use (command mnemonic), shortInfo (short info shown in help),
	// longDescription (long description shown in help), runFunction (function to be run when command is invoked)
	AddCommand(use string, shortInfo string, longDescription string, runFunction func()) Command

	// AddParameterString adds a string parameter to the command. Parameter value can be read using viper.GetString method.
	// Parameter "shorthand" can only be one letter string!
	AddParameterString(name string, defaultValue string, required bool, shorthand string, description string)

	// AddParameterBool adds a boolean parameter to the command. Parameter value can be read using viper.GetBool method.
	// Parameter "shorthand" can only be one letter string!
	AddParameterBool(name string, defaultValue bool, required bool, shorthand string, description string)

	// AddParameterInt adds an integer parameter to the command. Parameter value can be read using viper.GetInt method.
	// Parameter "shorthand" can only be one letter string!
	AddParameterInt(name string, defaultValue int, required bool, shorthand string, description string)

	// AddPersistentParameterString adds a string parameter to the command, which will also be available to all sub commands
	// as well. Value of the parameter can be read using viper.GetString method. Parameter "shorthand" can only be one letter string!
	AddPersistentParameterString(name string, defaultValue string, required bool, shorthand string, description string)

	// AddPersistentParameterBool adds a boolean parameter to the command, which will also be available to all sub commands
	// as well. Value of the parameter can be read using viper.GetBool method. Parameter "shorthand" can only be one letter string!
	AddPersistentParameterBool(name string, defaultValue bool, required bool, shorthand string, description string)

	// AddPersistentParameterInt adds an integer parameter to the command, which will also be available to all sub commands
	// as well. Value of the parameter can be read using viper.GetInt method. Parameter "shorthand" can only be one letter string!
	AddPersistentParameterInt(name string, defaultValue int, required bool, shorthand string, description string)

	// GetCobraCommand returns the underlying cobra.Command for this command (framework used under the hood)
	GetCobraCommand() *cobra.Command
}

type commandWrapper struct {
	cobraCommand *cobra.Command
}

func (command *commandWrapper) AddCommand(use string, shortInfo string, longDescription string, runFunction func()) Command {
	newCommand := &cobra.Command{
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

	command.cobraCommand.AddCommand(newCommand)

	return &commandWrapper{
		cobraCommand: newCommand,
	}
}

func (command *commandWrapper) AddParameterString(name string, defaultValue string, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().StringP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.Flags().Lookup(name))
}

func (command *commandWrapper) AddParameterBool(name string, defaultValue bool, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().BoolP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.Flags().Lookup(name))
}

func (command *commandWrapper) AddParameterInt(name string, defaultValue int, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().IntP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.Flags().Lookup(name))
}

func (command *commandWrapper) AddPersistentParameterString(name string, defaultValue string, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().StringP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.PersistentFlags().Lookup(name))
}

func (command *commandWrapper) AddPersistentParameterBool(name string, defaultValue bool, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().BoolP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.PersistentFlags().Lookup(name))
}

func (command *commandWrapper) AddPersistentParameterInt(name string, defaultValue int, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().IntP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	viper.BindPFlag(name, command.cobraCommand.PersistentFlags().Lookup(name))
}

func (command *commandWrapper) GetCobraCommand() *cobra.Command {
	return command.cobraCommand
}
