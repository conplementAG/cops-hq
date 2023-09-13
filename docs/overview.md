# Features

Cops-hq offers many features, and most can be used in isolation. For example, if you only wanted to write pure IaC scripts in Go, 
you can simply use the [Command Executor](features/02-command-execution.md) library. If you want a higher-level support with common tasks
already solved out of the box, you can use one of the recipes like [Terraform Recipe](features/09-terraform.md).

Note: all public exposed members are documented via GoDoc, so make sure you check the source code as well! Tests are also a nice
way to understand the features. 

# Contents

- [Overview](overview.md)
- [Logging](features/01-logging.md)
- [Command Execution](features/02-command-execution.md)
- [HQ](features/03-hq.md)
- [CLI](features/04-cli.md)
- [Configuration Management](features/05-configuration.md)
- [Naming Service](features/06-naming.md)
- [Error handling](features/07-error-handling.md)
- [Recipe - Azure Login](features/08-azure-login.md)
- [Recipe - Terraform](features/09-terraform.md)
- [Recipe - Helm](features/10-helm.md)
- [Recipe - Sops](features/11-sops.md)
- [Dockerfile for CI/CD](features/99-dockerfile.md)