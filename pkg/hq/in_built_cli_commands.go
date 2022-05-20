package hq

import (
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/sirupsen/logrus"
)

func addInbuiltHqCliCommands(cli cli.Cli, container *hqContainer) {
	hqBaseCommand := cli.AddBaseCommand("hq", "in-build HQ command group", "Command predefined by cops-hq.", nil)

	hqBaseCommand.AddCommand("check-dependencies", "Checks the installed tools and the project structure",
		"Use this command to check the installed versions of tools such as azure-cli or kubectl, and to check the project "+
			"structure for expected directories and files.", func() {
			err := container.CheckToolingDependencies()

			if err != nil {
				logrus.Error(err)
				panic(err)
			}
		})
}
