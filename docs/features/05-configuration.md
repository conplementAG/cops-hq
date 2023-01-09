# Configuration

As an interface for all configuration, irrelevant of the source, viper library should be used, for example:

``` go
value := viper.GetString("variable_name")
```

Viper is per default loaded with these sources:
- all environment variables (overrides other same named keys, priority source)
- parameters defined for the CLI

## Application configuration

There are many ways on providing application configuration parameters to your IaC code. For example, you can provide all 
parameters via environment variables, or via CLI parameters. However, this tends to lead to so-called "configuration sprawl",
in which parameters are repeated over and over again (first in you CI code, then as env variables, then as CLI parameters etc.)

A better way of doing this is having (almost) all the parameters checked in your source code repository, with secrets in encrypted 
form. For this purpose, we recommend the Mozilla Sops project, which is pretty each to use and set up. 

After you install sops on your machine, create a config directory in the root of your project. In it, create a file called .sops.yaml

```yaml 
# azure_keyvault is configured with an ID of the key which should be used to decrypt specified config files.
# Decryption will work, if the currently logged-in user has access to the specified key. All team developers should
# per default have access to the non-prod key, while prod key access should only be granted to the CI/CD service principal.
creation_rules:
  - path_regex: \.prod.yaml$
    azure_keyvault: prod_keyvault_key_id_from_azure_portal
    encrypted_suffix: _secret

  - path_regex: ""
    azure_keyvault: dev_keyvault_key_id_from_azure_portal
    encrypted_suffix: _secret
```

What this configuration does, is basically it says that sops should use the "prod_keyvault_key_id_from_azure_portal" key
to decrypt / encrypt production env configs, and  "dev_keyvault_key_id_from_azure_portal" for all other environments. 
Current logged-in user should of course have access to these key(s), which can be granted in Azure portal.   

Reading / writing of configuration files is always done via sops CLI, which uses the currently configured system text editor. 
If you want to override this, you could run
#### Linux and Mac
```shell
export EDITOR="code -w"
```
or
#### Windows (as Admin)
```shell
setx EDITOR "code -w" /m
```
to set VS Code as text editor before executing any sops command.

First, create a template file which developers can re-use in the future for own personal dev environments. You can store all 
the default configuration values here, and check this file in your source control. 

```shell
sops local-template.yaml
```

To create a new configuration file for your environment, follow the structure <<env-name>.yaml, e.g. you can create a config 
for production by typing `sops prod.yaml`. While having both editors open (for both local-template.yaml and prod.yaml), you 
can copy-paste the values. Closing the editor will automatically encrypt and save the file(s).

To load the contents for these configuration files into Viper, we provide a handy method `hq.LoadEnvironmentConfigFile()`. This
method expect a viper variable called "environment-tag" to be defined and set. 

To showcase the complete "solution", usual pattern is something like this:

```go
// since we use sops with azure key vault, we need to log into the Azure first, so that sops can crypt / decrypt the config files
login := azure_login.New(hq.GetExecutor())
login.Login()

// now we can load the env file. Usually, the environment-tag variable will be set through CLI at 
// runtime (e.g. infra create --environment-tag prod)
hq.LoadEnvironmentConfigFile()

// if we had the subscription ID stored in the env specific config file, we can set it now as the current azure subscription
login.SetSubscription(viper.GetString(config_file_keys.KeySubscriptionId))
```

