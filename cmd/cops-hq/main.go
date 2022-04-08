package main

import (
	"fmt"
	"github.com/denisbiondic/cops-hq/pkg/commands"
	"github.com/denisbiondic/cops-hq/pkg/logging"
	"github.com/sirupsen/logrus"
)

// main in the cops-hq project is a simple "consumer" CLI, which can be used for manual code tests.
// Logic and code here should be kept to minimum. Use automated tests instead.
func main() {
	logging.Init("hq.log")
	logrus.Info("Running the test CLI...")

	exec := commands.NewExecutor("hq.log")

	logrus.Info("Testing the docker correctly outputs in TTY...")
	exec.ExecuteTTY("docker build .")

	logrus.Info("Testing the output can be parsed...")
	out, _ := exec.Execute("echo test")
	fmt.Println(out)

	logrus.Info("Testing the stderr is also shown...")
	exec.Execute("ls this-file-does-not-exist")
}
