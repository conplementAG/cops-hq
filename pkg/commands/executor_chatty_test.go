//go:build linux || darwin

// go build directive above because most of the commands below do not work on Windows

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/denisbiondic/cops-hq/internal/testing_utils"
	"github.com/denisbiondic/cops-hq/pkg/logging"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

const testLogFileName = "exec_tests.log"

type ExecutorChattyTestSuite struct {
	suite.Suite
	exec Executor
}

func (s *ExecutorChattyTestSuite) SetupTest() {
	logger := logging.Init(testLogFileName)

	// executor initialized with same file as the logging system to test conflicts when writing to the same file
	s.exec = NewChatty(testLogFileName, logger)
}

func (s *ExecutorChattyTestSuite) AfterTest(suiteName string, testName string) {
	fmt.Printf("Cleanup after test %s/%s\n", suiteName, testName)
	os.Remove(testLogFileName)
}

func TestChattyTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutorChattyTestSuite))
}

func (s *ExecutorChattyTestSuite) Test_ExecuteNormalCommand() {
	s.exec.Execute("ls -la")
}

func (s *ExecutorChattyTestSuite) Test_ReturnsCommandStdoutOutput() {
	out, _ := s.exec.Execute("echo test")
	s.Equal("test", out)
}

func (s *ExecutorChattyTestSuite) Test_NotFoundCommandsFailWithErrorsAndNoOutput() {
	out, err := s.exec.Execute("no-such-thing-to-do bla")
	s.Error(err)
	s.Contains(err.Error(), "executable file not found")
	s.Equal("", out)
}

func (s *ExecutorChattyTestSuite) Test_ExecuteCommandInTTYMode() {
	s.exec.ExecuteTTY("ls -la") // simply run the command, confirm it does not fail
}

func (s *ExecutorChattyTestSuite) Test_ExecuteCommandWithArgumentsWithSpacesAndQuotations() {
	out, _ := s.exec.Execute("echo \"this is a long string\"")
	s.Equal("this is a long string", out)
}

func (s *ExecutorChattyTestSuite) Test_SuccessfulCommandReturnNoErrors() {
	_, err := s.exec.Execute("echo test")
	s.NoError(err)
}

func (s *ExecutorChattyTestSuite) Test_CommandStdErrIsNotCollectedForTheOutput() {
	out, err := s.exec.Execute("ls this-file-does-not-exist")
	s.Error(err)
	s.Contains(err.Error(), "exit status 1")
	s.NotContains(out, "No such file")
}

func (s *ExecutorChattyTestSuite) Test_Integration_ParsingComplexTypeFromCommandsIsPossible() {
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
