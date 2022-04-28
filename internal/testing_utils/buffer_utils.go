package testing_utils

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func ReadBuffer(t *testing.T, outputBuffer *bytes.Buffer) string {
	out, err := ioutil.ReadAll(outputBuffer)

	if err != nil {
		t.Fatal(err)
	}

	return string(out)
}
