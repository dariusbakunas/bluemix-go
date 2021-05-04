package appid

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type AccessTokenConfig struct {
	ExpiresIn int `json:"expires_in,omitempty"`
}

type RefreshTokenConfig struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn int   `json:"expires_in,omitempty"`
}

type AnonymusAccessConfig struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn int   `json:"expires_in,omitempty"`
}
type TokenClaim struct {
	Source           string  `json:"source"`
	SourceClaim      *string `json:"sourceClaim,omitempty"`
	DestinationClaim *string `json:"destinationClaim,omitempty"`
}

type TokenConfig struct {
	Access            *AccessTokenConfig    `json:"access,omitempty"`
	Refresh           *RefreshTokenConfig   `json:"refresh,omitempty"`
	AnonymousAccess   *AnonymusAccessConfig `json:"anonymousAccess,omitempty"`
	IDTokenClaims     []TokenClaim          `json:"idTokenClaims,omitempty"`
	AccessTokenClaims []TokenClaim          `json:"accessTokenClaims,omitempty"`
}

type config struct {
	client *client.Client
}

type Config interface {
	GetTokenConfig(tenantID string) (TokenConfig, error)
	UpdateTokenConfig(tenantID string, config TokenConfig) error
}

func newConfigAPI(c *client.Client) Config {
	return &config{
		client: c,
	}
}

func (c *config) GetTokenConfig(tenantID string) (TokenConfig, error) {
	tokenConfig := TokenConfig{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/tokens", url.QueryEscape(tenantID)), &tokenConfig)
	return tokenConfig, err
}

func (c *config) UpdateTokenConfig(tenantID string, config TokenConfig) error {
	_, err := c.client.Put(fmt.Sprintf("/management/v4/%s/config/tokens", url.QueryEscape(tenantID)), config, nil)
	return err
}
