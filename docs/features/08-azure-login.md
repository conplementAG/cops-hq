# Azure Login recipe

This recipe is pretty small and simple, but it captures the essence of one very important thing in infra as code: both
developer and technical accounts should always be supported to run the IaC application! This is important from security
perspective, to prevent technical account identifiers and secrets to be passed around by developers. 

In the default setup, calling Login() will log the user in (if not already logged in) via normal prompts via azure CLI. 
If specific viper variables are set (check the code docs on `azure_login.New` or `NewWithParams` methods), the following login methods are also supported (useful for CI systems).

## Managed identity
You can provide a flag whether to use azure managed identities for the login. 
### User assigned managed identity
Login via a user assigned managed identity can be done by additionally providing the client id.
### System assigned managed identity
The system assigned managed identity is used when the client id is ommitted

## Service Principal
By providing the client-id, client-secret, tenant-id you can login via a service principal as well. You also have to ommit the flag to use a managed identity.

## Usage

```go
login := azure_login.New(hq.GetExecutor())
login.Login()
```
The login mechanisms which will be attempted in the following order:
- User assigned managed identity
- System assigned managed identity
- Service Principal
- Normal user login