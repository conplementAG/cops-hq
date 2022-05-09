package azure_login

import (
	"encoding/json"
	"errors"
	"github.com/denisbiondic/cops-hq/pkg/commands"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func Test_TriggersServicePrincipalLogin_WhenIdProvided(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	azureLogin := NewWithParams(executor, "abcd", "secret", "tenantId")

	executor.On("ExecuteSilent", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "--service-principal") && strings.Contains(command, "abcd")
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

	executor.On("Execute", mock.MatchedBy(func(command string) bool {
		return !strings.Contains(command, "az login")
	}))

	// Act
	azureLogin.Login()

	// Assert
	executor.AssertExpectations(t)
}

func Test_TriggersUserLogin_WhenNoCredentialsProvidedAndNotLoggedIn(t *testing.T) {
	// Arrange
	executor := &loginExecutorMock{}
	executor.userLoggedIn = false

	azureLogin := New(executor)

	// sadly, testify expects that every call made has a matching expect, so we need to add this comamnd as well,
	// although it makes the test just more brittle
	executor.On("Execute", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "az account show")
	})).Once()

	executor.On("Execute", mock.MatchedBy(func(command string) bool {
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
	}

	return "unknown command for the Execute mock called, but let's return successfully anyways", nil
}

func (e *loginExecutorMock) ExecuteSilent(command string) (string, error) {
	e.Called(command)

	if strings.Contains(command, "--service-principal") {
		return "sp login output", nil
	} else if strings.Contains(command, "az account show") {
		a := &account{}
		a.User.Type = "User"
		b, err := json.Marshal(a)
		return string(b), err
	}

	return "unknown command for the ExecuteSilent mock called, but let's return successfully anyways", nil
}
