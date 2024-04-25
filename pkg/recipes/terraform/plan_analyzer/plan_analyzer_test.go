package plan_analyzer

import (
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform/file_paths"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestPlanAnalizer(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		setupFunc   func() error
		cleanupFunc func() error
		want        bool
	}{
		{
			name:        "Plan is dirty",
			setupFunc:   func() error { return nil },
			cleanupFunc: func() error { return nil },
			want:        true,
		},
		{
			name: "Plan is not dirty",
			setupFunc: func() error {
				// Create temporary directory and file
				err := os.MkdirAll(file_paths.PlansDirectory, 0755)
				if err != nil {
					return err
				}

				fileName := filepath.Join(file_paths.PlansDirectory, file_paths.GetPlanStateFileName("project", false))
				err = os.WriteFile(fileName, []byte("content"), 0644)
				if err != nil {
					return err
				}

				return nil
			},
			cleanupFunc: func() error {
				return os.RemoveAll(file_paths.PlansDirectory)
			},
			want: false,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.setupFunc()
			defer tc.cleanupFunc()
			assert.NoError(t, err)

			// Act
			actual, err := New("project", ".").IsDeployPlanDirty()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want, actual)
		})
	}
}
