package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conplementag/cops-hq/internal"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Terraform is a wrapper around common terraform functionality used in IaC projects with Azure. In includes remote state
// backend configuration and other best practices when dealing with terraform.
type Terraform interface {
	// Init initializes the Terraform project, also creating the backend remote storage, and the Azure resource group
	// (if not overridden via BackendStorageSettings)
	Init() error

	// SetVariables is required, before any of the following methods like PlanDeploy or Deploy are called. Variables set
	// will be applied on any subsequent operation. Parameters are:
	//     terraformVariables (this is a map of terraform variables, defined as string keys, and serialized via json serializer,
	//         which means any complex or simple type is supported. Simply set the value to whatever type you want, as long as
	//         it is properly serializable. For example, make sure complex types have the required json or mapstruct annotations,
	//         and keep in mind that only the public struct members will be serialized!)
	SetVariables(terraformVariables map[string]interface{}) error

	// DeployFlow is the method which most deployments should use. It provides a single method, which can be used for all
	// IaC cases, like local deployment, or deployment in CI systems.
	// Parameters support are
	//     planOnly (will only create the plan)
	//     useExistingPlan (will reuse existing plan from the disk)
	// Here is a short explanation:
	//     DeployFlow(false, false) will show the plan, prompt the user, and apply if confirmed
	//     DeployFlow(true, false) will only show the plan, but also persist it on disk (use GetDeployPlanFileName() for details)
	//     DeployFlow(false, true) will reuse the plan already saved on disk, and apply it without any user confirmations
	//     DeployFlow(true, true) will just show the plan persisted on the disk, without generating a new plan.
	// It is best practice to set the both planOnly and useExistingPlan from the CLI, so that CI scripts can simply override
	// the variables depending on the current CI step (usually a plan is presented, user is awaited for approval, then the existing
	// plan is applied).
	DeployFlow(planOnly bool, useExistingPlan bool) error

	// DestroyFlow is same as DeployFlow, but only for destroy.
	DestroyFlow(planOnly bool, useExistingPlan bool) error

	// PlanDeploy executes the terraform plan for deployment, returning the changes as a string. Plan output is always
	// saved to a file as well. Common pattern is to show the changes to the user, ask for confirmation, and then to Deploy the plan.
	// For this purpose, you could use the DeployFlow method.
	PlanDeploy() (string, error)

	// PlanDestroy is same as PlanDeploy, but only for destroy. Consider using DestroyFlow method as well.
	PlanDestroy() (string, error)

	// ForceDeploy deploys the plan persisted on disk via PlanDeploy. User will not be asked for any confirmations, so it is
	// your job in code to present the plan, and prompt for confirmation! For this purpose, you can use the DeployFlow method.
	ForceDeploy() error

	// ForceDestroy is same as ForceDeploy, but only for destroy. Consider using DestroyFlow method as well.
	ForceDestroy() error

	// GetBackendStorageSettings returns the backend remote state storage settings, which can be read or modified if desired
	GetBackendStorageSettings() *BackendStorageSettings

	// GetDeploymentSettings returns the current deployment settings, which can be read or modified if desired
	GetDeploymentSettings() *DeploymentSettings

	// GetDeployPlanFileName returns the file name in which the terraform deploy plan will be stored. This name is convention based
	// on the currently set project parameter while creating the terraform wrapper instance
	GetDeployPlanFileName() string

	// GetDestroyPlanFileName returns the file name in which the terraform destroy plan will be stored. This name is convention based
	// on the currently set project parameter while creating the terraform wrapper instance
	GetDestroyPlanFileName() string

	// GetVariablesFileName returns the file name in which the terraform variables will be stored. This name is convention based
	// on the currently set project parameter while creating the terraform wrapper instance
	GetVariablesFileName() string
}

type terraformWrapper struct {
	executor commands.Executor

	projectName             string
	subscriptionId          string
	tenantId                string
	region                  string
	resourceGroupName       string
	stateStorageAccountName string
	terraformDirectory      string

	variablesSet bool
	variables    map[string]interface{}

	storageSettings    BackendStorageSettings
	deploymentSettings DeploymentSettings
}

