package cmdutil

import (
	"errors"
	"strconv"
	"testing"
)

// MockFunction simulates a function that can fail up to a certain number of times before succeeding.
func MockFunction(maxFailures string) func(string) (string, error) {
	attempts := uint64(0)
	maxFailuresInternal, _ := strconv.ParseUint(maxFailures, 10, 64)
	return func(cmd string) (string, error) {
		if attempts < maxFailuresInternal {
			attempts++
			return "", errors.New("error")
		}
		return "success", nil
	}
}

func TestExecuteWithRetryCommand(t *testing.T) {
	tests := []struct {
		name        string
		maxFailures string
		maxAttempts uint
		expectError bool
	}{
		{"SuccessOnFirstTry", "0", 3, false},
		{"SuccessOnRetry", "2", 3, false},
		{"FailAfterMaxAttempts", "3", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := MockFunction(tt.maxFailures)
			err := ExecuteWithRetry(function, "cmd", tt.maxAttempts)
			if tt.expectError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			}
		})
	}
}
