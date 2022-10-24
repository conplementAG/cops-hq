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

## basic usage

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
// simple or complex object (check the tf.SetVariables code docs). Don't serialize your variables to JSON / string, 
// this is not required and it will not work!
vars := make(map[string]interface{})

vars["var_x"] = "some value"

// alternatively you can auto-populate your map of strings from your config (e.g. sops) with the help of viper,
// but keep in mind that terraform only supports flat key structures out of the box.
vars = viper.GetStringMap("infrastructure")

err = tf.SetVariables(vars)
... do something with the error

// then, we can start the real deployment flow. Via CLI flags plan-only, use-existing-plan and auto-approve (you need to define them yourself), 
// we can control the set up from the outside, for example, a local developer would have these flags set to false, so plan would 
// be created before deploying. In CI system, we would have separate jobs for creating and apply the plan, as a job
// to wait for user confirmation would have to be set in between. 
err = tf.DeployFlow(viper.GetBool(cli_flags.PlanOnly), viper.GetBool(cli_flags.UseExistingPlan), viper.GetBool(cli_flags.AutoApprove))
... do something with the error

```

## Terraform plan and detecting changes in CI/CD

Terraform recipe automatically persists the plans in the .plans directory, in the same place where you specified that your terraform 
sources are located. Additionally, to the plan file in terraform format, text and json representations of the same plan file are persisted
as well. These can easily be used in CI/CD for advanced use cases like automatic approval on no terraform changes. For this purpose, 
take a look at the plan_analyzer object and its IsDeployPlanDirty / IsDestroyPlanDirty methods, which following the example above, could be used as:

```go 

analyzer := plan_analyzer.New(yourProjectName, terraformSourcesDirectory)
result, err := analyzer.IsDeployPlanDirty()

fmt.Println(result)

```

For CI/CD usage, you can simply wrap required method into your own CLI method (e.g. `infrastructure is-deploy-plan-dirty` or similar), 
which calls this Go method under the hood. The true / false output can then be parsed with a simple bash command or similar. 

## Additional Terraform plan formats

The plans saved in .plans directory follow the naming convention: `<<project_name>.(destroy|deploy).tfplan`
Plaintext (command line output of terraform plan) and JSON versions (extra `terraform show -json` call output) of the same plan
file are saved in the same place, with .txt and .json extensions respectively.

## Variable auto-populate from config file

In the example above, a code snippet is included on how to auto-populate the variables from a config section. You could also combine
both auto-population and your own keys, since this is just simple map manipulation in Go!