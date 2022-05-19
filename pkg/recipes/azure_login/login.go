package azure_login

import (
	"encoding/json"
	"errors"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/sirupsen/logrus"
	"strings"
)

type Login struct {
	servicePrincipalId     string
	servicePrincipalSecret string
	tenant                 string
	executor               commands.Executor
}

// Login logs the currently configured user in Azure. If configured with service principal, it will attempt a non-interactive login,
// otherwise a normal user login will be started.
func (l *Login) Login() error {
	if l.useServicePrincipalLogin() {
		if l.servicePrincipalSecret == "" {
			return errors.New("service principal secret must be given, when using service principal credentials")
		}

		if l.tenant == "" {
			return errors.New("tenant must be given, when using service principal credentials")
		}

		logrus.Info("Login as service-principal: " + l.servicePrincipalId)
		return l.servicePrincipalLogin(l.servicePrincipalId, l.servicePrincipalSecret, l.tenant)
	} else {
		loggedIn, err := l.isUserAlreadyLoggedIn()

		if err != nil {
			logrus.Debug("Checking user already logged in returned error: " + err.Error() + ". " +
				"Will will to re-login the user.")
		}

		if !loggedIn {
			logrus.Info("Login as user interactive")
			return l.interactiveLogin()
		} else {
			logrus.Info("User is already logged in")
		}
	}

	return nil
}

// SetSubscription sets the current Azure subscription on the running system (for Azure CLI)
func (l *Login) SetSubscription(subscriptionId string) error {
	logrus.Info("Setting current Azure subscription to: " + subscriptionId)
	_, err := l.executor.Execute("az account set -s " + subscriptionId)
	return err
}

func (l *Login) useServicePrincipalLogin() bool {
	return l.servicePrincipalId != ""
}

func (l *Login) interactiveLogin() error {
	_, err := l.executor.Execute("az login")
	return err
}

func (l *Login) servicePrincipalLogin(servicePrincipal string, secret string, tenant string) error {
	commandText := "az login -u " + servicePrincipal + " -p " + secret + " -t " + tenant + " --service-principal"
	_, err := l.executor.ExecuteSilent(commandText)
	return err
}

func (l *Login) isUserAlreadyLoggedIn() (bool, error) {
	// since we actually rely on errors to test if user is logged in, we will shortly supress the executor panics
	previousPanicSetting := l.executor.GetPanicOnAnyError()
	l.executor.SetPanicOnAnyError(false)

	output, err := l.executor.ExecuteSilent("az account show")

	l.executor.SetPanicOnAnyError(previousPanicSetting)

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

func (l *Login) setSubscription(subscription string) error {
	commandText := "az account set -s " + subscription
	_, err := l.executor.Execute(commandText)
	return err
}

type account struct {
	User struct {
		Type string `json:"type"`
	} `json:"user"`
}
