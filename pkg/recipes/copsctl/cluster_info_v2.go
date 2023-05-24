package copsctl

import (
	"errors"
	"github.com/ahmetb/go-linq"
)

// InfoV2 object is a wrapper over the cluster-info response in version 2.
type InfoV2 struct {
	Version                        string               `json:"version"`
	SubscriptionId                 string               `json:"subscription_id"`
	TenantId                       string               `json:"tenant_id"`
	EgressStaticOutboundIpsEnabled bool                 `json:"egress_static_outbound_ips_enabled"`
	EgressStaticOutboundIps        []string             `json:"egress_static_outbound_ips"`
	TechnicalAccountName           string               `json:"technical_account_name"`
	TechnicalAccountNamespace      string               `json:"technical_account_namespace"`
	CoreOpsMonitoringAndLogging    string               `json:"coreops_monitoring_and_logging"`
	LogAnalyticsWorkspace          string               `json:"log_analytics_workspace"`
	NetworkingBlue                 Networking           `json:"networking_blue"`
	NetworkingGreen                Networking           `json:"networking_green"`
	ApplicationDnsZones            []ApplicationDnsZone `json:"application_dns_zones"`
	OidcIssuerProfileUrl           string               `json:"oidc_issuer_profile_url"`
}

type Networking struct {
	ClusterSubnetId    string           `json:"cluster_subnet_id"`
	ApplicationSubnets []Subnet         `json:"application_subnets"`
	PrivateDnsZones    []PrivateDnsZone `json:"private_dns_zones"`
}

type Subnet struct {
	Name                      string `json:"name"`
	Id                        string `json:"id"`
	Cidr                      string `json:"cidr"`
	EndpointResourceGroupName string `json:"endpoint_resourcegroup_name"`
}

type PrivateDnsZone struct {
	ResourceGroupName string `json:"resourcegroup_name"`
	ZoneName          string `json:"zone_name"`
	Name              string `json:"name"`
}

type ApplicationDnsZone struct {
	Name              string `json:"name"`
	SubscriptionId    string `json:"subscription_id"`
	ResourceGroupName string `json:"resourcegroup_name"`
	ZoneName          string `json:"zone_name"`
}

// GetDevOpsTeamSubnets finds both the blue & green subnets for a given DevOps team name
func (info *InfoV2) GetDevOpsTeamSubnets(devOpsTeamName string) (subnetBlue *Subnet, subnetGreen *Subnet, err error) {
	subnetBlueResult, ok := linq.From(info.NetworkingBlue.ApplicationSubnets).SingleWithT(func(subnet Subnet) bool {
		return subnet.Name == devOpsTeamName
	}).(Subnet)

	subnetBlue = &subnetBlueResult

	if !ok {
		return nil, nil, errors.New("Subnet blue for team " + devOpsTeamName + " not found!")
	}

	subnetGreenResult, ok := linq.From(info.NetworkingGreen.ApplicationSubnets).SingleWithT(func(subnet Subnet) bool {
		return subnet.Name == devOpsTeamName
	}).(Subnet)

	subnetGreen = &subnetGreenResult

	if !ok {
		return nil, nil, errors.New("Subnet green for team " + devOpsTeamName + " not found!")
	}

	return subnetBlue, subnetGreen, nil
}

// GetPrivateDnsZones finds both the blue & green zones for a given name
// (e.g. file.core.windows.net, returns privatelink.file.core.windows.net zones for both blue & green environments)
func (info *InfoV2) GetPrivateDnsZones(name string) (zoneBlue *PrivateDnsZone, zoneGreen *PrivateDnsZone, err error) {
	zoneBlueResult, ok := linq.From(info.NetworkingBlue.PrivateDnsZones).SingleWithT(func(zone PrivateDnsZone) bool {
		return zone.Name == name
	}).(PrivateDnsZone)

	zoneBlue = &zoneBlueResult

	if !ok {
		return nil, nil, errors.New("Private DNS Zone blue for name " + name + " not found!")
	}

	zoneGreenResult, ok := linq.From(info.NetworkingGreen.PrivateDnsZones).SingleWithT(func(zone PrivateDnsZone) bool {
		return zone.Name == name
	}).(PrivateDnsZone)

	zoneGreen = &zoneGreenResult

	if !ok {
		return nil, nil, errors.New("Private DNS Zone green for name " + name + " not found!")
	}

	return zoneBlue, zoneGreen, nil
}

// GetApplicationDnsZone finds the public DNS zone assigned to the DevOps team
func (info *InfoV2) GetApplicationDnsZone(devOpsTeamName string) (*ApplicationDnsZone, error) {
	zone, ok := linq.From(info.ApplicationDnsZones).SingleWithT(func(zone ApplicationDnsZone) bool {
		return zone.Name == devOpsTeamName
	}).(ApplicationDnsZone)

	if !ok {
		return nil, errors.New("Application DNS Zone for team " + devOpsTeamName + " not found!")
	}

	return &zone, nil
}
