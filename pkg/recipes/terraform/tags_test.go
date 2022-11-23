package terraform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TagSerialization(t *testing.T) {
	result := serializeTagsIntoCmdArgsList(map[string]string{
		"tag1":  "valueA",
		"tag_b": "long string value",
		"tag c": "123",
	})

	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
	// has to be asserted separately, because iterating go maps does not guarantee order!
	assert.Contains(t, result, "tag1=valueA", result)
	assert.Contains(t, result, "tag_b=long string value", result)
	assert.Contains(t, result, "tag c=123", result)
}

func Test_TagSerialization_Nil(t *testing.T) {
	result := serializeTagsIntoCmdArgsList(nil)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func Test_TagSerialization_Empty(t *testing.T) {
	result := serializeTagsIntoCmdArgsList(make(map[string]string))

	assert.NotNil(t, result)
	assert.Empty(t, result)
}
