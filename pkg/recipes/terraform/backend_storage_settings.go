package terraform

type BackendStorageSettings struct {
	CreateResourceGroup bool
	Tags                map[string]string
	BlobContainerName   string
	BlobContainerKey    string
	AllowedIpAddresses  []string
}

var DefaultBackendStorageSettings = BackendStorageSettings{
	CreateResourceGroup: true,
	Tags:                map[string]string{},
	BlobContainerName:   "tfstate",
	BlobContainerKey:    "terraform.tfstate",
	AllowedIpAddresses:  []string{},
}
