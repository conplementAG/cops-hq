package terraform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TagSerialization(t *testing.T) {
	result := serializeTagsIntoAzureCliString(map[string]string{
		"tag1":  "valueA",
		"tag_b": "long string value",
		"tag-c": "123",
	})

	// has to be asserted separately, because iterating go maps does not guarantee order!
	assert.Contains(t, result, "\"tag1\"=\"valueA\"", result)
	assert.Contains(t, result, "\"tag_b\"=\"long string value\"", result)
	assert.Contains(t, result, "\"tag-c\"=\"123\"", result)
}
