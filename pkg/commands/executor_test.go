package commands

import (
	"encoding/json"
	"fmt"
	"github.com/conplementag/cops-hq/internal/testing_utils"
	"github.com/conplementag/cops-hq/pkg/logging"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"strings"
	"testing"
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
	s.exec.Execute("ls -la")
}

func (s *ExecutorTestSuite) Test_ReturnsCommandStdoutOutput() {
	out, _ := s.exec.Execute("echo test")
	s.Equal("test", out)
}

func (s *ExecutorTestSuite) Test_NotFoundCommandsFailWithErrorsAndNoOutput() {
	out, err := s.exec.Execute("no-such-thing-to-do bla")
	s.Error(err)
	s.Contains(err.Error(), "executable file not found")
	s.Equal("", out)
}

func (s *ExecutorTestSuite) Test_ExecuteCommandInTTYMode() {
	s.exec.ExecuteTTY("ls -la") // simply run the command, confirm it does not fail
}

func (s *ExecutorTestSuite) Test_ExecuteCommandWithArgumentsWithSpacesAndQuotations() {
	out, _ := s.exec.Execute("echo \"this is a long string\"")
	s.Equal("this is a long string", out)
}

func (s *ExecutorTestSuite) Test_SuccessfulCommandReturnNoErrors() {
	_, err := s.exec.Execute("echo test")
	s.NoError(err)
}

func (s *ExecutorTestSuite) Test_CommandStdErrIsNotCollectedForTheOutput() {
	out, err := s.exec.Execute("ls this-file-does-not-exist")
	s.Error(err)
	s.NotContains(out, "No such file")
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
	e.Execute("ls -la")
}
