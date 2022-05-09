package cli

import (
	"github.com/denisbiondic/cops-hq/internal/testing_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FactoryCreatesFunctioningCli(t *testing.T) {
	// Arrange & Act
	cli := New("myprogram", "1.0.0")

	// Assert
	err := cli.Run()
	assert.NotNil(t, cli)
	assert.NoError(t, err)
}

func Test_CreatedCliSupportsVersionFlag(t *testing.T) {
	// Arrange
	cli := New("myprog", "0.0.1")
	outputBuffer := testing_utils.PrepareCommandForTesting(cli.GetRootCommand(), "--version")

	// Act
	cli.Run()

	// Assert
	assert.Contains(t, testing_utils.ReadBuffer(t, outputBuffer), "0.0.1")
}
