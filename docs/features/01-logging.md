# Logging

Logging in cops-hq is highly-opinionated, in which [logrus](https://github.com/sirupsen/logrus) library is expected to be directly used. Logging in cops-hq 
works by piping every log message to both stdout and configured file. Some common problems, like log rotation, or color
output on Windows, are solved out of the box. 

After the logging is initialized, logging is simply a matter of calling of the logrus methods, like `logrus.Info("message")`.

However, you usually don't initialize the Logging system directly via `logging.Init()` method, but you use the [HQ](03-hq.md) 
setup which has this out of the box. However, logging.Init() might be used directly in cases where you only use the
[Command Execution](02-command-execution.md) part of cops-hq. 

Note: logging should be initialized only once per application, since it uses a global singleton pattern. 