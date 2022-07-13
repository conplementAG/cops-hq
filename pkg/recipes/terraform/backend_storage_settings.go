package terraform

type BackendStorageSettings struct {
	CreateResourceGroup bool
	ResourceGroupTags   map[string]string
	BlobContainerName   string
	BlobContainerKey    string
}

var DefaultBackendStorageSettings = BackendStorageSettings{
	CreateResourceGroup: true,
	ResourceGroupTags:   map[string]string{},
	BlobContainerName:   "tfstate",
	BlobContainerKey:    "terraform.tfstate",
}
