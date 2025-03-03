package azure_login

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/sirupsen/logrus"
)

type Login struct {
	servicePrincipalId                  string
	servicePrincipalSecret              string
	userAssignedManagedIdentityClientId string
	useManagedIdentity                  bool
	tenant                              string
	executor                            commands.Executor
}

// Login logs the currently configured user in AzureCLI and Terraform
//
// Attempts the login in the following order if configured:
//   - user assigned managed identity
//   - system assigned managed identity
//   - service principal
//   - normal user login
func (l *Login) Login() error {
	if l.useUserAssignedManagedIdentityLogin() {
		if l.tenant == "" {
			return errors.New("tenant must be given, when using user assigned managed identity")
		}

		logrus.Info("Login as user assigned managed identity: " + l.userAssignedManagedIdentityClientId)
		err := l.userAssignedManagedIdentityLogin(l.userAssignedManagedIdentityClientId, l.tenant)
		return internal.ReturnErrorOrPanic(err)
	} else if l.useSystemAssignedManagedIdentityLogin() {
		if l.tenant == "" {
			return errors.New("tenant must be given, when using system assigned managed identity")
		}

		logrus.Info("Login as system assigned managed identity")
		err := l.systemAssignedManagedIdentityLogin(l.tenant)
		return internal.ReturnErrorOrPanic(err)
	} else if l.useServicePrincipalLogin() {
		if l.servicePrincipalSecret == "" {
			return internal.ReturnErrorOrPanic(errors.New("service principal secret must be given, when using service principal credentials"))
		}

		if l.tenant == "" {
			return errors.New("tenant must be given, when using service principal credentials")
		}

		logrus.Info("Login as service-principal: " + l.servicePrincipalId)
		err := l.servicePrincipalLogin(l.servicePrincipalId, l.servicePrincipalSecret, l.tenant)
		return internal.ReturnErrorOrPanic(err)
	} else {
		loggedIn, err := l.isUserAlreadyLoggedIn()

		if err != nil {
			logrus.Debug("Checking user already logged in returned error: " + err.Error() + ". " +
				"Will will to re-login the user.")
		}

		if !loggedIn {
			logrus.Info("Login as user interactive")
			return internal.ReturnErrorOrPanic(l.interactiveLogin())
		} else {
			logrus.Info("User is already logged in")
		}
	}

	return nil
}

// SetSubscription sets the current Azure subscription on the running system (for Azure CLI & Terraform)
func (l *Login) SetSubscription(subscriptionId string) error {
	logrus.Info("Setting current Azure subscription to: " + subscriptionId)
	_, err := l.executor.Execute("az account set -s " + subscriptionId)

	errEnvVar := os.Setenv("ARM_SUBSCRIPTION_ID", subscriptionId)

	if err != nil || errEnvVar != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("errors while setting the subscription: %v %v ",
			err, errEnvVar))
	}

	return nil
}

func (l *Login) useSystemAssignedManagedIdentityLogin() bool {
	return l.useManagedIdentity && l.userAssignedManagedIdentityClientId == ""
}

func (l *Login) useUserAssignedManagedIdentityLogin() bool {
	return l.useManagedIdentity && l.userAssignedManagedIdentityClientId != ""
}

func (l *Login) useServicePrincipalLogin() bool {
	return l.servicePrincipalId != ""
}

func (l *Login) interactiveLogin() error {
	_, err := l.executor.ExecuteLoud("az login")
	return err
}

func (l *Login) servicePrincipalLogin(servicePrincipal string, secret string, tenant string) error {
	// First, we log into the Azure CLI
	// see https://learn.microsoft.com/en-us/cli/azure/reference-index?view=azure-cli-latest#az-login hints for secrets starting with "-"
	commandText := "az login -u " + servicePrincipal + " -p=" + secret + " -t " + tenant + " --service-principal"
	_, err := l.executor.ExecuteSilent(commandText)

	// Then, we also need to set the env variables required for Terraform if working with service principals
	err1 := os.Setenv("ARM_CLIENT_ID", servicePrincipal)
	err2 := os.Setenv("ARM_CLIENT_SECRET", secret)
	err3 := os.Setenv("ARM_TENANT_ID", tenant)

	if err != nil || err1 != nil || err2 != nil || err3 != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("errors while logging in via azure service principal: %v %v %v %v",
			err, err1, err2, err3))
	}

	return nil
}

func (l *Login) userAssignedManagedIdentityLogin(userAssignedManagedIdentityClientId string, tenant string) error {
	// First, we log into the Azure CLI
	// see https://learn.microsoft.com/en-us/cli/azure/reference-index?view=azure-cli-latest#az-login hints for secrets starting with "-"
	commandText := "az login --identity --username " + userAssignedManagedIdentityClientId
	_, err := l.executor.Execute(commandText)

	// Then, we also need to set the env variables required for Terraform if working with user assigned managed identities
	err1 := os.Setenv("ARM_CLIENT_ID", userAssignedManagedIdentityClientId)
	err2 := os.Setenv("ARM_USE_MSI", "true")
	err3 := os.Setenv("ARM_TENANT_ID", tenant)

	if err != nil || err1 != nil || err2 != nil || err3 != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("errors while logging in via user assigned managed identity: %v %v %v %v",
			err, err1, err2, err3))
	}

	return nil
}

func (l *Login) systemAssignedManagedIdentityLogin(tenant string) error {
	// First, we log into the Azure CLI
	// see https://learn.microsoft.com/en-us/cli/azure/reference-index?view=azure-cli-latest#az-login hints for secrets starting with "-"
	commandText := "az login --identity"
	_, err := l.executor.Execute(commandText)

	// Then, we also need to set the env variables required for Terraform if working with system assigned managed identities
	err1 := os.Setenv("ARM_USE_MSI", "true")
	err2 := os.Setenv("ARM_TENANT_ID", tenant)

	if err != nil || err1 != nil || err2 != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("errors while logging in via system assigned managed identity: %v %v %v",
			err, err1, err2))
	}

	return nil
}

func (l *Login) isUserAlreadyLoggedIn() (bool, error) {
	// since we actually rely on errors to test if user is logged in, we will shortly suppress the executor panics
	previousPanicSetting := error_handling.PanicOnAnyError
	error_handling.PanicOnAnyError = false

	output, err := l.executor.ExecuteSilent("az account show")

	error_handling.PanicOnAnyError = previousPanicSetting

	if err != nil {
		return false, err
	}

	var response account
	err = json.Unmarshal([]byte(output), &response)

	if err != nil {
		return false, err
	}

	// case-insensitive comparison because Azure CLI is known to introduce these breaking changes sometimes
	return strings.EqualFold(response.User.Type, "user"), nil
}

type account struct {
	User struct {
		Type string `json:"type"`
	} `json:"user"`
}
