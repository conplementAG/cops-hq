package azure_login

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/stretchr/testify/mock"
)

func Test_TriggersServicePrincipalLogin_WhenIdProvided(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	azureLogin := NewWithParams(executor, "abcd", "secret", "tenantId", "", false)

	executor.On("ExecuteSilent", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "--service-principal") && strings.Contains(command, "abcd")
	}))

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertExpectations(t)
}

func Test_TriggersUserAssignedManagedIdentityLogin_WhenClientIdAndFlagProvided(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	azureLogin := NewWithParams(executor, "abcd", "secret", "tenantId", "umi-clientid", true)

	executor.On("Execute", mock.MatchedBy(func(command string) bool {
		return command == "az login --identity --username umi-clientid"
	}))

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertExpectations(t)
}

func Test_TriggersSystemAssignedManagedIdentityLogin_WhenOnlyFlagProvided(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	azureLogin := NewWithParams(executor, "abcd", "secret", "tenantId", "", true)

	executor.On("Execute", mock.MatchedBy(func(command string) bool {
		return command == "az login --identity"
	}))

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertExpectations(t)
}

func Test_TriggersServicePrincipalLogin_WhenIdProvidedAndUamIdProvidedButMiFlagNotProvided(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	azureLogin := NewWithParams(executor, "abcd", "secret", "tenantId", "umi-clientid", false)

	executor.On("ExecuteSilent", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "--service-principal") && strings.Contains(command, "abcd") && !strings.Contains(command, "--identity")
	}))

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertExpectations(t)
}

func Test_TriggersNoLogin_WhenUserAlreadyLoggedIn(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	executor.userLoggedIn = true

	azureLogin := New(executor)
	executor.On("ExecuteSilent", mock.Anything)

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertNotCalled(t, "Execute", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "az login")
	}))
	executor.AssertExpectations(t)
}

func Test_TriggersUserLogin_WhenNoCredentialsProvidedAndNotLoggedIn(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	executor.userLoggedIn = false

	azureLogin := New(executor)

	// sadly, testify expects that every call made has a matching expect, so we need to add this command as well,
	// although it makes the test just more brittle
	executor.On("ExecuteSilent", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "az account show")
	})).Once()

	executor.On("ExecuteLoud", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "az login")
	})).Once()

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertExpectations(t)
}

type loginExecutorMock struct {
	mock.Mock
	commands.Executor
	userLoggedIn bool
}

func (e *loginExecutorMock) setUserLoggedIn(userLoggedIn bool) {
	e.userLoggedIn = userLoggedIn
}

func (e *loginExecutorMock) Execute(command string) (string, error) {
	e.Called(command)

	if strings.Contains(command, "az account show") {
		a := &account{}

		if e.userLoggedIn {
			a.User.Type = "User"
			b, err := json.Marshal(a)
			return string(b), err
		} else {
			return "not logged in", errors.New("not logged int")
		}
	} else if strings.Contains(command, "--identity") {
		return "logged in with managed identity", nil
	}

	return "unknown command for the Execute mock called, but let's return successfully anyways", nil
}

func (e *loginExecutorMock) ExecuteSilent(command string) (string, error) {
	e.Called(command)

	if strings.Contains(command, "--service-principal") {
		return "sp login output", nil
	} else if strings.Contains(command, "az account show") {
		a := &account{}

		if e.userLoggedIn {
			a.User.Type = "User"
			b, err := json.Marshal(a)
			return string(b), err
		} else {
			return "not logged in", errors.New("not logged int")
		}
	}

	return "unknown command for the ExecuteSilent mock called, but let's return successfully anyways", nil
}

func (e *loginExecutorMock) ExecuteLoud(command string) (string, error) {
	e.Called(command)

	return "mock return - success", nil
}
