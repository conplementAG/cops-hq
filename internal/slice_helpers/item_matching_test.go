package slice_helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindItemsToAddAndRemove(t *testing.T) {
	data := []struct {
		testName       string
		actualCurrent  []string
		actualExpected []string
		expectedAdd    []string
		expectedRemove []string
	}{
		// for ease of use letters are used and not ip addresses
		{
			testName:       "green field - nothing to do",
			actualCurrent:  nil,
			actualExpected: nil,
			expectedAdd:    nil,
			expectedRemove: nil,
		},
		{
			testName:       "green field",
			actualCurrent:  nil,
			actualExpected: []string{"A", "B"},
			expectedAdd:    []string{"A", "B"},
			expectedRemove: nil,
		},
		{
			testName:       "brown field - nothing to do",
			actualCurrent:  []string{"A", "B"},
			actualExpected: []string{"A", "B"},
			expectedAdd:    nil,
			expectedRemove: nil,
		},
		{
			testName:       "brown field - nothing to do - order",
			actualCurrent:  []string{"A", "B"},
			actualExpected: []string{"B", "A"},
			expectedAdd:    nil,
			expectedRemove: nil,
		},
		{
			testName:       "brown field - add and remove",
			actualCurrent:  []string{"A", "B", "C"},
			actualExpected: []string{"A", "B", "D"},
			expectedAdd:    []string{"D"},
			expectedRemove: []string{"C"},
		},
	}

	for _, test := range data {
		t.Run(test.testName, func(t *testing.T) {

			// Act
			actualAdd, actualRemove := FindItemsToAddAndRemove(test.actualCurrent, test.actualExpected)

			// Assert
			assert.Equal(t, test.expectedAdd, actualAdd)
			assert.Equal(t, test.expectedRemove, actualRemove)
		})
	}
}
