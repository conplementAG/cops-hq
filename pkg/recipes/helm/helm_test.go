package helm

import (
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type executorMock struct {
	mock.Mock
	commands.Executor
}

func Test_SetVariablesCanCreateVariablesFile(t *testing.T) {
	// Arrange
	h, _ := createSimpleHelmWithDefaultSettings("test")

	var variables = make(map[string]interface{})

	variables["simple"] = "value"
	variables["simpleNumber"] = 3
	variables["truth"] = false
	variables["nested"] = map[string]interface{}{
		"key1": "value1_nested",
		"key2": 3}
	variables["list_of_strings"] = []string{"10.0.0.1", "20.20.20.4"}
	variables["empty_array"] = []string{}

	// Act
	err := h.SetVariables(variables)

	// Assert
	assert.NoError(t, err)

	filePath := filepath.Join(".", h.GetVariablesOverrideFileName())
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		assert.NoError(t, err)
	}

	// the strings below are the exact terraform format
	assert.Contains(t, string(fileBytes), "simple: value")
	assert.Contains(t, string(fileBytes), "simpleNumber: 3")
	assert.Contains(t, string(fileBytes), "truth: false")
	assert.Contains(t, string(fileBytes), "nested:")
	assert.Contains(t, string(fileBytes), "key1: value1_nested")
	assert.Contains(t, string(fileBytes), "key2: 3")
	assert.Contains(t, string(fileBytes), "list_of_strings:")
	assert.Contains(t, string(fileBytes), "- 10.0.0.1")
	assert.Contains(t, string(fileBytes), "empty_array")
	assert.Contains(t, string(fileBytes), "[]")

	// Cleanup
	if existsFile(filePath) {
		os.Remove(filePath)
	}
}

func Test_DeployExecutesExpectedCommandWithOverrideValues(t *testing.T) {
	h, executorMock := createSimpleHelmWithDefaultSettings("project")
	// helm upgrade with one value file is expected
	executorMock.On("Execute", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "helm upgrade") && strings.Contains(command, "values.yaml") && strings.Contains(command, "values.override.yaml") && strings.Contains(command, "--timeout 5m0s")
	})).Once()

	helmVariables := make(map[string]interface{})
	helmVariables["test_key"] = "test_value"

	h.SetVariables(helmVariables)

	// Act
	err := h.Deploy()

	// Assert
	assert.NoError(t, err)

	executorMock.AssertExpectations(t)

	// Cleanup
	filePath := filepath.Join(".", h.GetVariablesOverrideFileName())
	if existsFile(filePath) {
		os.Remove(filePath)
	}
}

func Test_DeployExecutesExpectedCommandWithoutOverrideValues(t *testing.T) {
	h, executorMock := createSimpleHelmWithDefaultSettings("project")
	// helm upgrade with one value file is expected
	executorMock.On("Execute", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "helm upgrade") && strings.Contains(command, "values.yaml") && !strings.Contains(command, "values.override.yaml")
	})).Once()

	// Act
	err := h.Deploy()

	// Assert
	assert.NoError(t, err)

	executorMock.AssertExpectations(t)
}

func Test_DeployExecutesExpectedCommandWithWaitAndTimeout(t *testing.T) {
	var deploymentSettings = DefaultDeploymentSettings
	deploymentSettings.Wait = true
	deploymentSettings.Timeout = 2 * time.Minute
	h, executorMock := createWithDeploymentSettings("project", deploymentSettings)
	// helm upgrade with wait and timeout value is expected
	executorMock.On("ExecuteWithProgressInfo", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "helm upgrade") && strings.Contains(command, "--wait") && strings.Contains(command, "--timeout 2m0s")
	})).Once()

	// Act
	err := h.Deploy()

	// Assert
	assert.NoError(t, err)

	executorMock.AssertExpectations(t)
}

func Test_GetVariablesOverrideFileNameReturnsSomething(t *testing.T) {
	h, executorMock := createSimpleHelmWithDefaultSettings("project")
	// helm upgrade with one value file is expected
	executorMock.On("Execute", mock.MatchedBy(func(command string) bool {
		return strings.Contains(command, "helm upgrade") && strings.Contains(command, "values.yaml") && !strings.Contains(command, "values.override.yaml")
	})).Once()

	// Act
	actual := h.GetVariablesOverrideFileName()

	// Assert
	assert.NotEmpty(t, actual)
}

func (e *executorMock) Execute(command string) (string, error) {
	e.Called(command)
	return "success", nil
}

func (e *executorMock) ExecuteWithProgressInfo(command string) (string, error) {
	e.Called(command)
	return "success", nil
}

func createSimpleHelmWithDefaultSettings(projectName string) (Helm, *executorMock) {
	executor := &executorMock{}

	return New(executor, projectName+"test", projectName+"test", filepath.Join(".")), executor
}

func createWithDeploymentSettings(projectName string, settings DeploymentSettings) (Helm, *executorMock) {
	executor := &executorMock{}

	return NewWithSettings(executor, projectName+"test", projectName+"test", filepath.Join("."), settings), executor
}

func existsFile(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
