package terraform

type BackendStorageSettings struct {
	CreateResourceGroup bool
	// RequireInfrastructureEncryption adds another layer of encryption to the storage account
	RequireInfrastructureEncryption bool
	Tags                            map[string]string
	BlobContainerName               string
	BlobContainerKey                string
	AllowedIpAddresses              []string
}

var DefaultBackendStorageSettings = BackendStorageSettings{
	CreateResourceGroup:             true,
	RequireInfrastructureEncryption: true,
	Tags:                            map[string]string{},
	BlobContainerName:               "tfstate",
	BlobContainerKey:                "terraform.tfstate",
	AllowedIpAddresses:              []string{},
}