func (tf *terraformWrapper) Init() error {
	if tf.storageSettings.CreateResourceGroup {
		logrus.Info("Deploying the project " + tf.projectName + " resource group " + tf.resourceGroupName + "...")
		_, err := tf.executor.Execute("az group create -l " + tf.region + " -n " + tf.resourceGroupName)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	}

	logrus.Info("Deploying the " + tf.projectName + " terraform state storage account " + tf.stateStorageAccountName + "...")
	_, err := tf.executor.Execute("az storage account create" +
		" --name " + tf.stateStorageAccountName +
		" --resource-group " + tf.resourceGroupName +
		" --location " + tf.region +
		" --sku Standard_LRS" +
		" --access-tier Hot" +
		" --require-infrastructure-encryption" + // infra encryption will add another layer of encryption at rest
		" --kind StorageV2" +
		" --min-tls-version TLS1_2" +
		" --https-only true")

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	logrus.Info("Reading the storage account key, which will be give to terraform to initialize the remote state...")
	storageAccountKey, err := tf.executor.ExecuteSilent("az storage account keys list" +
		" --resource-group " + tf.resourceGroupName +
		" --account-name " + tf.stateStorageAccountName +
		" --query [0].value -o tsv")

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	logrus.Info("Creating the remote state blob container named " + tf.storageSettings.BlobContainerName + "...")
	_, err = tf.executor.Execute("az storage container create" +
		" --account-name " + tf.stateStorageAccountName +
		" --account-key " + storageAccountKey +
		" --name " + tf.storageSettings.BlobContainerName)

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	if tf.deploymentSettings.AlwaysCleanLocalCache {
		logrus.Info("Clearing the terraform cache...")
		err1 := os.RemoveAll(filepath.Join(tf.terraformDirectory, ".terraform"))
		err2 := os.RemoveAll(filepath.Join(tf.terraformDirectory, ".terraform.lock.hcl"))
		err3 := os.RemoveAll(filepath.Join(tf.terraformDirectory, tf.GetVariablesFileName()))

		if err1 != nil || err2 != nil || err3 != nil {
			return internal.ReturnErrorOrPanic(fmt.Errorf("errors while clearing terraform cache: %v %v %v", err1, err2, err3))
		}
	}

	logrus.Info("Terraform init...")
	_, err = tf.executor.Execute("terraform" +
		" -chdir=" + tf.terraformDirectory +
		" init -upgrade " +
		" --backend-config=\"subscription_id=" + tf.subscriptionId + "\"" +
		" --backend-config=\"tenant_id=" + tf.tenantId + "\"" +
		" --backend-config=\"storage_account_name=" + tf.stateStorageAccountName + "\"" +
		" --backend-config=\"access_key=" + storageAccountKey + "\"" +
		" --backend-config=\"container_name=" + tf.storageSettings.BlobContainerName + "\"" +
		" --backend-config=\"key=" + tf.storageSettings.BlobContainerKey + "\"")

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	return nil
}

func (tf *terraformWrapper) SetVariables(terraformVariables map[string]interface{}) error {
	logrus.Info("Setting the terraform variables...")

	variablesPath := filepath.Join(tf.terraformDirectory, tf.GetVariablesFileName())

	// terraform will get the variables via a file on disk, so we create it and fill out the values
	f, err := os.Create(variablesPath)

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	defer f.Close()

	for key, value := range terraformVariables {
		valueJson, err := json.Marshal(value)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}

		_, err = f.WriteString(fmt.Sprintf("%s=%s\n", key, string(valueJson)))

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	}

	tf.variablesSet = true
	return nil
}

func (tf *terraformWrapper) DeployFlow(planOnly bool, useExistingPlan bool) error {
	return tf.applyFlow(false, planOnly, useExistingPlan)
}

func (tf *terraformWrapper) DestroyFlow(planOnly bool, useExistingPlan bool) error {
	return tf.applyFlow(true, planOnly, useExistingPlan)
}

func (tf *terraformWrapper) PlanDeploy() (string, error) {
	return tf.plan(false)
}

