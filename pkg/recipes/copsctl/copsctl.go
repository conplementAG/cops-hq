package copsctl

import (
	"encoding/json"
	"github.com/conplementag/cops-hq/internal"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/sirupsen/logrus"
)

type Copsctl interface {
	// Connect sets the local kubectl connection to the configured cluster. Parameters:
	//     clusterName - cluster environment tag for the cluster you want to connect to
	//     clusterConnectionString - connection string of the cluster (either a user or technical connection string)
	//     isTechnicalAccountConnect - set to true, if the clusterConnectionString used is technical account connection string
	//     connectToSecondaryCluster - if set, will connect to secondary cluster (use for cluster migrations)
	Connect(clusterName string, clusterConnectionString string, isTechnicalAccountConnect bool, connectToSecondaryCluster bool) error

	// GetClusterInfo returns the cluster info for the currently connected cluster
	GetClusterInfo() (*Info, error)
}

type copsctl struct {
	executor commands.Executor
}

func (c *copsctl) Connect(clusterName string, clusterConnectionString string, isTechnicalAccountConnect bool, connectToSecondaryCluster bool) error {
	logrus.Info("[Cluster] Connecting to cluster " + clusterName + " ...")

	copsConnectCmd := "copsctl connect -e " + clusterName + " -c \"" + clusterConnectionString + "\" -a"

	if isTechnicalAccountConnect {
		copsConnectCmd = copsConnectCmd + " -t"
	}

	if connectToSecondaryCluster {
		copsConnectCmd = copsConnectCmd + " -s"
	}

	_, err := c.executor.Execute(copsConnectCmd)

	if err != nil {
		return internal.ReturnErrorOrPanic(err)
	}

	// workaround to force kubectl login for interactive-login mode
	if !isTechnicalAccountConnect {
		err = c.executor.ExecuteTTY("kubectl auth can-i list copsnamespaces.coreops.conplement.cloud") // this query should always work in copsctl context

		if err != nil {
			return internal.ReturnErrorOrPanic(err)
		}
	}

	logrus.Info("Done.")
	return nil
}

func (c *copsctl) GetClusterInfo() (*Info, error) {
	logrus.Info("Receiving cluster info...")

	clusterInfoJson, err := c.executor.ExecuteSilent("copsctl cluster-info --print-to-stdout-silence-everything-else")

	if err != nil {
		return nil, internal.ReturnErrorOrPanic(err)
	}

	var clusterInfo Info
	err = json.Unmarshal([]byte(clusterInfoJson), &clusterInfo)

	if err != nil {
		return nil, internal.ReturnErrorOrPanic(err)
	}

	logrus.Info("Done.")

	return &clusterInfo, nil
}
