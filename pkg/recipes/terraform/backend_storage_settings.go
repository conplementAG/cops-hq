package terraform

type BackendStorageSettings struct {
	CreateResourceGroup bool
	BlobContainerName   string
	BlobContainerKey    string
}

var DefaultBackendStorageSettings = BackendStorageSettings{
	CreateResourceGroup: true,
	BlobContainerName:   "tfstate",
	BlobContainerKey:    "terraform.tfstate",
}
