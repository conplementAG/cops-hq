package main

import (
	"bufio"
	"fmt"
	"github.com/denisbiondic/cops-hq/pkg/cli"
	"github.com/denisbiondic/cops-hq/pkg/commands"
	"github.com/denisbiondic/cops-hq/pkg/hq"
	"github.com/denisbiondic/cops-hq/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// main in the cops-hq project is a simple "consumer" CLI, which is used as a showcase.
// Logic and code here should be kept to minimum. For testing, use automated tests instead of extending the showcase here.
func main() {
	fmt.Print("Run hq or step by step setup? Type hq or step to proceed: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()

	if text == "hq" {
		simpleHqSetup()
	} else if text == "step" {
		stepByStepSetup()
	}
}

func simpleHqSetup() {
	hq := hq.New("hq", "0.0.1", "hq.log")
	hq.Cli.AddBaseCommand("infrastructure", "Infrastructure command", "Infrastructure command", func() {
		logrus.Info("infra command running...")
	})

	hq.Run()

	hq.Executor.ExecuteWithProgressInfo("echo hello")
}

func stepByStepSetup() {
	// Logging features -> initializing logging system will log to both console and the configured file
	logger := logging.Init("hq.log")
	logrus.Info("Running the test CLI...")

	// Executor features -> this is a wrapper which can be used for running any command, incl. output parsing, error handling etc.
	exec := commands.NewChatty("hq.log", logger)

	// for example, TTY works as well
	logrus.Info("Testing the docker correctly outputs in TTY...")
	exec.ExecuteTTY("docker build .")

	logrus.Info("Testing the output can be parsed...")
	out, _ := exec.Execute("echo test")
	fmt.Println(out)

	logrus.Info("Testing the stderr is also shown...")
	exec.Execute("ls this-file-does-not-exist")

	// CLI features -> cops-hq offers a simple CLI builder, based on Cobra
	// to test the code below, call 'infrastructure create -e abc -a def'
	cli := cli.New("hq", "0.0.1")
	infraCommand := cli.AddBaseCommand("infrastructure", "command group for infrastructure commands", "command group for infrastructure commands", nil)
	infraCommand.AddPersistentParameterString("environment-tag", "", true, "e", "")

	createCommand := infraCommand.AddCommand("create", "create the infrastructure", "create the infrastructure", func() {
		logrus.Info("Creating infrastructure...")

		// Configuration management features -> CLIs are fully integrated with Viper library, so it can directly be used to retrieve parameters
		logrus.Info("Got parameter: " + viper.GetString("environment-tag"))
		logrus.Info("Got parameter: " + viper.GetString("account"))
	})

	createCommand.AddParameterString("account", "", true, "a", "sample parameter")

	cli.Run()
}
