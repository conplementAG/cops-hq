package hq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InitWillCreateFunctioningHQInstance(t *testing.T) {
	hq := Init("hq", "0.0.1", "test-logs.txt")
	assert.NotNil(t, hq.Cli)
	assert.NotNil(t, hq.Executor)
}
