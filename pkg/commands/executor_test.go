package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/conplementag/cops-hq/v2/internal/testing_utils"
	"github.com/conplementag/cops-hq/v2/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const testLogFileName = "exec_tests.log"

type ExecutorTestSuite struct {
	suite.Suite
	exec Executor
}

func (s *ExecutorTestSuite) SetupTest() {
	logger := logging.Init(testLogFileName)

	// executor initialized with same file as the logging system to test conflicts when writing to the same file
	s.exec = NewChatty(testLogFileName, logger)
}

func (s *ExecutorTestSuite) AfterTest(suiteName string, testName string) {
	fmt.Printf("Cleanup after test %s/%s\n", suiteName, testName)
	os.Remove(testLogFileName)
}

func TestChattyTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutorTestSuite))
}

func (s *ExecutorTestSuite) Test_ExecuteNormalCommand() {
	// Use cross-platform command
	if runtime.GOOS == "windows" {
		s.exec.Execute("dir")
	} else {
		s.exec.Execute("ls -la")
	}
}

func (s *ExecutorTestSuite) Test_ReturnsCommandStdoutOutput() {
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "cmd /c echo test"
	} else {
		cmd = "echo test"
	}
	out, _ := s.exec.Execute(cmd)
	s.Equal("test", strings.TrimSpace(out))
}

func (s *ExecutorTestSuite) Test_OsExecCommandsAreCorrectlyExecuted() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo test")
	} else {
		cmd = exec.Command("echo", "test")
	}
	out, _ := s.exec.ExecuteCmd(cmd)
	s.Equal("test", strings.TrimSpace(out))
}

func (s *ExecutorTestSuite) Test_NotFoundCommandsFailWithErrorsAndNoOutput() {
	out, err := s.exec.Execute("no-such-thing-to-do bla")
	s.Error(err)
	s.Contains(err.Error(), "executable file not found")
	s.Equal("", out)
}

func (s *ExecutorTestSuite) Test_ExecuteCommandInTTYMode() {
	// Use cross-platform command - simply run the command, confirm it does not fail
	if runtime.GOOS == "windows" {
		s.exec.ExecuteTTY("dir")
	} else {
		s.exec.ExecuteTTY("ls -la")
	}
}

// func (s *ExecutorTestSuite) Test_ExecuteCommandWithArgumentsWithSpacesAndQuotations() {
// 	out, err := s.exec.Execute("echo \"this is a long string\"")
// 	s.Nil(err)
// 	s.Equal("this is a long string", out)
// }

func (s *ExecutorTestSuite) Test_SuccessfulCommandReturnNoErrors() {
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "cmd /c echo test"
	} else {
		cmd = "echo test"
	}
	_, err := s.exec.Execute(cmd)
	s.NoError(err)
}

// func (s *ExecutorTestSuite) Test_CommandStdErrIsNotCollectedForTheOutput() {
// 	out, err := s.exec.Execute("ls this-file-does-not-exist")
// 	s.Error(err)
// 	s.Contains(err.Error(), "No such file")
// 	s.NotContains(out, "No such file")
// }

func (s *ExecutorTestSuite) Test_CollectsBothStdErrAndStdOutOnError() {
	var cmd string
	if runtime.GOOS == "windows" {
		// Windows: Use PowerShell or cmd to produce stdout, stderr, and error
		cmd = "cmd /c \"echo This is standard output && echo This is standard error 1>&2 && dir this-file-does-not-exist\""
	} else {
		cmd = "bash -c \"{ echo 'This is standard output'; echo 'This is standard error' >&2; ls this-file-does-not-exist; }\""
	}
	_, err := s.exec.Execute(cmd)
	s.Error(err)
	s.Contains(err.Error(), "This is standard output")
	s.Contains(err.Error(), "This is standard error")
}

func (s *ExecutorTestSuite) Test_ErrorHasExitErrorWrapped() {
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "cmd /c \"exit 5\""
	} else {
		cmd = "bash -c \"exit 5\""
	}
	_, err := s.exec.Execute(cmd)

	var exitErr *exec.ExitError
	errors.As(err, &exitErr)

	assert.Equal(s.T(), 5, exitErr.ExitCode())
}

func (s *ExecutorTestSuite) Test_Integration_ParsingComplexTypeFromCommandsIsPossible() {
	// the two methods here can be further optimized in the future if we have more integrations tests, for example
	// by having a list of conditions passed to a single CheckIntegrationTestPrerequisites method?
	testing_utils.SkipTestIfOnlyShortTests(s.T())
	testing_utils.SkipTestIfAzureCliMissing(&s.Suite, s.exec.Execute)

	// Act
	out, _ := s.exec.Execute("az version -o json")

	// Assert
	var resultingMap map[string]interface{}
	json.Unmarshal([]byte(out), &resultingMap)

	_, ok := resultingMap["azure-cli"] // this key is always expected
	s.True(ok)
}

func (s *ExecutorTestSuite) Test_AskUserToConfirm() {
	tests := []struct {
		testName       string
		userInput      string
		expectedResult bool
	}{
		{"false for random input", "$blabla", false},
		{"true for confirmation with yes", "yes", true},
		{"true for confirmation with Y", "Y", true},
		{"true for confirmation with YES", "YES", true},
		{"false for confirmation with no", "no", false},
		{"false for confirmation with newline", "\n", false},
		{"false for confirmation with no input", "", false},
	}

	for _, tt := range tests {
		fmt.Println("Running test: " + tt.testName)
		var reader io.Reader = strings.NewReader(tt.userInput)
		s.exec.(*executor).OverrideStdIn(reader)

		s.Equal(s.exec.AskUserToConfirm("Should I?"), tt.expectedResult)
	}
}

func (s *ExecutorTestSuite) Test_AskUserToConfirmWithKeyword() {
	tests := []struct {
		testName       string
		userInput      string
		keyword        string
		expectedResult bool
	}{
		{"false for random input", "$blabla", "test", false},
		{"true for correct input", "core", "core", true},
	}

	for _, tt := range tests {
		fmt.Println("Running test: " + tt.testName)
		var reader io.Reader = strings.NewReader(tt.userInput)
		s.exec.(*executor).OverrideStdIn(reader)

		s.Equal(s.exec.AskUserToConfirmWithKeyword("Should I?", tt.keyword), tt.expectedResult)
	}
}

func Test_QuietExecutorWorksAsWell(t *testing.T) {
	logger := logging.Init(testLogFileName)
	e := NewQuiet(testLogFileName, logger)
	if runtime.GOOS == "windows" {
		e.Execute("dir")
	} else {
		e.Execute("ls -la")
	}
}
