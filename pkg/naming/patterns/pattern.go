package patterns

type Pattern string

const (
	// Normal pattern creates name such as context-module-color-name-region-environment-resource, e.g. cops-controller-green-cache-weu-prod-rg
	// undefined variables, such as {module}, will be skipped if not provided
	Normal Pattern = "{context}{module}{color}{name}{region}{environment}{resource_suffix}"
)
