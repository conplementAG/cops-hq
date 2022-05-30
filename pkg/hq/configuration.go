package hq

import (
	"fmt"
	"github.com/conplementag/cops-hq/internal"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

func (hq *hqContainer) LoadConfigFile(filePath string) error {
	var err error

	// this should be kept as ExecuteSilent for security reasons, not to leak the whole config file in plaintext
	// to the log file!
	configFile, err := hq.Executor.ExecuteSilent("sops -d " + filePath)

	if err != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("error recieved while reading the config file: %w", err))
	}

	viper.SetConfigType("yaml")
	err = viper.MergeConfig(strings.NewReader(configFile))

	if err != nil {
		return internal.ReturnErrorOrPanic(fmt.Errorf("error recieved while reading the config file: %w", err))
	}

	return nil
}

func (hq *hqContainer) LoadEnvironmentConfigFile() error {
	configFilePath := filepath.Join(ProjectBasePath, "config", viper.GetString("environment-tag")+".yaml")
	return hq.LoadConfigFile(configFilePath)
}
