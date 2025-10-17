package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/internal/cmdutil"
	"github.com/conplementag/cops-hq/v2/internal/file_handling"
	"github.com/conplementag/cops-hq/v2/internal/slice_helpers"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform/file_paths"
	"github.com/sirupsen/logrus"
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
	//     autoApprove (will skip asking user questions, if required, like approving a plan before apply. Defaults to false.
	//                  Can be safely set to true on operations which don't prompt any user inputs, it will just have zero effect on the behaviour).
	// Here is a short explanation:
	//     DeployFlow(false, false, false) will show the plan, prompt the user, and apply if confirmed (setting autoApprove to true will skip confirmation)
	//     DeployFlow(true, false, false) will only show the plan, but also persist it on disk (use GetDeployPlanFileName() for details)
	//     DeployFlow(false, true, false) will reuse the plan already saved on disk, and apply it without any user confirmations (autoApprove makes no impact here)
	//     DeployFlow(true, true, false) will just show the plan persisted on the disk, without generating a new plan.
	// It is best practice to set the both planOnly and useExistingPlan from the CLI, so that CI scripts can simply override
	// the variables depending on the current CI step (usually a plan is presented, user is awaited for approval, then the existing
	// plan is applied). The autoApprove parameter is useful in local deployment scenarios, where you plan / deploy everything as one step,
	// and perhaps do not want to be prompted.
	DeployFlow(planOnly bool, useExistingPlan bool, autoApprove bool) error

	// DestroyFlow is same as DeployFlow, but only for destroy.
	DestroyFlow(planOnly bool, useExistingPlan bool, autoApprove bool) error

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

	// GetVariablesFileName returns the file name in which the terraform variables will be stored. This name is convention based
	// on the currently set project parameter while creating the terraform wrapper instance
	GetVariablesFileName() string

	// Output returns the terraform provided output if any
	// Parameters are:
	// 		parameterName will restrict output to single output parameter and output the parameter in raw mode. If not set all available output is given
	Output(parameterName *string) (string, error)

	// OutputAsJson returns the terraform provided output as json string
	//
	// CAUTION: Returns sensitive values in plain text
	OutputAsJson() (string, error)
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

	tags := serializeTagsIntoCmdArgsList(tf.storageSettings.Tags)

	if tf.storageSettings.CreateResourceGroup {
		logrus.Info("Deploying the project " + tf.projectName + " resource group " + tf.resourceGroupName + "...")

		groupCreateCmd := exec.Command("az", "group", "create",
			"-l", tf.region,
			"-n", tf.resourceGroupName,
			"--tags")

		if len(tags) > 0 {
			groupCreateCmd.Args = append(groupCreateCmd.Args, tags...)
		} else {
			// reset tags, if there were already any assigned
			groupCreateCmd.Args = append(groupCreateCmd.Args, "")
		}

		// using ExecuteCmd to skip the escaping logic of "Execute"
		// the --tags argument does not follow the usual --argument value semantics (value(s) contain equal (=) sign and could also contain spaces)
		_, err := tf.executor.ExecuteCmd(groupCreateCmd)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	}

	logrus.Info("Deploying the " + tf.projectName + " terraform state storage account " + tf.stateStorageAccountName + "...")

	defaultAction := "Allow"
	if len(tf.storageSettings.AllowedIpAddresses) > 0 {
		defaultAction = "Deny"
	}

	storageAccountCreateCmd := exec.Command("az", "storage", "account", "create",
		"--name", tf.stateStorageAccountName,
		"--resource-group", tf.resourceGroupName,
		"--location", tf.region,
		"--default-action", defaultAction,
		"--sku", "Standard_LRS",
		"--access-tier", "Hot",
		"--kind", "StorageV2",
		"--min-tls-version", "TLS1_2",
		"--https-only", "true",
		"--tags")

	if len(tags) > 0 {
		storageAccountCreateCmd.Args = append(storageAccountCreateCmd.Args, tags...)
	} else {
		// reset tags, if there were already any assigned
		storageAccountCreateCmd.Args = append(storageAccountCreateCmd.Args, "")
	}

	if tf.storageSettings.RequireInfrastructureEncryption {
		storageAccountCreateCmd.Args = append(storageAccountCreateCmd.Args, "--require-infrastructure-encryption") // infra encryption will add another layer of encryption at rest
	}

	// using ExecuteCmd to skip the escaping logic of "Execute"
	// the --tags argument does not follow the usual --argument value semantics (value(s) contain equal (=) sign and could also contain spaces)
	_, err := tf.executor.ExecuteCmd(storageAccountCreateCmd)

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	err = tf.addStorageAccountNetworkRules()
	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	defaultPortalAccessCmd := exec.Command("az", "storage", "account", "update",
		"--name", tf.stateStorageAccountName,
		"--resource-group", tf.resourceGroupName,
		"--set", "defaultToOAuthAuthentication=true")

	_, err = tf.executor.ExecuteCmd(defaultPortalAccessCmd)
	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	fileSharePropertiesCmd := exec.Command("az", "storage", "account", "file-service-properties", "update",
		"--account-name", tf.stateStorageAccountName,
		"--resource-group", tf.resourceGroupName,
		"--versions", "SMB3.1.1",
		"--auth-methods", "Kerberos",
		"--kerb-ticket-encryption", "AES-256",
		"--channel-encryption", "AES-256-GCM")

	_, err = tf.executor.ExecuteCmd(fileSharePropertiesCmd)
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

	storageAccountKey = trimLinebreakSuffixes(storageAccountKey)

	logrus.Info("Creating the remote state blob container named " + tf.storageSettings.BlobContainerName + "...")
	err = cmdutil.ExecuteWithRetry(
		tf.executor.ExecuteSilent,
		"az storage container create"+
			" --account-name "+tf.stateStorageAccountName+
			" --account-key "+storageAccountKey+
			" --name "+tf.storageSettings.BlobContainerName,
		tf.storageSettings.ContainerCreateRetryCount)

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
		" --backend-config=subscription_id=" + tf.subscriptionId +
		" --backend-config=tenant_id=" + tf.tenantId +
		" --backend-config=storage_account_name=" + tf.stateStorageAccountName +
		" --backend-config=access_key=" + storageAccountKey +
		" --backend-config=container_name=" + tf.storageSettings.BlobContainerName +
		" --backend-config=key=" + tf.storageSettings.BlobContainerKey)

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

