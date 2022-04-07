package testing_utils

import (
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func SkipTestIfOnlyShortTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test (only short tests specified via -short)")
	}
}

func SkipTestIfAzureCliMissing(suite *suite.Suite, executor ExecuteCommandAndGetString) {
	out, _ := executor("az")

	// highly implementation dependant, but ok
	if !strings.Contains(out, "Welcome to the cool new Azure CLI!") {
		suite.T().Skip("Test skipped because Azure CLI was not detected")
	}
}

// ExecuteCommandAndGetString this type only exists to remove any circular dependencies between packages, for example
// executor_tests depends on this, and in turn this packages uses the executor to run commands
type ExecuteCommandAndGetString func(command string) (output string, err error)
