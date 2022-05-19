package terraform

import "github.com/conplementag/cops-hq/pkg/commands"

// New creates a new instance of Terraform, which is a wrapper around common Terraform functionality. In includes remote state
// backend configuration and other best practices when dealing with terraform.
// Parameters:
//     executor (can be provided from hq.GetExecutor() or by instantiating your own),
//     projectName (a pseudo-name for your terraform deployment, useful in case you have multiple terraform projects
//         in one IaC. In doubt, set this to "app" if you only have a single terraform project),
//     subscriptionId and tenantId (required for terraform state setup)
//     region (Azure region, e.g. westeurope, required for terraform state setup)
//     resourceGroupName (name of the resource group where the terraform state will be stored. The resource group will
//         also be created per default, unless overridden via backendStorageSettings. It is recommended to
//         use the naming.Service to generate this name)
//     stateStorageAccountName (name of the storage account where the terraform state will be stored. It is recommended to
//         use the naming.Service to generate this name)
//     terraformDirectory (directory where your terraform resources are stored. To construct the full path, simply use
//         filepath.Join() method and the hq.ProjectBasePath)
//     backendStorageSettings (various settings which can be read or set for terraform backend setup, but it is best not
//         to override these. Simply set to terraform.DefaultBackendStorageSettings)
//     deploymentSettings (various settings which can be read or set for terraform deployments, but it is best not
//         to override these. Simply set to terraform.DefaultDeploymentSettings
func New(executor commands.Executor, projectName string,
	subscriptionId string, tenantId string, region string,
	resourceGroupName string, stateStorageAccountName string, terraformDirectory string,
	backendStorageSettings BackendStorageSettings, deploymentSettings DeploymentSettings) Terraform {

	return &terraformWrapper{
		executor:                executor,
		projectName:             projectName,
		subscriptionId:          subscriptionId,
		tenantId:                tenantId,
		region:                  region,
		resourceGroupName:       resourceGroupName,
		stateStorageAccountName: stateStorageAccountName,
		terraformDirectory:      terraformDirectory,
		storageSettings:         backendStorageSettings,
		deploymentSettings:      deploymentSettings,
	}
}
