# Command execution with the Executor

Common advantage of Shell scripting over other languages like Go or Python for scripting, is that commands can be simply written in plaintext,
without any hassle. Well, cops-hq `commands.Executor` enables you to do the same, for example:

``` go
result, err := executor.Execute("echo test")
```

With `commands.Executor` you get full native output and error parsing, in statically typed objects in Go. The result, which is always a
string, can simply be converted to a struct or any other object if parsable (like commands outputting json or yaml). Error checking is 
also much simpler then in shell scripting, because both command errors and stderr outputs are considered as errors. 

Make sure to check out the code docs on the `commands.Executor`interface for available methods and settings!

## Instantiation

Usually, you will use the [HQ](03-hq.md) setup, in which you can call the `hq.GetExecutor()` method to get the currently configured 
executor. Executor can be configured in one of two modes:
- chatty, which is similar to shell scripting - all commands and their outputs are shown on stdout / stderr
- quiet, in which no command output will be piped to stdout/stderr (but it is always written to log file), unless a verbose flag
  is provided via viper (can be loaded from CLI arguments, more info about that in next sections of the documentation). Quiet
  mode is interesting because it removes all the clutter of command output from the console, and you then only show semantic messages
  like `logrus.Info("deploying something important...)`. Pattern with configuring the executor in quiet mode is keeping it quiet
  for local developers, but setting the `--verbose` in CI systems so that the full output is always shown. 

You can instantiate the Executor yourself using one of the factory methods as well. Multiple instance of Executors are supported in 
your applications, in case you want multiple executors with different setups. Usually however, just ask your hq.HQ instance for one.
