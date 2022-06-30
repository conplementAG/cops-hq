package hq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LogFileNameNotRequiredWhenLoggingToFileDisabled(t *testing.T) {
	options := &HqOptions{DisableFileLogging: true}
	assert.NoError(t, options.Validate())
}

func Test_LogFileNameRequired(t *testing.T) {
	options := &HqOptions{DisableFileLogging: false}
	assert.Error(t, options.Validate())

	optionsWithFileName := &HqOptions{DisableFileLogging: false, LogFileName: "bla.log"}
	assert.NoError(t, optionsWithFileName.Validate())
}
