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
	cobraCommand         *cobra.Command
	parameters           []string
	persistentParameters []string
	parentCommand        *commandWrapper
}

func (command *commandWrapper) AddCommand(use string, shortInfo string, longDescription string, runFunction func()) Command {
	cw := &commandWrapper{
		parentCommand: command,
	}

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
		PreRun: func(cmd *cobra.Command, args []string) {
			// we have to map the viper parameters on runtime, when the command is executing, to prevent
			// overwriting of viper mappings in case multiple commands have the same named parameters
			for _, p := range cw.parameters {
				viper.BindPFlag(p, cmd.Flags().Lookup(p))
			}

			// as for the persistent parameters, only one PreRun is running at the time (of the executing command),
			// so we have to traverse and find all persistent parameters of the parent commands as well
			traverseAndBindFlagsToViper(cw)
		},
	}

	cw.cobraCommand = newCommand
	command.cobraCommand.AddCommand(newCommand)

	return cw
}

func traverseAndBindFlagsToViper(command *commandWrapper) {
	if command == nil {
		return
	}

	for _, p := range command.persistentParameters {
		viper.BindPFlag(p, command.GetCobraCommand().PersistentFlags().Lookup(p))
	}

	traverseAndBindFlagsToViper(command.parentCommand)
}

func (command *commandWrapper) AddParameterString(name string, defaultValue string, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().StringP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	command.parameters = append(command.parameters, name)
}

func (command *commandWrapper) AddParameterBool(name string, defaultValue bool, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().BoolP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	command.parameters = append(command.parameters, name)
}

func (command *commandWrapper) AddParameterInt(name string, defaultValue int, required bool, shorthand string, description string) {
	command.cobraCommand.Flags().IntP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkFlagRequired(name)
	}

	command.parameters = append(command.parameters, name)
}

func (command *commandWrapper) AddPersistentParameterString(name string, defaultValue string, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().StringP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	command.persistentParameters = append(command.persistentParameters, name)
}

func (command *commandWrapper) AddPersistentParameterBool(name string, defaultValue bool, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().BoolP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	command.persistentParameters = append(command.persistentParameters, name)
}

func (command *commandWrapper) AddPersistentParameterInt(name string, defaultValue int, required bool, shorthand string, description string) {
	command.cobraCommand.PersistentFlags().IntP(name, shorthand, defaultValue, description)

	if required {
		command.cobraCommand.MarkPersistentFlagRequired(name)
	}

	command.persistentParameters = append(command.persistentParameters, name)
}

func (command *commandWrapper) GetCobraCommand() *cobra.Command {
	return command.cobraCommand
}
