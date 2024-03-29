package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	authservice "github.com/rishabhkailey/media-vault/internal/services/authService"
	"github.com/rishabhkailey/media-vault/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type OidcClient struct {
	Provider         oidc.Provider
	Verfier          oidc.IDTokenVerifier
	Oauth2Config     oauth2.Config
	ProviderMetadata ProviderMetadata
}

type ProviderMetadata struct {
	Issuer                string   `json:"issuer"`
	AuthURL               string   `json:"authorization_endpoint"`
	TokenURL              string   `json:"token_endpoint"`
	JWKSURL               string   `json:"jwks_uri"`
	UserInfoURL           string   `json:"userinfo_endpoint"`
	IntrospectionEndpoint string   `json:"introspection_endpoint"`
	Algorithms            []string `json:"id_token_signing_alg_values_supported"`
}

func NewOidcClient(issuerUrl, discoveryEndpoint, clientID, clientSecret string) (*OidcClient, error) {
	var err error

	oidcProvider, err := oidc.NewProvider(context.Background(), issuerUrl)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   issuerUrl,
			"error": err,
		}).Error("failed to get oidc provider config")
		return nil, fmt.Errorf("failed to get oidc provider config")
	}
	resp, err := http.Get(discoveryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("http request to %s failed", discoveryEndpoint)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	// todo validate all urls
	var providerMetadata ProviderMetadata
	if err := json.Unmarshal(body, &providerMetadata); err != nil {
		return nil, fmt.Errorf("could not unmarshal json response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	oidcVerifier := oidcProvider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "roles", "email"},
		Endpoint:     oidcProvider.Endpoint(),
	}

	return &OidcClient{
		Provider:         *oidcProvider,
		Verfier:          *oidcVerifier,
		Oauth2Config:     oauth2Config,
		ProviderMetadata: providerMetadata,
	}, nil
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type TokenInfo struct {
	Active      bool                      `json:"active"`
	ClientID    string                    `json:"client_id"`
	Subject     string                    `json:"sub"`
	Scope       string                    `json:"scope"`
	IssuedTime  int64                     `json:"iat"`
	ExpireTime  int64                     `json:"exp"`
	UserName    string                    `json:"username"`
	RealmAccess RealmAccess               `json:"realm_access"`
	Audience    utils.StringOrStringSlice `json:"aud"`
	// TODO: the above RealmAccess claim name is configurable in mappers settings of keycloak
}

func (client *OidcClient) IntrospectToken(token string) (*TokenInfo, error) {
	// todo set basic auth header? not sure if it is required
	data := url.Values{
		"token": []string{token},
	}
	req, err := http.NewRequest(http.MethodPost, client.ProviderMetadata.IntrospectionEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create interospectToken request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	encodedClientIdSecret := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", client.Oauth2Config.ClientID, client.Oauth2Config.ClientSecret),
	))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedClientIdSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("introspect token request failed: %w", err)
	}
	// todo how to differentiate between client auth issue and token issue
	// we get 401 if either token is invalid or client id/secret is invalid
	// maybe we can validate client creds at app start?
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, authservice.ErrUnauthorized
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}
	fmt.Println(string(body))
	var tokenInfo TokenInfo
	if err = json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("could not unmarshal json response: %w", err)
	}
	if !tokenInfo.Active || tokenInfo.ExpireTime < time.Now().Unix() {
		return nil, fmt.Errorf("either token is expired or it is not active: %w", authservice.ErrUnauthorized)
	}
	if !slices.Contains(tokenInfo.Audience.List(), client.Oauth2Config.ClientID) {
		return nil, fmt.Errorf("audience mismatch: %w", authservice.ErrForbidden)
	}
	return &tokenInfo, nil
}
