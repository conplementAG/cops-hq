package sops

import (
	"errors"
	"fmt"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Sops is a wrapper around common sops functionality used for config management.
type Sops interface {
	// RegenerateMacValues regenerates the config MAC values for all yaml files in a given directory. You can use this
	// function to fix the sops generated config files with invalid mac values. This usually occurs after merging the
	// configs in Git. This function works on a directory, and it will iterate over all .yaml files it finds.
	//
	// Parameters support are
	//     filePath (path to the config directory)
	RegenerateMacValues(filePath string) error
}

type sopsWrapper struct {
	executor commands.Executor
}

func (s *sopsWrapper) RegenerateMacValues(filePath string) error {
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") && !strings.HasPrefix(info.Name(), ".") {
			logrus.Infof("Checking file: %s\n", path)

			// we check by simply opening the file
			checkMacCommand := fmt.Sprintf("sops -d %s", path)

			if _, err := s.executor.Execute(checkMacCommand); err != nil {
				var exitError *exec.ExitError

				if errors.As(err, &exitError) {
					exitCode := exitError.ExitCode()

					// See https://github.com/mozilla/sops/blob/v3.6.1/cmd/sops/codes/codes.go#L19
					if exitCode == 51 {
						logrus.Infof("Regenerating sops MAC for: %s\n", path)

						vimCommand, err := getVimCommand()
						if err != nil {
							return err
						}

						err = os.Setenv("EDITOR", vimCommand+" -es +\"norm Go\" +\":wq\"")
						if err != nil {
							return err
						}

						_, err = s.executor.Execute(fmt.Sprintf("sops --ignore-mac %s", path))
						if err != nil {
							return err
						}
					} else if exitCode == 128 {
						logrus.Error("CouldNotRetrieveKey - are you logged in to encrypt sops file?")
					}
				}
			} else {
				logrus.Infof("No regeneration needed for %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		logrus.Error("Error:", err)
	}

	return nil
}

func getVimCommand() (string, error) {
	_, err := exec.LookPath("vim")
	if err != nil {
		_, errVi := exec.LookPath("vi")
		if errVi != nil {
			return "", fmt.Errorf("neither vim nor vi could be found - if you are a windows user, please install vim")
		}
		return "vi", nil
	}
	return "vim", nil
}
