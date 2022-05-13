package cli

import (
	"fmt"
	"github.com/conplementag/cops-hq/internal/testing_utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SubcommandsWithParametersAndViper(t *testing.T) {
	// Arrange
	correctCommandActionExecuted := false
	cli := New("myprog", "0.0.1")

	// Act
	command := cli.AddBaseCommand("test", "Simple test command", "big description", func() {
		fmt.Println("some action")
	})

	command.AddParameterString("login-user", "", true, "u", "first test arg")
	command.AddParameterString("login-pass", "", true, "p", "second test arg")
	command.AddParameterBool("admin", false, false, "a", "third bool flag test")
	command.AddParameterInt("retries", 1, false, "r", "integer flag test")

	subCommand := command.AddCommand("me", "a test subcommand", "will do a lot of stuff", func() {
		correctCommandActionExecuted = true
	})

	subCommand.AddParameterString("argX", "Y", true, "a", "first arg")
	subCommand.AddParameterString("argY", "X", false, "b", "second arg")

	testing_utils.PrepareCommandForTesting(cli.GetRootCommand(), "test", "me", "--argX", "W")
	cli.Run()

	// Assert
	assert.True(t, correctCommandActionExecuted)
	assert.Equal(t, "W", viper.GetString("argX")) // should have the new value from args
	assert.Equal(t, "X", viper.GetString("argY")) // should keep the default value
}

func Test_CommandShouldShowHelpWhenNoRunFunctionGiven(t *testing.T) {
	// Arrange
	cli := New("myprog", "0.0.1")
	cli.AddBaseCommand("test", "Simple test command", "big description", nil)
	outputBuffer := testing_utils.PrepareCommandForTesting(cli.GetRootCommand(), "test")

	// Act
	cli.Run()

	// Assert
	assert.Contains(t, testing_utils.ReadBuffer(t, outputBuffer), "big description")
}

func Test_RequiredParametersShouldPreventCommandExecutionWhenNotProvided(t *testing.T) {
	// Arrange
	commandActionCalled := false

	cli := New("myprog", "0.0.1")
	testCommand := cli.AddBaseCommand("test", "Simple test command", "big description", func() {
		commandActionCalled = true
	})
	testCommand.AddParameterString("param", "", true, "p", "test arg")
	outputBuffer := testing_utils.PrepareCommandForTesting(cli.GetRootCommand(), "test")

	// Act
	cli.Run()

	// Assert
	assert.False(t, commandActionCalled)
	assert.Contains(t, testing_utils.ReadBuffer(t, outputBuffer), "required flag")
}

func Test_PersistentParametersAreSharedWithSubcommands(t *testing.T) {
	// Arrange
	cli := New("myprog", "0.0.1")

	testCommand := cli.AddBaseCommand("infra", "infra command group", "", nil)
	testCommand.AddPersistentParameterString("environment-tag", "", true, "e", "env tag")

	testCommand.AddCommand("create", "create infra", "", func() {})

	outputBuffer := testing_utils.PrepareCommandForTesting(cli.GetRootCommand(), "infra", "create")

	// Act
	cli.Run()

	// Assert
	output := testing_utils.ReadBuffer(t, outputBuffer)
	assert.Contains(t, output, "required flag")
	assert.Contains(t, output, "environment-tag")
}

func Test_PersistentParametersAreAvailableThroughViperInSubcommands(t *testing.T) {
	// Arrange
	expectedActionCalled := false
	cli := New("myprog", "0.0.1")

	testCommand := cli.AddBaseCommand("infra", "infra command group", "", nil)
	testCommand.AddPersistentParameterString("environment-tag", "", true, "e", "env tag")

	testCommand.AddCommand("create", "create infra", "", func() {
		expectedActionCalled = true
		assert.Equal(t, "prod", viper.GetString("environment-tag"))
	})

	testing_utils.PrepareCommandForTesting(cli.GetRootCommand(), "infra", "create", "-e", "prod")

	// Act
	cli.Run()

	// Assert
	assert.True(t, expectedActionCalled)
}
