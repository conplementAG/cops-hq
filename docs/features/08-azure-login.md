# Azure Login recipe

This recipe is pretty small and simple, but it captures the essence of one very important thing in infra as code: both
developer and technical accounts should always be supported to run the IaC application! This is important from security
perspective, to prevent technical account identifiers and secrets to be passed around by developers. 

In the default setup, calling Login() will log the user in (if not already logged in) via normal prompts via azure CLI. 
If viper variables for service principal info are set (check the code docs on `azure_login.New` or `NewWithParams` methods), 
then the service principal login will be used (useful for CI systems).

# usage

```go
login := azure_login.New(hq.GetExecutor())
login.Login()
```