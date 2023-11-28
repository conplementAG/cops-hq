package copsctl

// ClusterInfoV1 object is a wrapper over the info cluster response in version 1.
type ClusterInfoV1 struct {
	Version       string `json:"version"`
	Description   string `json:"description"`
	Identifier    string `json:"identifier"`
	OidcIssuerUrl string `json:"oidc_issuer_url"`
}
