package main

import (
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/sirupsen/logrus"
)

// main in the example-infra project is a simple "consumer" CLI, which is used as a showcase.
// Logic and code here should be kept to minimum. For testing, use automated tests instead of extending the showcase here.
func main() {
	// We create a HQ instance in normal "chatty" mode, in which all commands and outputs are written on the console.
	// We could also use the hq.NewQuiet to override this behaviour, to keep the console cleaner
	hq := hq.New("example-infra", "0.0.1", "example-infra.log")

	hq.GetCli().AddBaseCommand("infrastructure", "Infrastructure command", "Example infrastructure command", func() {
		createInfrastructure(hq)
	})

	// this will start the parsing the os.Args given to the application, and execute the matching CLI command
	hq.Run()
}

func createInfrastructure(hq hq.HQ) {
	logrus.Info("infra command running...")

	// to showcase the correct output piping to stdout, we execute the most complex TTY command we know :)
	err := hq.GetExecutor().ExecuteTTY("docker build . -f example.Dockerfile")
	panicOnError(err)

	// ... or we can just use normal commands
	result, err := hq.GetExecutor().Execute("echo test")
	panicOnError(err)

	if result != "test" {
		panic("Expected to get 'test'")
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
