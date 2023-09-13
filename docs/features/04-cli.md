# CLI

As we believe that self-documenting CLIs are the way to go when writing infrastructure-as-code, the support for CLI commands is 
provided out of the box. Libraries we use are Cobra and Viper, in which any parameter declared on Cobra will also be available 
via Viper for querying in your code. 

Viper and Cobra usage is wrapped in cops-hq, but the underlying Cobra object can always be accessed for any command. In most cases, 
this should however not be necessary. 

## Setup

In most cases, you should use the CLI instance provided via your HQ instance: `hq.GetCLI()`. You could however also use the `cli.New()`
factory method directly. 

Setting up your commands is as simple as:

```go 
hq.GetCli().AddBaseCommand("infrastructure", "Infrastructure command", "Example infrastructure command", func() {
    // do something
})
```

Every `AddXXX()` method returns a Command instance, which can be used to add additional child commands or parameters. 

Per default, CLI will set up two parameters which will be available on each and every command in the CLI. These are the `verbose` 
and `silence-long-running-progress-indicators` flags, which executor uses under the hood.

You can add additional parameters for any command you want, by calling either `AddParameterXXX` or `AddPersistentParameterXXX`. The
latter will add the parameter not only on this command, but any child command as well. 

## In-build commands

The CLI automatically adds a command group called "hq" with subcommands documented below.

### Dependency check

Check the 03-hq.md section on dependency checking.