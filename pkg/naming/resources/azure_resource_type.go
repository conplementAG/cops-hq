package resources

type AzureResourceType string

const (
	ResourceGroup              AzureResourceType = "rg"
	StorageAccount             AzureResourceType = "sa"
	SqlServer                  AzureResourceType = "sqls"
	SqlDatabase                AzureResourceType = "sqldb"
	MongoAtlasCluster          AzureResourceType = "mac"
	SqlManagedInstance         AzureResourceType = "sqlmi"
	SqlElasticPool             AzureResourceType = "sqlep"
	KeyVault                   AzureResourceType = "kv"
	KeyVaultWithoutHyphens     AzureResourceType = "v"
	IotHub                     AzureResourceType = "ioth"
	RecoveryServicesVault      AzureResourceType = "rsv"
	AKSCluster                 AzureResourceType = "aks"
	VirtualNetwork             AzureResourceType = "vn"
	VirtualNetworkGateway      AzureResourceType = "vngw"
	RouteTable                 AzureResourceType = "rt"
	ApplicationGateway         AzureResourceType = "agw"
	PublicIp                   AzureResourceType = "pip"
	PrivateEndpoint            AzureResourceType = "pe"
	Bastion                    AzureResourceType = "bn"
	UserAssignedIdentity       AzureResourceType = "uai"
	LogAnalyticsWorkspace      AzureResourceType = "law"
	NetworkSecurityGroup       AzureResourceType = "nsg"
	CosmosDB                   AzureResourceType = "cdb"
	AzureCacheForRedis         AzureResourceType = "acfr"
	LoadBalancer               AzureResourceType = "lb"
	VirtualMachineScalSetLinux AzureResourceType = "vmss"
	EventGridNamespace         AzureResourceType = "egns"
	EventgridTopicSpace        AzureResourceType = "egts"
	StorageBackupVault         AzureResourceType = "sbv"
	DataFactory                AzureResourceType = "df"
	ElasticSearch              AzureResourceType = "es"
	ServicebusNamespace        AzureResourceType = "sbus"
)
