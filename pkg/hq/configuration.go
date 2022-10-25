package hq

import (
	"fmt"
	"github.com/conplementag/cops-hq/v2/internal"
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

	hq.RawConfiguration = configFile
	return nil
}

func (hq *hqContainer) LoadEnvironmentConfigFile() error {
	configFilePath := filepath.Join(ProjectBasePath, "config", viper.GetString("environment-tag")+".yaml")
	return hq.LoadConfigFile(configFilePath)
}

func (hq *hqContainer) GetRawConfigurationFile() (string, error) {
	if hq.RawConfiguration == "" {
		return "", internal.ReturnErrorOrPanic(fmt.Errorf("configuration was not loaded yet. load configfile first."))
	}

	return hq.RawConfiguration, nil
}
