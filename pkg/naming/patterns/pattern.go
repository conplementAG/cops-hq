package patterns

type Pattern string

const (
	// Normal pattern creates name such as context-module-name-region-environment-resource, e.g. cops-controller-cache-weu-prod-rg
	// undefined variables, such as {module}, will be skipped if not provided
	Normal Pattern = "{context}{module}{name}{region}{environment}{resource_suffix}"
)
