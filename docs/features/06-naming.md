# Naming Service

Naming is often a topic when dealing with Azure resources. Many resources have different naming limitations, so having the 
naming topic solved out of the box is a good thing. 

`naming.Service` can be instantiated via its factory method `naming.New()`. Check the code docs on the explanation of the
parameters. Multiple instance of the naming service are supported. 

## usage

Basic usage will look something like this:

```go
	environmentTag := viper.GetString(cli_flags.KeyEnvironmentTag)
	region := viper.GetString(config_file_keys.KeyRegion)

	namingService, err := naming.New("myproject", region, environmentTag, "myservice")
    ... do something about the error
	
	resourceGroupName, err := namingService.GenerateResourceName(resources.ResourceGroup, "")
    ... do something about the error
```

You can check the tests for ideas about more advanced usage.

It is not a bad idea to wrap the naming service setup and calls in your own object, and expose only methods for your 
application specific resources (e.g. GetCoreResourceGroup(), GetBackupStorageAccountName() etc.)

## patterns

The naming service comes with a default naming pattern, but you can change this to support any naming schema you want. 