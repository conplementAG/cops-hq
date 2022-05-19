package hq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SetPanicOnAnyError_AlsoSetsTheSameModeForInnerFunctionality(t *testing.T) {
	// Arrange
	hq := New("hq", "0.0.1", "test-logs.txt")

	// Act
	hq.SetPanicOnAnyError(true)

	// Assert
	assert.True(t, hq.GetPanicOnAnyError())
	assert.True(t, hq.GetCli().GetPanicOnAnyError())
	assert.True(t, hq.GetExecutor().GetPanicOnAnyError())
}
