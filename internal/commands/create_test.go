package commands

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CommandParsing(t *testing.T) {
	tests := []struct {
		testName       string
		command        string
		expectedResult []string
	}{
		{"normal command", "run some stuff", []string{"run", "some", "stuff"}},
		{"arguments with spaces", "run -a \"some stuff\"", []string{"run", "-a", "some stuff"}},
		{"multiple arguments with spaces", "run -a \"some stuff\" -b \"other stuff\"", []string{"run", "-a", "some stuff", "-b", "other stuff"}},
		{"long command", "az role assignment create --role \"Network Contributer\" --assignee ABC --scope abc",
			[]string{"az", "role", "assignment", "create", "--role", "Network Contributer", "--assignee", "ABC", "--scope", "abc"}},
		{"arguments with escaped backslash", "run -a \"bla\\bla\"", []string{"run", "-a", "bla\\bla"}},
		{"arguments with quotations", "run -a object[\"property\"]", []string{"run", "-a", "object[\"property\"]"}},
	}

	for _, tt := range tests {
		fmt.Println("Running test: " + tt.testName)
		resultCmd := Create(tt.command)

		assert.Equal(t, tt.expectedResult, resultCmd.Args)
	}
}
