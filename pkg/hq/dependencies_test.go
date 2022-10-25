package hq

import (
	"encoding/json"
	"fmt"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestHQ_CheckToolingDependencies_ShouldPassWhenDependenciesSatisfied(t *testing.T) {
	// Arrange
	executorMock := &versionCheckExecutorMock{}
	executorMock.SetVersionsToExpected()
	executorMock.On("Execute", mock.Anything).Maybe()

	hq := New("hq", "0.0.1", "test.logs")
	hq.(*hqContainer).Executor = executorMock

	// Act & Assert
	error := hq.CheckToolingDependencies()
	assert.NoError(t, error)
	executorMock.AssertExpectations(t)
}

func TestHQ_CheckToolingDependencies_ShouldFailWhenDependencyOutOfDate(t *testing.T) {
	// Arrange
	executorMock := &versionCheckExecutorMock{}
	executorMock.SetVersionsToExpected()
	executorMock.azureCliVersion = "2.15.0"
	executorMock.On("Execute", mock.Anything).Maybe()

	hq := New("hq", "0.0.1", "test.logs")
	hq.(*hqContainer).Executor = executorMock

	// Act & Assert
	err := hq.CheckToolingDependencies()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "azure cli version mismatch")
	}
	executorMock.AssertExpectations(t)
}

type versionCheckExecutorMock struct {
	mock.Mock
	commands.Executor

	azureCliVersion  string
	terraformVersion string
	kubectlVersion   string
	helmVersion      string
	copsctlVersion   string
}

func (e *versionCheckExecutorMock) SetVersionsToExpected() {
	e.azureCliVersion = ExpectedMinAzureCliVersion
	e.terraformVersion = ExpectedMinTerraformVersion
	e.kubectlVersion = ExpectedMinKubectlVersion
	e.helmVersion = ExpectedMinHelmVersion
	e.copsctlVersion = ExpectedMinCopsctlVersion
}

func (e *versionCheckExecutorMock) Execute(command string) (string, error) {
	e.Called(command)

	if strings.Contains(command, "az") {
		response := map[string]string{
			"azure-cli": e.azureCliVersion,
		}
		return serializeToJson(response), nil
	}

	if strings.Contains(command, "terraform") {
		response := map[string]string{
			"terraform_version": e.terraformVersion,
		}
		return serializeToJson(response), nil
	}

	if strings.Contains(command, "kubectl") {
		response := new(kubectlVersionResponse)
		response.ClientVersion.GitVersion = e.kubectlVersion

		return serializeToJson(response), nil
	}

	if strings.Contains(command, "helm") {
		return e.helmVersion, nil
	}

	if strings.Contains(command, "copsctl") {
		return e.copsctlVersion, nil
	}

	// we explicitly don't set sops as installed, because it should just issue a warning
	return "unknown command for the Execute mock called, but let's return successfully anyways", nil
}

func serializeToJson(input interface{}) string {
	encoded, err := json.Marshal(input)

	if err != nil {
		panic(fmt.Errorf("problem serializing to json: " + err.Error()))
	}

	return string(encoded)
}
