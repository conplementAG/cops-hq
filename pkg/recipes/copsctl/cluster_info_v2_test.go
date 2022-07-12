package copsctl

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_Subnets_FindForExistingTeam(t *testing.T) {
	// Arrange
	clusterInfo := loadClusterInfoTestFile(t)

	// Act
	subnetBlue, subnetGreen, err := clusterInfo.GetDevOpsTeamSubnets("ateam")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, subnetBlue.Id, "subnet_id_blue")
	assert.Equal(t, subnetGreen.Id, "subnet_id_green")
}

func Test_Subnets_SearchingForNonExistingDevOpsTeam(t *testing.T) {
	// Arrange
	clusterInfo := loadClusterInfoTestFile(t)

	// Act
	subnetBlue, subnetGreen, err := clusterInfo.GetDevOpsTeamSubnets("does-not-exist")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, subnetBlue)
	assert.Nil(t, subnetGreen)
}

func Test_PrivateDnsZones_FindForExistingZone(t *testing.T) {
	// Arrange
	clusterInfo := loadClusterInfoTestFile(t)

	// Act
	zoneGreen, zoneBlue, err := clusterInfo.GetPrivateDnsZones("file.core.windows.net")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, zoneBlue.Name, "file.core.windows.net")
	assert.Equal(t, zoneGreen.Name, "file.core.windows.net")
}

func Test_PrivateDnsZones_SearchingForNonExistingZone(t *testing.T) {
	// Arrange
	clusterInfo := loadClusterInfoTestFile(t)

	// Act
	zoneGreen, zoneBlue, err := clusterInfo.GetPrivateDnsZones("something.net")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, zoneBlue)
	assert.Nil(t, zoneGreen)
}

func Test_ApplicationDnsZone_FindForExistingTeam(t *testing.T) {
	// Arrange
	clusterInfo := loadClusterInfoTestFile(t)

	// Act
	zone, err := clusterInfo.GetApplicationDnsZone("ateam")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, zone.ZoneName, "ateam.local")
}

func Test_ApplicationDnsZone_SearchingForNonExistingDevOpsTeam(t *testing.T) {
	// Arrange
	clusterInfo := loadClusterInfoTestFile(t)

	// Act
	zone, err := clusterInfo.GetApplicationDnsZone("does-not-exist")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, zone)
}

func loadClusterInfoTestFile(t *testing.T) *InfoV2 {
	fileBytes, err := ioutil.ReadFile(filepath.Join(".", "cluster_info_v2_test.json"))
	assert.NoError(t, err)

	var clusterInfo InfoV2
	err = json.Unmarshal(fileBytes, &clusterInfo)
	assert.NoError(t, err)

	return &clusterInfo
}
