//go:build linux || darwin

// go build directive above because most of the commands below do not work on Windows

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/denisbiondic/cops-hq/internal/testing_utils"
	"github.com/denisbiondic/cops-hq/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

const testLogFileName = "exec_tests.log"

type ExecutorTestSuite struct {
	suite.Suite
	exec *Executor
}

func (s *ExecutorTestSuite) SetupTest() {
	logging.Init(testLogFileName)
	logrus.Info("I am running from the tests...")

	// executor initialized with same file as the logging system to test conflicts when writing to the same file
	s.exec = NewExecutor(testLogFileName)
}

func (s *ExecutorTestSuite) AfterTest(suiteName string, testName string) {
	fmt.Printf("Cleanup after test %s/%s\n", suiteName, testName)
	os.Remove(testLogFileName)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutorTestSuite))
}

func (s *ExecutorTestSuite) Test_ExecuteNormalCommand() {
	s.exec.Execute("ls -la")
}

func (s *ExecutorTestSuite) Test_ReturnsCommandStdoutOutput() {
	out, _ := s.exec.Execute("echo test")
	s.Equal("test", out)
}

func (s *ExecutorTestSuite) Test_FailWithError() {
	_, err := s.exec.Execute("no-such-thing-to-do bla")
	s.Error(err)
	s.Contains(err.Error(), "executable file not found")
}

func (s *ExecutorTestSuite) Test_ExecuteCommandInTTYMode() {
	s.exec.ExecuteTTY("ls -la") // simply run the command, confirm it does not fail
}

func (s *ExecutorTestSuite) Test_ExecuteCommandWithArgumentsWithSpacesAndQuotations() {
	out, _ := s.exec.Execute("echo \"this is a long string\"")
	s.Equal("this is a long string", out)
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
