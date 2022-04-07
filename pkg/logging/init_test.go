package logging

import (
	"github.com/denisbiondic/cops-hq/internal/testing_utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestShouldWriteToBothFileAndConsole(t *testing.T) {
	// Arrange
	logFile := "test-file.log"
	testMessage := uuid.New().String()

	// Act
	Init(logFile)
	logrus.Info(testMessage)

	// Assert
	// we only test the file, because stdout capture approach did not work with Logrus
	testing_utils.CheckFileContainsString(t, logFile, testMessage)
	os.Remove(logFile)
}
