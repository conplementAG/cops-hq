# cops-hq

Base infrastructure as code libraries for projects using Golang instead of Bash and other languages.

## Concept

This is an opinionated library, which sets the following goals and hard dependencies:
- scripted command's output should be shown on stdout, but also recorded to a log file
- logging is done via logrus
- configuration management is done via Viper, which is also often used as the injection mechanism for required recipe 
parameters. Check the recipe code comment documentation for reference. 
- 

## Architecture and module dependencies

Library is envisioned as a collection of packages that can be independently used. There are however some

## Contribution

Check [contribution](docs/contribution.md).  

