package appid

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type IDPS interface {
	// GetCloudDirectoryIDP returns the Cloud Directory identity provider configuration
	GetCloudDirectoryIDP(tenantID string) (CloudDirectoryIDP, error)
	// GetCustomIDP returns custom identity provider configuration
	GetCustomIDP(tenantID string) (CustomIDP, error)
	// GetFacebookIDP returns the Facebook identity provider configuration
	GetFacebookIDP(tenantID string) (GenericIDP, error)
	// GetGoogleIDP returns the Google identity provider configuration
	GetGoogleIDP(tenantID string) (GenericIDP, error)
	// GetSAMLIDP returns the SAML identity provider configuration, including status and credentials.
	GetSAMLIDP(tenantID string) (SAMLIDP, error)
}

type idps struct {
	client *client.Client
}

func newIDPSAPI(c *client.Client) IDPS {
	return &idps{
		client: c,
	}
}

type GenericIDPConfig struct {
	IDPID  string `json:"idpId"`
	Secret string `json:"secret"`
}

type GenericIDP struct {
	IsActive    bool              `json:"isActive"`
	Config      *GenericIDPConfig `json:"config,omitempty"`
	RedirectURL string            `json:"redirectURL"`
}

type GoogleIDP struct {
	IsActive bool `json:"isActive"`
}

type AuthNContext struct {
	Class      []string `json:"class,omitempty"`
	Comparison string   `json:"comparison,omitempty"`
}

type SAMLConfig struct {
	EntityID        string        `json:"entityID"`
	DisplayName     string        `json:"displayName,omitempty"`
	SignInURL       string        `json:"signInUrl"`
	Certificates    []string      `json:"certificates"`
	AuthNContext    *AuthNContext `json:"authnContext,omitempty"`
	SignRequest     *bool         `json:"signRequest,omitempty"`
	EncryptResponse *bool         `json:"encryptResponse,omitempty"`
	IncludeScoping  *bool         `json:"includeScoping,omitempty"`
}

type SAMLIDP struct {
	IsActive bool        `json:"isActive"`
	Config   *SAMLConfig `json:"config,omitempty"`
}

type IdentityConfirmation struct {
	AccessMode string   `json:"accessMode"`
	Methods    []string `json:"methods"`
}

type CloudDirectoryInteractions struct {
	WelcomeEnabled                   bool                 `json:"welcomeEnabled"`
	ResetPasswordEnabled             bool                 `json:"resetPasswordEnabled"`
	ResetPasswordNotificationEnabled bool                 `json:"resetPasswordNotificationEnable"`
	IdentityConfirmation             IdentityConfirmation `json:"identityConfirmation"`
}
type CloudDirectoryConfig struct {
	SelfServiceEnabled bool                       `json:"selfServiceEnabled"`
	SignupEnabled      *bool                      `json:"signupEnabled,omitempty"`
	Interactions       CloudDirectoryInteractions `json:"interactions"`
	IdentityField      string                     `json:"itentityField,omitempty"`
}
type CloudDirectoryIDP struct {
	IsActive bool                  `json:"isActive"`
	Config   *CloudDirectoryConfig `json:"config,omitempty"`
}

type CustomIDPConfig struct {
	PublicKey string `json:"publicKey,omitempty"`
}
type CustomIDP struct {
	IsActive bool             `json:"isActive"`
	Config   *CustomIDPConfig `json:"config,omitempty"`
}

// GetFacebookIDP ...
func (i *idps) GetFacebookIDP(tenantID string) (GenericIDP, error) {
	fb := GenericIDP{}
	_, err := i.client.Get(fmt.Sprintf("/management/v4/%s/config/idps/facebook", url.QueryEscape(tenantID)), &fb)
	return fb, err
}

// GetGoogleIDP ...
func (i *idps) GetGoogleIDP(tenantID string) (GenericIDP, error) {
	google := GenericIDP{}
	_, err := i.client.Get(fmt.Sprintf("/management/v4/%s/config/idps/google", url.QueryEscape(tenantID)), &google)
	return google, err
}

// GetSAMLIDP ...
func (i *idps) GetSAMLIDP(tenantID string) (SAMLIDP, error) {
	saml := SAMLIDP{}
	_, err := i.client.Get(fmt.Sprintf("/management/v4/%s/config/idps/saml", url.QueryEscape(tenantID)), &saml)
	return saml, err
}

// GetCloudDirectoryIDP
func (i *idps) GetCloudDirectoryIDP(tenantID string) (CloudDirectoryIDP, error) {
	cd := CloudDirectoryIDP{}
	_, err := i.client.Get(fmt.Sprintf("/management/v4/%s/config/idps/cloud_directory", url.QueryEscape(tenantID)), &cd)
	return cd, err
}

// GetCustomIDP
func (i *idps) GetCustomIDP(tenantID string) (CustomIDP, error) {
	custom := CustomIDP{}
	_, err := i.client.Get(fmt.Sprintf("/management/v4/%s/config/idps/custom", url.QueryEscape(tenantID)), &custom)
	return custom, err
}
