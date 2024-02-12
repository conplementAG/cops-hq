package helm

import "github.com/conplementag/cops-hq/v2/pkg/commands"

// New creates a new instance of Helm, which is a wrapper around common Helm functionality. The deployment will run with default DeploymentSettings
//
// Parameters:
//
//	executor (can be provided from hq.GetExecutor() or by instantiating your own),
//	namespace (the kubernetes namespace to deploy to),
//	chartName (the chart name of the helm deployment),
//	helmDirectory (directory where your helm resources are stored. To construct the full path, simply use
//	    filepath.Join() method and the hq.ProjectBasePath)
func New(executor commands.Executor, namespace string, chartName string, helmDirectory string) Helm {
	return NewWithSettings(executor, namespace, chartName, helmDirectory, DefaultDeploymentSettings)
}

// NewWithSettings creates a new instance of Helm, which is a wrapper around common Helm functionality.
//
// Parameters:
//
//	executor (can be provided from hq.GetExecutor() or by instantiating your own),
//	namespace (the kubernetes namespace to deploy to),
//	chartName (the chart name of the helm deployment),
//	helmDirectory (directory where your helm resources are stored. To construct the full path, simply use
//	   filepath.Join() method and the hq.ProjectBasePath)
//	deploymentSettings (deployment specific settings, e.g. to wait for completion of the deployment)
func NewWithSettings(executor commands.Executor, namespace string, chartName string, helmDirectory string, deploymentSettings DeploymentSettings) Helm {

	return &helmWrapper{
		executor:           executor,
		namespace:          namespace,
		chartName:          chartName,
		helmDirectory:      helmDirectory,
		deploymentSettings: deploymentSettings,
	}
}
