package testing_utils

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func CheckFileContainsString(t *testing.T, fileName string, search string) {
	fileContents, err := ioutil.ReadFile(fileName)

	if err != nil {
		t.Fatal(err)
	}

	fileContentsString := string(fileContents)

	assert.Contains(t, fileContentsString, search)
}
