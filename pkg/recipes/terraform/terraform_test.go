package terraform

import (
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/conplementag/cops-hq/pkg/hq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_SettingsAreNotSharedByReferenceBetweenMultipleInstances(t *testing.T) {
	// we want to make sure that the design of non-shared settings is always kept, even after refactorings
	// Arrange
	tf1 := New(&executorMock{}, "app1", "1234", "3214",
		"westeurope", "testrg", "storeaccount", hq.ProjectBasePath,
		DefaultBackendStorageSettings, DefaultDeploymentSettings)

	tf2 := New(&executorMock{}, "app2", "1234", "3214",
		"westeurope", "testrg", "storeaccount", hq.ProjectBasePath,
		DefaultBackendStorageSettings, DefaultDeploymentSettings)

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
	tf := New(&executorMock{}, "test", "1234", "3214",
		"westeurope", "testrg", "storeaccount",
		filepath.Join("."),
		DefaultBackendStorageSettings, DefaultDeploymentSettings)

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

	fileBytes, err := ioutil.ReadFile(filepath.Join(".", tf.GetVariablesFileName()))

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

type executorMock struct {
	mock.Mock
	commands.Executor
}

type variablesStruct struct {
	AString string `mapstructure:"aString" json:"aString" yaml:"aString"`
	AnInt   int    `mapstructure:"anInt" json:"anInt" yaml:"anInt"`
	ABool   bool   `mapstructure:"aBool" json:"aBool" yaml:"aBool"`
}
