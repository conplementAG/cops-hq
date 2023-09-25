package terraform

type BackendStorageSettings struct {
	CreateResourceGroup bool
	// RequireInfrastructureEncryption adds another layer of encryption to the storage account
	RequireInfrastructureEncryption bool
	Tags                            map[string]string
	BlobContainerName               string
	BlobContainerKey                string
	// List of IPs or CIDRs to be added to network accesslist. Networking restrictions are applied when first IP or CIDR given
	// Small address ranges using "/31" or "/32" prefix sizes are not supported.
	// These ranges should be configured using individual IP address rules without prefix specified.
	AllowedIpAddresses        []string
	ContainerCreateRetryCount uint
}

var DefaultBackendStorageSettings = BackendStorageSettings{
	CreateResourceGroup:             true,
	RequireInfrastructureEncryption: true,
	Tags:                            map[string]string{},
	BlobContainerName:               "tfstate",
	BlobContainerKey:                "terraform.tfstate",
	AllowedIpAddresses:              []string{},
	ContainerCreateRetryCount:       10,
}
