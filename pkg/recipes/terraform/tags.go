package terraform

import (
	"fmt"
)

// serializeTagsIntoCmdArgsList converts the given tags from the map to a slice that they can be passed as args to a command
func serializeTagsIntoCmdArgsList(tags map[string]string) []string {
	tagTexts := make([]string, 0)
	for key, value := range tags {
		tagTexts = append(tagTexts, fmt.Sprintf("%s=%s", key, value))
	}

	return tagTexts
}
