package terraform

import (
	"fmt"
	"strings"
)

// serializeTagsIntoAzureCliString converts the given tags from the map to a text that can be passed to an azure cli command
// with the --tags parameter
// The key and value will be quoted e.g. "This is a key"="This is a value" (to support key and value with spaces)
func serializeTagsIntoAzureCliString(tags map[string]string) string {
	var tagTexts []string
	for key, value := range tags {
		tagTexts = append(tagTexts, fmt.Sprintf("%q=%q", key, value))
	}

	return strings.Join(tagTexts, " ")
}
