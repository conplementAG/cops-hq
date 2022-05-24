# cops-hq

Base infrastructure as code libraries for projects using Golang instead of Bash, Python and other languages. Highly opinionated, 
but useful if you are using Azure and Kubernetes. 

## Concepts

This is an opinionated library, which sets the following goals and hard dependencies:
- This library is a base library for orchestration of tools commonly used in the DevOps daily business (terraform, azure cli, etc.)
- These tools are expected to be pre-installed (there is a method in code and a CLI command to test this).
- Entry point for any IaC application should be the CLI, with self-explanatory commands. For this purpose, Cobra is used, but wrapped
  to simplify the setup and solve the common problems. CLI should represent all the tasks (runbooks) the DevOps team has in daily operations
  (e.g. infrastructure create, application deployment, configuring access to a certain resource etc.)
- scripted command's output should be shown on stdout, but also recorded to a log file. This behavior can also be customized. 
- logging is done via logrus, and logrus only
- configuration management is done via Viper and Mozilla Sops. Configuration is checked into the code in an encrypted form
  (usually encrypted via Azure KeyVault) 
- naming of Azure resources should be done in a standardized (although customizable) way.  

## Usage

For basic usage, check out the `main.go` file in `cmd/cops-hq`.

Explanation of all concepts and features can be found in [Features Overview](docs/features/00-overview.md)

## Contribution

Check [contribution](docs/contribution.md).