func (tf *terraformWrapper) PlanDestroy() (string, error) {
	return tf.plan(true)
}

func (tf *terraformWrapper) ForceDeploy() error {
	return tf.forceApply(false)
}

func (tf *terraformWrapper) ForceDestroy() error {
	return tf.forceApply(true)
}

func (tf *terraformWrapper) GetBackendStorageSettings() *BackendStorageSettings {
	return &tf.storageSettings
}

func (tf *terraformWrapper) GetDeploymentSettings() *DeploymentSettings {
	return &tf.deploymentSettings
}

func (tf *terraformWrapper) GetDeployPlanFileName() string {
	return tf.projectName + ".deploy.tfplan"
}

func (tf *terraformWrapper) GetDestroyPlanFileName() string {
	return tf.projectName + ".destroy.tfplan"
}

func (tf *terraformWrapper) GetVariablesFileName() string {
	return tf.projectName + ".tfvars"
}

func (tf *terraformWrapper) guardAgainstUnsetVariables() error {
	if !tf.variablesSet {
		return errors.New("you should call SetVariables() before executing any of the terraform functions")
	}

	return nil
}

func (tf *terraformWrapper) plan(isDestroy bool) (string, error) {
	if isDestroy {
		logrus.Info("Creating the terraform destroy plan...")
	} else {
		logrus.Info("Creating the terraform deployment plan...")
	}

	err := tf.guardAgainstUnsetVariables()
	if err != nil {
		return "", internal.ReturnErrorOrPanic(err)
	}

	tfCommand := "terraform" +
		" -chdir=" + tf.terraformDirectory +
		" plan -input=false " +
		" -var-file=\"" + tf.GetVariablesFileName() + "\""

	if isDestroy {
		tfCommand += " -destroy"
		tfCommand += " -out=\"" + tf.GetDestroyPlanFileName() + "\""
	} else {
		tfCommand += " -out=\"" + tf.GetDeployPlanFileName() + "\""
	}

	planOutput, err := tf.executor.Execute(tfCommand)

	if err != nil {
		return "", internal.ReturnErrorOrPanic(err)
	}

	return planOutput, nil
}

func (tf *terraformWrapper) applyFlow(isDestroy bool, planOnly bool, useExistingPlan bool) error {
	var err error
	var plan string

	if planOnly && useExistingPlan {
		return errors.New("planOnly with useExistingPlan makes no sense as a combination")
	}

	if useExistingPlan {
		if isDestroy {
			err = tf.ForceDestroy()
		} else {
			err = tf.ForceDeploy()
		}

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	} else {
		if isDestroy {
			plan, err = tf.PlanDestroy()
		} else {
			plan, err = tf.PlanDeploy()
		}

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}

		// we show the plan to the user, but since the command output already logged it to the file, it is enough to pipe it
		// to console directly, so that duplicate log line in the file is avoided
		fmt.Println(plan)

		if !planOnly {
			if tf.executor.AskUserToConfirm("Do you want to apply the plan?") {
				if isDestroy {
					err = tf.ForceDestroy()
				} else {
					err = tf.ForceDeploy()
				}

				if err != nil {
					return internal.ReturnErrorOrPanic(err)
				}
			} else {
				logrus.Info("Plan was not approved.")
			}
		}
	}

	return nil
}

func (tf *terraformWrapper) forceApply(isDestroy bool) error {
	if isDestroy {
		logrus.Info("Starting the terraform destroy...")
	} else {
		logrus.Info("Starting the terraform apply")
	}

	err := tf.guardAgainstUnsetVariables()

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	tfCommand := "terraform" +
		" -chdir=" + tf.terraformDirectory +
		" apply" +
		" -auto-approve -input=false"

	if isDestroy {
		tfCommand += " -destroy"
		tfCommand += " \"" + tf.GetDestroyPlanFileName() + "\""
	} else {
		tfCommand += " \"" + tf.GetDeployPlanFileName() + "\""
	}

	_, err = tf.executor.Execute(tfCommand)

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	return nil
}
