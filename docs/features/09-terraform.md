# Terraform recipe

Dealing with terraform in IaC code is not difficult, but there are common tasks that should (almost) always be done:
- set up the state management in a remote storage location (e.g. Azure Storage account)
- initializing terraform before any other command
- clear the local cache before re-initializing terraform
- common flows like plan, save plan to file, ask user to confirm, apply the saved file
- CI setups

To have these topics solved out of the box, use the `terraform.New` to create a new wrapper around this functionality. 
Make sure to read the code docs for many parameters which are required. You can also create multiple instances of terraform
recipe, which might be useful if you have multiple separate terraform projects in your IaC code (e.g. core terraform, for 
things common to all environments, and app terraform which contains env specific resources).

Terraform recipe supports two ways to operations:
- you can call the methods which more-or-less map 1:1 to terraform functionality, like `PlanDeploy` and `ForceDeploy`, or
- you can call a higher-order method like `DeployFlow` which contains best practice around deploying (e.g. plan, ask user
  for confirmation, apply changes etc.) `DeployFlow` is an example of a method written with both local development and CI 
  in mind, because it provides a single entry point configurable for all scenarios.

# basic usage

```go 

tf := terraform.New(o.executor, "your-tf-project-name",
    // usually, these three are loaded from your sops based config file
    viper.GetString(config_file_keys.KeySubscriptionId),
    viper.GetString(config_file_keys.KeyTenantId),
    viper.GetString(config_file_keys.KeyRegion),
    
    // resource group and the storage account where this terraform state will be stored
    naming.GetCoreResourceGroup(), 
    storageAccountName,
    
    // path to your terraform files
    filepath.Join(copshq.ProjectBasePath, "terraform", "core"),
    
    // settings which can be overriden
    terraform.DefaultBackendStorageSettings,
    terraform.DefaultDeploymentSettings)

// first step is always to init the terraform
err = tf.Init()
... do something with the error

// then, we need to set all the terraform variables. This is always a simple map of strings, but supporting any
// simple or complex object (check the tf.SetVariables code docs)
vars := make(map[string]interface{})

vars["var_x"] = "some value"

err = tf.SetVariables(vars)
... do something with the error

// then, we can start the real deployment flow. Via CLI flags plan-only and use-existing-plan (you need to define them yourself), 
// we can control the set up from the outside, for example, a local developer would have these flags set to false, so plan would 
// be created before deploying. In CI system, we would have separate jobs for creating and apply the plan, as a job
// to wait for user confirmation would have to be set in between. 
err = tf.DeployFlow(viper.GetBool(cli_flags.PlanOnly), viper.GetBool(cli_flags.UseExistingPlan))
... do something with the error

```
