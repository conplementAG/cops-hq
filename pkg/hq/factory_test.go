package hq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewWillCreateFunctioningHQInstance(t *testing.T) {
	hq := New("hq", "0.0.1", "test-logs.txt")
	assert.NotNil(t, hq.GetCli())
	assert.NotNil(t, hq.GetExecutor())
}