func (tf *terraformWrapper) DeployFlow(planOnly bool, useExistingPlan bool, autoApprove bool) error {
	return tf.applyFlow(false, planOnly, useExistingPlan, autoApprove)
}

func (tf *terraformWrapper) DestroyFlow(planOnly bool, useExistingPlan bool, autoApprove bool) error {
	return tf.applyFlow(true, planOnly, useExistingPlan, autoApprove)
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

func (tf *terraformWrapper) GetVariablesFileName() string {
	return tf.projectName + ".tfvars"
}

func (tf *terraformWrapper) Output(parameterName *string) (string, error) {

	tfCommand := "terraform" +
		" -chdir=" + tf.terraformDirectory +
		" output"

	if parameterName != nil {
		tfCommand += " -raw " + *parameterName
	}

	return tf.executor.Execute(tfCommand)
}

func (tf *terraformWrapper) OutputAsJson() (string, error) {

	tfCommand := "terraform" +
		" -chdir=" + tf.terraformDirectory +
		" output" +
		" -json"

	return tf.executor.ExecuteSilent(tfCommand)
}

func (tf *terraformWrapper) addStorageAccountNetworkRules() error {
	existingIpAddresses, err := tf.determineCurrentAllowedIpAddresses()
	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	addIpAddresses, removeIpAddresses := slice_helpers.FindItemsToAddAndRemove(existingIpAddresses, tf.storageSettings.AllowedIpAddresses)

	// add new rules
	for _, ipAddress := range addIpAddresses {
		err = tf.addOrRemoveStorageAccountNetworkRule("add", ipAddress)
		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	}

	// remove rules
	for _, ipAddress := range removeIpAddresses {
		err = tf.addOrRemoveStorageAccountNetworkRule("remove", ipAddress)
		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	}

	// ensure rules are applied to allow further processing
	retryErrorText := "network rules not equal"
	err = cmdutil.ExecuteFunctionWithRetry(
		func() error {
			currentAllowedIpAddresses, _ := tf.determineCurrentAllowedIpAddresses()
			if reflect.DeepEqual(currentAllowedIpAddresses, tf.storageSettings.AllowedIpAddresses) {
				return nil
			} else {
				return errors.New(retryErrorText)
			}
		}, tf.storageSettings.ContainerCreateRetryCount)
	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	return nil
}

func (tf *terraformWrapper) addOrRemoveStorageAccountNetworkRule(method, value string) error {
	cmd := "az storage account network-rule " + method +
		" --resource-group " + tf.resourceGroupName +
		" --account-name " + tf.stateStorageAccountName + " " +
		" --ip-address " + value

	_, err := tf.executor.Execute(cmd)
	return internal.ReturnErrorOrPanic(err)
}

func (tf *terraformWrapper) determineCurrentAllowedIpAddresses() ([]string, error) {
	networkRuleListCmd := "az storage account network-rule list" +
		" --resource-group " + tf.resourceGroupName +
		" --account-name " + tf.stateStorageAccountName +
		" --query ipRules[].ipAddressOrRange -o json"

	currentAllowedIpAddresses, err := tf.executor.Execute(networkRuleListCmd)

	if err != nil {
		return []string{}, err
	}

	var currentAllowedIpAddressesMapped []string
	err = json.Unmarshal([]byte(currentAllowedIpAddresses), &currentAllowedIpAddressesMapped)

	if err != nil {
		return []string{}, err
	}

	return currentAllowedIpAddressesMapped, nil
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
		" -var-file=" + tf.GetVariablesFileName() +
		" -detailed-exitcode"

	var localTerraformRelativePlanFilePath string

	if isDestroy {
		tfCommand += " -destroy"
		localTerraformRelativePlanFilePath, err = file_paths.GetLocalTerraformRelativePlanFilePath(tf.projectName, tf.terraformDirectory, true)
		if err != nil {
			return "", internal.ReturnErrorOrPanic(err)
		}

		tfCommand += " -out=" + localTerraformRelativePlanFilePath
	} else {
		localTerraformRelativePlanFilePath, err = file_paths.GetLocalTerraformRelativePlanFilePath(tf.projectName, tf.terraformDirectory, false)
		if err != nil {
			return "", internal.ReturnErrorOrPanic(err)
		}

		tfCommand += " -out=" + localTerraformRelativePlanFilePath
	}

	// before creating a new plan, make sure all files for this project are removed
	err = file_handling.DeleteFilesStartingWith(
		file_paths.GetPlanFileName(tf.projectName, isDestroy),
		filepath.Join(tf.terraformDirectory, file_paths.PlansDirectory))
	if err != nil {
		return "", internal.ReturnErrorOrPanic(err)
	}

	// disable global setting panic on error to allow terraform plan with
	// detailed-exitcode for plan analyzing
	panicOnError := error_handling.PanicOnAnyError
	error_handling.PanicOnAnyError = false
	defer func(panicOnError bool) {
		error_handling.PanicOnAnyError = panicOnError
	}(panicOnError)

	plaintextPlanOutput, err := tf.executor.Execute(tfCommand)
	// terraform plan with -detailed-exitcode results in the following exit codes
	// 0 = Succeeded with empty diff (no changes)
	// 1 = Error
	// 2 = Succeeded with non-empty diff (changes present)
	var planIsDirty bool
	switch exitCode := getExitCode(err); exitCode {
	case 0:
		planIsDirty = false
		break
	case 1:
		return "", internal.ReturnErrorOrPanic(err)
	case 2:
		planIsDirty = true
		break
	default:
		return "", internal.ReturnErrorOrPanic(fmt.Errorf("unexpected exit code %d in terraform plan command %w", exitCode, err))
	}

	err = tf.persistPlanInAdditionalFormatsOnDisk(plaintextPlanOutput, localTerraformRelativePlanFilePath)
	if err != nil {
		return "", internal.ReturnErrorOrPanic(err)
	}

	err = tf.persistAnalysisResultOnDisk(localTerraformRelativePlanFilePath, isDestroy, planIsDirty)
	if err != nil {
		return "", internal.ReturnErrorOrPanic(err)
	}

	return plaintextPlanOutput, nil
}

func (tf *terraformWrapper) applyFlow(isDestroy bool, planOnly bool, useExistingPlan bool, autoApprove bool) error {
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
			approved := false

			if autoApprove {
				approved = true
			}

			if !approved {
				approved = tf.executor.AskUserToConfirm("Do you want to apply the plan?")
			}

			if !approved {
				err := errors.New("plan was not approved")
				logrus.Error(err)
				return internal.ReturnErrorOrPanic(err)
			}

			if isDestroy {
				err = tf.ForceDestroy()
			} else {
				err = tf.ForceDeploy()
			}

			if err != nil {
				return internal.ReturnErrorOrPanic(err)
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
		path, err := file_paths.GetLocalTerraformRelativePlanFilePath(tf.projectName, tf.terraformDirectory, true)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}

		tfCommand += " \"" + path + "\""
	} else {
		path, err := file_paths.GetLocalTerraformRelativePlanFilePath(tf.projectName, tf.terraformDirectory, false)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}

		tfCommand += " \"" + path + "\""
	}

	_, err = tf.executor.Execute(tfCommand)

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	return nil
}

func trimLinebreakSuffixes(storageAccountKey string) string {
	return strings.TrimRight(storageAccountKey, "\r\n")
}

func getExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}
