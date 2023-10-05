package file_handling

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_DeleteFilesStartingWith(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-dir")
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create some files in the temporary directory

	err = os.WriteFile(filepath.Join(tempDir, "prefix_file1.txt"), []byte("content"), 0644)
	assert.Nil(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "prefix_file2.txt"), []byte("content"), 0644)
	assert.Nil(t, err)

	// prefix added in the middle to make sure it is not simple string matching done
	err = os.WriteFile(filepath.Join(tempDir, "file1_prefix.txt"), []byte("content"), 0644)
	assert.Nil(t, err)

	// Delete files starting with "prefix_"
	err = DeleteFilesStartingWith("prefix_", tempDir)
	assert.Nil(t, err)

	// Check if the files were deleted
	files, err := os.ReadDir(tempDir)
	assert.Nil(t, err)

	_, err = os.Stat(filepath.Join(tempDir, "file1_prefix.txt"))
	assert.NoError(t, err)

	for _, file := range files {
		assert.False(t, strings.HasPrefix(file.Name(), "prefix_"))
	}
}
