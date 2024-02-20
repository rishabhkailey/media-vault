package v1models

type SpaConfig struct {
	OidcServerUrl               string `json:"oidc_server_url"`
	OidcServerPublicClientId    string `json:"oidc_server_public_client_id"`
	OidcServerDiscoveryEndpoint string `json:"oidc_server_discovery_endpoint"`
}
