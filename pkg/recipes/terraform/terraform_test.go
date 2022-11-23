package terraform

import (
	"errors"
	"fmt"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_SettingsAreNotSharedByReferenceBetweenMultipleInstances(t *testing.T) {
	// we want to make sure that the design of non-shared settings is always kept, even after refactorings
	// Arrange
	tf1, _ := createSimpleTerraformWithDefaultSettings("app1")
	tf2, _ := createSimpleTerraformWithDefaultSettings("app2")

	// Act
	tf1.GetDeploymentSettings().AlwaysCleanLocalCache = false
	tf1.GetBackendStorageSettings().BlobContainerKey = "acme"
	tf2.GetBackendStorageSettings().BlobContainerKey = "lock"

	// Assert
	assert.True(t, tf2.GetDeploymentSettings().AlwaysCleanLocalCache)
	assert.Equal(t, "acme", tf1.GetBackendStorageSettings().BlobContainerKey)
	assert.Equal(t, "lock", tf2.GetBackendStorageSettings().BlobContainerKey)

	// we also check the defaults object was not modified by accident
	assert.True(t, DefaultDeploymentSettings.AlwaysCleanLocalCache)
}

func Test_SetVariablesCanSerializeAnySimpleOrComplexValue(t *testing.T) {
	// Arrange
	tf, _ := createSimpleTerraformWithDefaultSettings("test")

	var variables = make(map[string]interface{})

	variables["simple"] = "value"
	variables["simpleNumber"] = 3
	variables["truth"] = false
	variables["list_of_numbers"] = []int{1, 2, 3, 4}
	variables["list_of_strings"] = []string{"10.0.0.1", "20.20.20.4"}
	variables["my_complex_type"] = &variablesStruct{
		AString: "What is the answer to life?",
		AnInt:   42,
		ABool:   true,
	}

	// Act
	err := tf.SetVariables(variables)

	// Assert
	assert.NoError(t, err)

	fileBytes, err := os.ReadFile(filepath.Join(".", tf.GetVariablesFileName()))

	if err != nil {
		assert.NoError(t, err)
	}

	// the strings below are the exact terraform format
	assert.Contains(t, string(fileBytes), "simple=\"value\"")
	assert.Contains(t, string(fileBytes), "simpleNumber=3")
	assert.Contains(t, string(fileBytes), "truth=false")
	assert.Contains(t, string(fileBytes), "list_of_numbers=[1,2,3,4]")
	assert.Contains(t, string(fileBytes), "list_of_strings=[\"10.0.0.1\",\"20.20.20.4\"]")
	assert.Contains(t, string(fileBytes), "my_complex_type={\"aString\":\"What is the answer to life?\",\"anInt\":42,\"aBool\":true}")
}

func Test_DeployFlow(t *testing.T) {
	tests := []struct {
		testName        string
		mockSetup       func(executor *executorMock)
		planOnly        bool
		useExistingPlan bool
		autoApprove     bool
		expectedError   error
	}{
		{"apply with existing plan",
			func(executor *executorMock) {
				// we only expect the existing plan to be applied
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "apply -auto-approve") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()
			},
			false,
			true,
			false,
			nil,
		},

		{"full apply flow",
			func(executor *executorMock) {
				// first, a plan should be executed, saving the file
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "plan -input=false") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// we expect the separate plan directory is always there (created)
				_, err := os.Stat(".plans")
				assert.NoError(t, err)

				// then the user confirmation is expected
				executor.On("AskUserToConfirm", mock.Anything).Once()

				// then the plan json will be created
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "show -json") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// then the fully apply is expected
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "apply -auto-approve") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()
			},
			false,
			false,
			false,
			nil,
		},

		{"auto approve does not prompt the user",
			func(executor *executorMock) {
				// first, a plan should be executed, saving the file
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "plan -input=false") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// then the plan json will be created
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "show -json") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// then the fully apply is expected
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "apply -auto-approve") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()
			},
			false,
			false,
			true,
			nil,
		},

		{"only execute the plan",
			func(executor *executorMock) {
				// a plan should be executed, saving the file
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "plan -input=false") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// then the plan json will be created
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "show -json") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()
			},
			true,
			false,
			false,
			nil,
		},

		{"only plan with existing plan throws error since it makes no sense",
			func(executor *executorMock) {},
			true,
			true,
			false,
			errors.New("planOnly with useExistingPlan makes no sense as a combination"),
		},
	}

	for _, tt := range tests {
		// Arrange
		fmt.Println("Executing test " + tt.testName)

		tf, executorMock := createSimpleTerraformWithDefaultSettings("test")
		tt.mockSetup(executorMock)
		tf.SetVariables(nil)

		// Act
		err := tf.DeployFlow(tt.planOnly, tt.useExistingPlan, tt.autoApprove)

		// Assert
		if tt.expectedError == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, tt.expectedError, err)
		}

		executorMock.AssertExpectations(t)
	}
}

type executorMock struct {
	mock.Mock
	commands.Executor
}

func (e *executorMock) Execute(command string) (string, error) {
	e.Called(command)
	return "success", nil
}

func (e *executorMock) ExecuteSilent(command string) (string, error) {
	e.Called(command)
	return "success", nil
}

func (e *executorMock) AskUserToConfirm(displayMessage string) bool {
	e.Called(displayMessage)
	return true
}

type variablesStruct struct {
	AString string `mapstructure:"aString" json:"aString" yaml:"aString"`
	AnInt   int    `mapstructure:"anInt" json:"anInt" yaml:"anInt"`
	ABool   bool   `mapstructure:"aBool" json:"aBool" yaml:"aBool"`
}

func createSimpleTerraformWithDefaultSettings(projectName string) (Terraform, *executorMock) {
	executor := &executorMock{}

	return New(executor, projectName, "1234", "3214",
		"westeurope", "testrg", "storeaccount",
		// important to keep the current directory configured, since some tests rely on this
		// location to verify that expected directories / files are created
		filepath.Join("."),
		DefaultBackendStorageSettings,
		DefaultDeploymentSettings), executor
}
