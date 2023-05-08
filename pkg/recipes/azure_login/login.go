package azure_login

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Login struct {
	servicePrincipalId     string
	servicePrincipalSecret string
	tenant                 string
	executor               commands.Executor
}

// Login logs the currently configured user in AzureCLI and Terraform. If configured with service principal, it will
// attempt a non-interactive login, otherwise a normal user login will be started.
func (l *Login) Login() error {
	if l.useServicePrincipalLogin() {
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

func (l *Login) useServicePrincipalLogin() bool {
	return l.servicePrincipalId != ""
}

func (l *Login) interactiveLogin() error {
	_, err := l.executor.ExecuteLoud("az login")
	return err
}

func (l *Login) servicePrincipalLogin(servicePrincipal string, secret string, tenant string) error {
	// First, we log into the Azure CLI
	commandText := "az login -u " + servicePrincipal + " -p " + secret + " -t " + tenant + " --service-principal"
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

func (l *Login) isUserAlreadyLoggedIn() (bool, error) {
	// since we actually rely on errors to test if user is logged in, we will shortly supress the executor panics
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
