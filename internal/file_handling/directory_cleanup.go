package file_handling

import (
	"os"
	"path/filepath"
	"strings"
)

// DeleteFilesStartingWith deletes all files in a given directory, ignoring if none exist or the directory itself does not exist
func DeleteFilesStartingWith(filePrefix, directoryPath string) error {
	dirExists, err := directoryExists(directoryPath)

	if !dirExists {
		return nil
	}

	if err != nil {
		return err
	}

	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), filePrefix) {
			err := os.Remove(filepath.Join(directoryPath, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
