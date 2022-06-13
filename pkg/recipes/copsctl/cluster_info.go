package copsctl

type Info struct {
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
