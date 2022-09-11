package hq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHQ_GetRawConfigurationFile_ShouldFailWhenConfigurationNotLoaded(t *testing.T) {
	// Arrange
	hq := New("hq", "0.0.1", "test-logs.txt")

	// Act & Assert
	config, err := hq.GetRawConfigurationFile()

	assert.Equal(t, "", config)
	assert.Equal(t, "configuration was not loaded yet. load configfile first.", err.Error())
}
