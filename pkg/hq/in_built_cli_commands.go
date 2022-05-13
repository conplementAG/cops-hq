package hq

import "github.com/conplementag/cops-hq/pkg/cli"

func addInbuiltHqCliCommands(cli cli.Cli, container *hqContainer) {
	hqBaseCommand := cli.AddBaseCommand("hq", "in-build HQ command group", "Command predefined by cops-hq.", nil)

	hqBaseCommand.AddCommand("check-setup-and-dependencies", "Checks the installed tools and the project structure",
		"Use this command to check the installed versions of tools such as azure-cli or kubectl, and to check the project "+
			"structure for expected directories and files.", func() {
			container.CheckSetupAndDependencies()
		})
}