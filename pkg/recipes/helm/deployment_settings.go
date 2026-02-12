package helm

import "time"

type DeploymentSettings struct {
	// enable verbose output.
	Debug bool
	// simulate an install.
	DryRun bool
	// if set, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment, StatefulSet, or ReplicaSet are in a ready state before marking the release as successful.
	Wait bool
	// if set, will wait until all Jobs have been completed before marking the release as successful.
	WaitForJobs bool
	// time to wait for any individual Kubernetes operation (like Jobs for hooks)
	Timeout time.Duration
	// Override path to values.yaml/values.override.yaml
	OverrideValuePath string
}

var DefaultDeploymentSettings = DeploymentSettings{
	Debug:       false,
	DryRun:      false,
	Wait:        false,
	WaitForJobs: false,
	Timeout:     5 * time.Minute,
	OverrideValuePath: "",
}
