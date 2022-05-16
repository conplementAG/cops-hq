package terraform

type DeploymentSettings struct {
	AlwaysCleanLocalCache bool
}

var DefaultDeploymentSettings = DeploymentSettings{
	AlwaysCleanLocalCache: true,
}
