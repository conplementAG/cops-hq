package resources

type AzureResourceType string

const (
	ResourceGroup         AzureResourceType = "rg"
	StorageAccount        AzureResourceType = "sa"
	SqlServer             AzureResourceType = "sqls"
	SqlDatabase           AzureResourceType = "sqldb"
	SqlManagedInstance    AzureResourceType = "sqlmi"
	SqlElasticPool        AzureResourceType = "sqlep"
	KeyVault              AzureResourceType = "kv"
	IotHub                AzureResourceType = "ioth"
	RecoveryServicesVault AzureResourceType = "rsv"
	AKSCluster            AzureResourceType = "aks"
	VirtualNetwork        AzureResourceType = "vn"
	VirtualNetworkGateway AzureResourceType = "vngw"
	RouteTable            AzureResourceType = "rt"
	ApplicationGateway    AzureResourceType = "agw"
	PublicIp              AzureResourceType = "pip"
	PrivateEndpoint       AzureResourceType = "pe"
	Bastion               AzureResourceType = "bn"
	UserAssignedIdentity  AzureResourceType = "uai"
	LogAnalyticsWorkspace AzureResourceType = "law"
	NetworkSecurityGroup  AzureResourceType = "nsg"
)
