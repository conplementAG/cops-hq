//go:build linux || darwin

// go build directive above because most of the commands below do not work on Windows

package commands

import (
	"fmt"
	"github.com/denisbiondic/cops-hq/pkg/logging"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ExecutorQuietTestSuite struct {
	suite.Suite
	exec Executor
}

func (s *ExecutorQuietTestSuite) SetupTest() {
	logger := logging.Init(testLogFileName)

	// executor initialized with same file as the logging system to test conflicts when writing to the same file
	s.exec = NewQuiet(testLogFileName, logger)
}

func (s *ExecutorQuietTestSuite) AfterTest(suiteName string, testName string) {
	fmt.Printf("Cleanup after test %s/%s\n", suiteName, testName)
	os.Remove(testLogFileName)
}

func TestQuietTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutorQuietTestSuite))
}

func (s *ExecutorQuietTestSuite) Test_ExecuteNormalCommand() {
	s.exec.Execute("ls -la")
}
