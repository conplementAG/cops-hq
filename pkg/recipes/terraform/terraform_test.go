package terraform

import (
	"errors"
	"fmt"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const projectName = "test"

func Test_SettingsAreNotSharedByReferenceBetweenMultipleInstances(t *testing.T) {
	// we want to make sure that the design of non-shared settings is always kept, even after refactorings
	// Arrange
	tf1 := createSimpleTerraformWithDefaultSettings(&executorMock{}, "app1")
	tf2 := createSimpleTerraformWithDefaultSettings(&executorMock{}, "app2")

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
	tf := createSimpleTerraformWithDefaultSettings(&executorMock{}, projectName)

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
		testName string
		// make sure the mock setup only contains the actual setup (like the On() calls) and not the additional validation
		// and verifications which have nothing to do with the mocks. Use postRunAdditionalVerifications for this. Main reason
		// for this is that mockSetup occurs before the actual test Act is performed, therefore it is too early to actually
		// verify anything.
		mockSetup                      func(executor *executorMock)
		planOnly                       bool
		useExistingPlan                bool
		planHasChanges                 bool
		autoApprove                    bool
		expectedError                  error
		postRunAdditionalVerifications func()
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
			false,
			nil,
			func() {
				assertPlanFilesPresence(t, false, false, false) // no plan should be created!
			},
		},

		{"full apply flow",
			func(executor *executorMock) {
				// a plan should be executed, saving the file
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "plan -input=false") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// the user confirmation is expected
				executor.On("AskUserToConfirm", mock.Anything).Once()

				// the plan json will also be created
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "show -json") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()

				// the fully apply is expected
				executor.On("Execute", mock.MatchedBy(func(command string) bool {
					return strings.Contains(command, "apply -auto-approve") && strings.Contains(command, "test.deploy.tfplan")
				})).Once()
			},
			false,
			false,
			true,
			false,
			nil,
			func() {
				assertPlanFilesPresence(t, true, true, false) // plan expected to be dirty!
			},
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
			false,
			true,
			nil,
			nil,
		},

		{"only plan with existing plan throws error since it makes no sense",
			func(executor *executorMock) {},
			true,
			true,
			false,
			false,
			errors.New("planOnly with useExistingPlan makes no sense as a combination"),
			nil,
		},

		{"only perform the plan",
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
			false,
			nil,
			func() {
				assertPlanFilesPresence(t, true, true, true)
			},
		},
	}

	for _, tt := range tests {
		// Arrange
		fmt.Println("Executing test " + tt.testName)

		// always clean the .plans to make sure the tests are not dependent on each other
		err := deleteDirectoryIfExists(".plans")
		assert.NoError(t, err)

		executor := &executorMock{}
		executor.planHasChanges = tt.planHasChanges

		tf := createSimpleTerraformWithDefaultSettings(executor, projectName)
		tt.mockSetup(executor)
		tf.SetVariables(nil)

		// Act
		err = tf.DeployFlow(tt.planOnly, tt.useExistingPlan, tt.autoApprove)

		// Assert
		if tt.expectedError == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, tt.expectedError, err)
		}

		executor.AssertExpectations(t)

		if tt.postRunAdditionalVerifications != nil {
			tt.postRunAdditionalVerifications()
		}
	}
}

func Test_MultiplePlansCleanupCorrectly(t *testing.T) {
	// we should start from a clean slate
	err := deleteDirectoryIfExists(".plans")
	assert.NoError(t, err)

	// first run, plan has no changes, all files are expected
	executor1 := &executorMock{isLooseMock: true, planHasChanges: false}
	tf1 := createSimpleTerraformWithDefaultSettings(executor1, projectName)
	tf1.SetVariables(nil)
	err = tf1.DeployFlow(true, false, false)
	assert.NoError(t, err)

	// assert the first run that the files are there
	assertPlanFilesPresence(t, true, true, true)

	// second run, plan has changes, VERY IMPORTANT that we don't find the has-no-changes file here!
	executor2 := &executorMock{isLooseMock: true, planHasChanges: true}
	tf2 := createSimpleTerraformWithDefaultSettings(executor2, projectName)
	tf2.SetVariables(nil)
	err = tf2.DeployFlow(true, false, false)
	assert.NoError(t, err)

	// assert the second run
	assertPlanFilesPresence(t, true, true, false)
}

type executorMock struct {
	mock.Mock
	commands.Executor
	planHasChanges bool
	isLooseMock    bool
}

func (e *executorMock) Execute(command string) (string, error) {
	if !e.isLooseMock {
		e.Called(command)
	}

	if strings.Contains(command, " plan ") {
		if e.planHasChanges {
			return "Terraform will perform the following actions ... To perform exactly these actions, run the following command to apply", exec.Command("bash", "-c", "exit 2").Run()
		} else {
			return "Your infrastructure matches the configuration. ... found no differences, so no changes are needed", nil
		}
	}

	return "success - this output does not matter", nil
}

func (e *executorMock) ExecuteSilent(command string) (string, error) {
	if !e.isLooseMock {
		e.Called(command)
	}

	return "success - this output does not matter", nil
}

func (e *executorMock) AskUserToConfirm(displayMessage string) bool {
	if !e.isLooseMock {
		e.Called(displayMessage)
	}

	return true
}

type variablesStruct struct {
	AString string `mapstructure:"aString" json:"aString" yaml:"aString"`
	AnInt   int    `mapstructure:"anInt" json:"anInt" yaml:"anInt"`
	ABool   bool   `mapstructure:"aBool" json:"aBool" yaml:"aBool"`
}

func createSimpleTerraformWithDefaultSettings(executorMock *executorMock, projectName string) Terraform {
	return New(executorMock, projectName, "1234", "3214",
		"westeurope", "testrg", "storeaccount",
		// important to keep the current directory configured, since some tests rely on this
		// location to verify that expected directories / files are created
		filepath.Join("."),
		DefaultBackendStorageSettings,
		DefaultDeploymentSettings)
}

func deleteDirectoryIfExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %v\n", dirPath)
		return nil
	} else if err != nil {
		return err
	}

	// Directory exists, so attempt to delete it
	err = os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	fmt.Printf("Directory deleted: %v\n", dirPath)
	return nil
}

func assertPlanFilesPresence(t *testing.T, jsonFile bool, txtFile bool, hasNoChangesFile bool) {
	_, err := os.Stat(filepath.Join(".plans", projectName+".deploy.tfplan.json"))
	if jsonFile {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}

	_, err = os.Stat(filepath.Join(".plans", projectName+".deploy.tfplan.txt"))
	if txtFile {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}

	_, err = os.Stat(filepath.Join(".plans", projectName+".deploy.tfplan.plan-has-no-changes"))
	if hasNoChangesFile {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}
}
