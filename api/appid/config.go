package appid

import (
	"bytes"
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
	Source           string `json:"source"`
	SourceClaim      string `json:"sourceClaim,omitempty"`
	DestinationClaim string `json:"destinationClaim,omitempty"`
}

type TokenConfig struct {
	Access            *AccessTokenConfig    `json:"access,omitempty"`
	Refresh           *RefreshTokenConfig   `json:"refresh,omitempty"`
	AnonymousAccess   *AnonymusAccessConfig `json:"anonymousAccess,omitempty"`
	IDTokenClaims     []TokenClaim          `json:"idTokenClaims,omitempty"`
	AccessTokenClaims []TokenClaim          `json:"accessTokenClaims,omitempty"`
}

type UsersProfileSettings struct {
	IsActive bool `json:"isActive"`
}

type ThemeText struct {
	Footnote string `json:"footnote"`
	TabTitle string `json:"tabTitle"`
}

type ThemeColors struct {
	HeaderColor string `json:"headerColor"`
}

type config struct {
	client *client.Client
}

type Config interface {
	GetRedirectUris(tenantID string) ([]string, error)
	GetSAMLMetadata(tenantID string) (string, error)
	GetThemeColors(tenantID string) (ThemeColors, error)
	GetThemeText(tenantID string) (ThemeText, error)
	GetTokenConfig(tenantID string) (TokenConfig, error)
	GetUsersProfileSettings(tenantID string) (UsersProfileSettings, error)
	GetWidgetLogoURI(tenantID string) (string, error)
	UpdateTokenConfig(tenantID string, config TokenConfig) error
}

func newConfigAPI(c *client.Client) Config {
	return &config{
		client: c,
	}
}

// GetTokenConfig returns the token configuration
func (c *config) GetTokenConfig(tenantID string) (TokenConfig, error) {
	tokenConfig := TokenConfig{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/tokens", url.QueryEscape(tenantID)), &tokenConfig)
	return tokenConfig, err
}

// GetRedirectUris returns the list of the redirect URIs that can be used as callbacks of App ID authentication flow
func (c *config) GetRedirectUris(tenantID string) ([]string, error) {
	response := struct {
		RedirectUris []string `json:"redirectUris"`
	}{}

	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/redirect_uris", url.QueryEscape(tenantID)), &response)

	if err != nil {
		return nil, err
	}

	return response.RedirectUris, nil
}

// GetUsersProfileSettings returns user profile settings
func (c *config) GetUsersProfileSettings(tenantID string) (UsersProfileSettings, error) {
	settings := UsersProfileSettings{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/users_profile", url.QueryEscape(tenantID)), &settings)
	return settings, err
}

// GetThemeText returns widget texts
func (c *config) GetThemeText(tenantID string) (ThemeText, error) {
	text := ThemeText{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/ui/theme_text", url.QueryEscape(tenantID)), &text)
	return text, err
}

// GetThemeColors returns widget colors
func (c *config) GetThemeColors(tenantID string) (ThemeColors, error) {
	colors := ThemeColors{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/ui/theme_color", url.QueryEscape(tenantID)), &colors)
	return colors, err
}

// GetWidgetLogoURI ...
func (c *config) GetWidgetLogoURI(tenantID string) (string, error) {
	response := struct {
		Image string `json:"image"`
	}{}

	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/ui/media", url.QueryEscape(tenantID)), &response)
	return response.Image, err
}

// GetSAMLMetdata returns the SAML metadata required in order to integrate App ID with a SAML identity provider
func (c *config) GetSAMLMetadata(tenantID string) (string, error) {
	var buf bytes.Buffer
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/saml_metadata", url.QueryEscape(tenantID)), &buf)
	return buf.String(), err
}

// UpdateTokenConfig updates the token configuration
func (c *config) UpdateTokenConfig(tenantID string, config TokenConfig) error {
	_, err := c.client.Put(fmt.Sprintf("/management/v4/%s/config/tokens", url.QueryEscape(tenantID)), config, nil)
	return err
}
