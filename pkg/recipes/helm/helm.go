package helm

import (
	"errors"
	"fmt"
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

// Helm is a wrapper around common helm functionality used for automated app deployments to kubernetes.
type Helm interface {

	// SetVariables set the variables for the helm deployment. Variables set will be applied on any subsequent operation. If you do not provide
	// any variables with SetVariables, the default variables in file values.yaml are used.
	//
	// Parameters support are
	//     helmVariables (this is a map of helm variables, defined as string keys and interface values. Nested structures are supported)
	SetVariables(helmVariables map[string]interface{}) error

	// Deploy deploys the helm charts provided in the helmDirectory to the configured kubernetes namespace.
	Deploy() error

	// GetVariablesOverrideFileName returns the file name in which the helm variables will be stored. This name is convention based on the helm tool chain.
	GetVariablesOverrideFileName() string
}

type helmWrapper struct {
	executor commands.Executor

	namespace          string
	chartName          string
	helmDirectory      string
	deploymentSettings DeploymentSettings

	variablesSet bool
}

func (h *helmWrapper) SetVariables(helmVariables map[string]interface{}) error {
	if helmVariables != nil {
		logrus.Info("Setting the helm variables...")

		data, err := yaml.Marshal(&helmVariables)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}

		// file permission: owner: r,w - group: r - other: r -> 0644
		err = ioutil.WriteFile(h.getValuesOverrideFilePath(), data, 0644)

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}

		h.variablesSet = true
	}

	return nil
}

func (h *helmWrapper) Deploy() error {
	if h.deploymentSettings.WaitForJobs && !h.deploymentSettings.Wait {
		return errors.New("[CopsHq][Helm] deployment setting 'WaitForJobs' could not be enabled without enabled 'Wait' flag")
	}

	helmCmd := fmt.Sprintf("helm upgrade --namespace %s --install %s %s -f %s --timeout %s", h.namespace, h.chartName, h.helmDirectory, h.getValuesFilePath(), h.deploymentSettings.Timeout)

	if h.variablesSet {
		helmCmd = fmt.Sprintf("%s -f %s", helmCmd, h.getValuesOverrideFilePath())
	}

	if h.deploymentSettings.Debug {
		helmCmd = fmt.Sprintf("%s --debug", helmCmd)
	}

	if h.deploymentSettings.DryRun {
		helmCmd = fmt.Sprintf("%s --dry-run", helmCmd)
	}

	if h.deploymentSettings.Wait {
		helmCmd = fmt.Sprintf("%s --wait", helmCmd)
	}

	if h.deploymentSettings.WaitForJobs {
		helmCmd = fmt.Sprintf("%s --wait-for-jobs", helmCmd)
	}

	var err error
	if h.deploymentSettings.Wait {
		_, err = h.executor.ExecuteWithProgressInfo(helmCmd)
	} else {
		_, err = h.executor.Execute(helmCmd)
	}

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	return nil
}

func (h *helmWrapper) GetVariablesOverrideFileName() string {
	return "values.override.yaml"
}

func (h *helmWrapper) getValuesFilePath() string {
	if( h.deploymentSettings.OverrideValuePath != "" ) {
		return filepath.Join(h.deploymentSettings.OverrideValuePath, getVariablesFileName())
	} else {
		return filepath.Join(h.helmDirectory, getVariablesFileName())
	}
}

func (h *helmWrapper) getValuesOverrideFilePath() string {
	if( h.deploymentSettings.OverrideValuePath != "" ) {
		return filepath.Join(h.deploymentSettings.OverrideValuePath, h.GetVariablesOverrideFileName())
	} else {
		return filepath.Join(h.helmDirectory, h.GetVariablesOverrideFileName())
	}
}

func getVariablesFileName() string {
	return "values.yaml"
}
