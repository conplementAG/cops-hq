package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// SetVariables is required, before any of the following methods, like PlanDeploy or Deploy are called. Variables set
	// will be applied on any subsequent operation. Parameters are:
	//     terraformVariables (this is a map of terraform variables, defined as string keys, and serialized via json serializer,
	//         which means any complex or simple type is supported)
	SetVariables(terraformVariables map[string]interface{}) error

	// PlanDeploy executes the terraform plan for deployment, returning the changes as a string. Plan output is always
	// saved to a file as well, which is mandatory for Deploy. Common pattern is to show the changes to the user, ask
	// for confirmation, and then to Deploy the plan.
	PlanDeploy() (string, error)

	// PlanDestroy executes the terraform plan for destroy, returning the changes as a string. Plan output is always
	// saved to a file as well, which is mandatory for Destroy. Common pattern is to show the changes to the user, ask
	// for confirmation, and then to Destroy the plan.
	PlanDestroy() (string, error)

	// Deploy deploys the plan persisted on disk via PlanDeploy
	Deploy() error

	// Destroy deploys the plan persisted on disk via PlanDestroy
	Destroy() error

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
			return err
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
		return err
	}

	logrus.Info("Reading the storage account key, which will be give to terraform to initialize the remote state...")
	storageAccountKey, err := tf.executor.ExecuteSilent("az storage account keys list" +
		" --resource-group " + tf.resourceGroupName +
		" --account-name " + tf.stateStorageAccountName +
		" --query [0].value -o tsv")

	if err != nil {
		return err
	}

	logrus.Info("Creating the remote state blob container named " + tf.storageSettings.BlobContainerName + "...")
	_, err = tf.executor.Execute("az storage container create" +
		" --account-name " + tf.stateStorageAccountName +
		" --account-key " + storageAccountKey +
		" --name " + tf.storageSettings.BlobContainerName)

	if err != nil {
		return err
	}

	if tf.deploymentSettings.AlwaysCleanLocalCache {
		logrus.Info("Clearing the terraform cache...")
		err1 := os.RemoveAll(filepath.Join(tf.terraformDirectory, ".terraform"))
		err2 := os.RemoveAll(filepath.Join(tf.terraformDirectory, ".terraform.lock.hcl"))
		err3 := os.RemoveAll(filepath.Join(tf.terraformDirectory, tf.GetVariablesFileName()))

		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("errors while clearing terraform cache: %v %v %v", err1, err2, err3)
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
		return err
	}

	return nil
}

func (tf *terraformWrapper) SetVariables(terraformVariables map[string]interface{}) error {
	logrus.Info("Setting the terraform variables...")

	// add the default variables which should always be available
	err := tf.addDefaultVariables(terraformVariables)

	if err != nil {
		return err
	}

	variablesPath := filepath.Join(tf.terraformDirectory, tf.GetVariablesFileName())

	f, err := os.Create(variablesPath)

	if err != nil {
		return err
	}

	defer f.Close()

	for key, value := range terraformVariables {
		valueJson, err := json.Marshal(value)

		if err != nil {
			return err
		}

		_, err = f.WriteString(fmt.Sprintf("%s=%s\n", key, string(valueJson)))

		if err != nil {
			return err
		}
	}

	tf.variablesSet = true
	return nil
}

func (tf *terraformWrapper) PlanDeploy() (string, error) {
	logrus.Info("Calculating terraform deployment plan...")
	err := tf.guardAgainstUnsetVariables()

	if err != nil {
		return "", err
	}

	planOutput, err := tf.executor.Execute("terraform" +
		" -chdir=" + tf.terraformDirectory +
		" plan -input=false " +
		" -var-file=\"" + tf.GetVariablesFileName() + "\"" +
		" -out=\"" + tf.GetDeployPlanFileName() + "\"")

	if err != nil {
		return "", err
	}

	return planOutput, err
}

func (tf *terraformWrapper) PlanDestroy() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (tf *terraformWrapper) Deploy() error {
	//TODO implement me
	panic("implement me")
}

func (tf *terraformWrapper) Destroy() error {
	//TODO implement me
	panic("implement me")
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

func (tf *terraformWrapper) addDefaultVariables(terraformVariables map[string]interface{}) error {
	if _, ok := terraformVariables["subscription_id"]; ok {
		return errors.New("key subscription_id was already set in terraform variables. This was unexpected, as this variable " +
			"should be set by cops-hq. Please rename your variable to something else")
	}

	if _, ok := terraformVariables["tenant_id"]; ok {
		return errors.New("key tenant_id was already set in terraform variables. This was unexpected, as this variable " +
			"should be set by cops-hq. Please rename your variable to something else")
	}

	if _, ok := terraformVariables["region"]; ok {
		return errors.New("key region was already set in terraform variables. This was unexpected, as this variable " +
			"should be set by cops-hq. Please rename your variable to something else")
	}

	terraformVariables["subscription_id"] = tf.subscriptionId
	terraformVariables["tenant_id"] = tf.tenantId
	terraformVariables["region"] = tf.region

	return nil
}
