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

type EmailTemplate struct {
	Subject       string `json:"subject"`
	HTMLBody      string `json:"html_body,omitempty"`
	B64Body       string `json:"base64_encoded_html_body,omitempty"`
	PlainTextBody string `json:"plain_text_body,omitempty"`
}

type SenderNameEmail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SenderDetails struct {
	From           SenderNameEmail `json:"from"`
	ReplyTo        SenderNameEmail `json:"reply_to"`
	CC             SenderNameEmail `json:"cc"`
	BCC            SenderNameEmail `json:"bcc"`
	LinkExpiration *int            `json:"linkExpirationSec,omitempty"`
}

type PasswordRegex struct {
	Regex        string `json:"regex"`
	B64Regex     string `json:"base64_encoded_regex"`
	ErrorMessage string `json:"error_message"`
}

type SendgridMailerSettings struct {
	ApiKey string `json:"apiKey"`
}

type CustomMailerAuth struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomMailerSettings struct {
	URL  string           `json:"url"`
	Auth CustomMailerAuth `json:"authorization"`
}

type EmailDispatcherSettings struct {
	Provider string                  `json:"provider"`
	SendGrid *SendgridMailerSettings `json:"sendgrid,omitempty"`
	Custom   *CustomMailerSettings   `json:"custom,omitempty"`
}

type APMSettings struct {
	Enabled       bool `json:"enabled"`
	PasswordReuse struct {
		Enabled bool `json:"enabled"`
		Config  struct {
			MaxReuse int `json:"maxPasswordReuse"`
		} `json:"config,omitempty"`
	} `json:"passwordReuse"`
	PreventPasswordWithUsername struct {
		Enabled bool `json:"enabled"`
	} `json:"preventPasswordWithUsername"`
	PasswordExpiration struct {
		Enabled bool `json:"enabled"`
		Config  struct {
			DaysToExpire int `json:"daysToExpire"`
		} `json:"config,omitempty"`
	} `json:"passwordExpiration"`
	LockoutPolicy struct {
		Enabled bool `json:"enabled"`
		Config  struct {
			LockoutTime   int `json:"lockOutTimeSec"`
			NumOfAttempts int `json:"numOfAttempts"`
		} `json:"config,omitempty"`
	} `json:"lockOutPolicy"`
	MinPasswordChangeInterval struct {
		Enabled bool `json:"enabled"`
		Config  struct {
			MinHours int `json:"minHoursToChangePassword"`
		} `json:"config,omitempty"`
	} `json:"minPasswordChangeInterval,omitempty"`
}

type config struct {
	client *client.Client
}

type Config interface {
	// GetActionURL returns the custom url to redirect to when action is executed, supported actions: on_user_verified, on_reset_password
	GetActionURL(tenantID string, action string) (string, error)
	// GetAPMSettings returns the configuration of the advanced password management
	GetAPMSettings(tenantID string) (APMSettings, error)
	// GetEmailDispatcherSettings returns configuration of email dispatcher that is used by Cloud Directory when sending emails
	GetEmailDispatcherSettings(tenantID string) (EmailDispatcherSettings, error)
	// GetEmailSenderDetails returns the sender details configuration that is used by Cloud Directory when sending emails
	GetEmailSenderDetails(tenantID string) (SenderDetails, error)
	// GetEmailTemplate returns the content of a custom email template or the default template in case it wasn't customized
	GetEmailTemplate(tenantID string, templateName string, language string) (EmailTemplate, error)
	// GetPasswordRegex returns the regular expression used by App ID for password strength validation
	GetPasswordRegex(tenantID string) (PasswordRegex, error)
	// GetRedirectUris returns the list of the redirect URIs that can be used as callbacks of App ID authentication flow
	GetRedirectUris(tenantID string) ([]string, error)
	// GetSAMLMetdata returns the SAML metadata required in order to integrate App ID with a SAML identity provider
	GetSAMLMetadata(tenantID string) (string, error)
	// GetTemplateLanguages returns the list of languages that can be used to customize email templates for Cloud Directory
	GetTemplateLanguages(tenantID string) ([]string, error)
	// GetThemeColors returns widget colors
	GetThemeColors(tenantID string) (ThemeColors, error)
	// GetThemeText returns widget texts
	GetThemeText(tenantID string) (ThemeText, error)
	// GetTokenConfig returns the token configuration
	GetTokenConfig(tenantID string) (TokenConfig, error)
	// GetUsersProfileSettings returns user profile settings
	GetUsersProfileSettings(tenantID string) (UsersProfileSettings, error)
	// GetWidgetLogoURI returns the link to the custom logo image of the login widget
	GetWidgetLogoURI(tenantID string) (string, error)
	// UpdateTokenConfig updates the token configuration
	UpdateTokenConfig(tenantID string, config TokenConfig) error
}

func newConfigAPI(c *client.Client) Config {
	return &config{
		client: c,
	}
}

// GetTokenConfig ...
func (c *config) GetTokenConfig(tenantID string) (TokenConfig, error) {
	tokenConfig := TokenConfig{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/tokens", url.QueryEscape(tenantID)), &tokenConfig)
	return tokenConfig, err
}

// GetRedirectUris ...
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

// GetUsersProfileSettings ...
func (c *config) GetUsersProfileSettings(tenantID string) (UsersProfileSettings, error) {
	settings := UsersProfileSettings{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/users_profile", url.QueryEscape(tenantID)), &settings)
	return settings, err
}

// GetThemeText ...
func (c *config) GetThemeText(tenantID string) (ThemeText, error) {
	text := ThemeText{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/ui/theme_text", url.QueryEscape(tenantID)), &text)
	return text, err
}

// GetThemeColors ...
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

// GetSAMLMetdata ...
func (c *config) GetSAMLMetadata(tenantID string) (string, error) {
	var buf bytes.Buffer
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/saml_metadata", url.QueryEscape(tenantID)), &buf)
	return buf.String(), err
}

// GetEmailTemplate ...
func (c *config) GetEmailTemplate(tenantID string, templateName string, language string) (EmailTemplate, error) {
	template := EmailTemplate{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/cloud_directory/templates/%s/%s", url.QueryEscape(tenantID), url.QueryEscape(templateName), url.QueryEscape(language)), &template)
	return template, err
}

// GetTemplateLanguages ...
func (c *config) GetTemplateLanguages(tenantID string) ([]string, error) {
	response := struct {
		Languages []string `json:"languages"`
	}{}

	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/ui/languages", url.QueryEscape(tenantID)), &response)

	if err != nil {
		return nil, err
	}

	return response.Languages, nil
}

// GetEmailSenderDetails ...
func (c *config) GetEmailSenderDetails(tenantID string) (SenderDetails, error) {
	details := struct {
		SenderDetails SenderDetails `json:"senderDetails"`
	}{}

	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/cloud_directory/sender_details", url.QueryEscape(tenantID)), &details)
	return details.SenderDetails, err
}

// GetActionURL ...
func (c *config) GetActionURL(tenantID string, action string) (string, error) {
	response := struct {
		URL string `json:"actionUrl"`
	}{}

	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/cloud_directory/action_url/%s", url.QueryEscape(tenantID), url.QueryEscape(action)), &response)
	return response.URL, err
}

// GetPasswordRegex ...
func (c *config) GetPasswordRegex(tenantID string) (PasswordRegex, error) {
	regex := PasswordRegex{}

	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/cloud_directory/password_regex", url.QueryEscape(tenantID)), &regex)
	return regex, err
}

// GetEmailDispatcherSettings ...
func (c *config) GetEmailDispatcherSettings(tenantID string) (EmailDispatcherSettings, error) {
	settings := EmailDispatcherSettings{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/cloud_directory/email_dispatcher", url.QueryEscape(tenantID)), &settings)
	return settings, err
}

// GetAPMSettings ...
func (c *config) GetAPMSettings(tenantID string) (APMSettings, error) {
	settings := struct {
		APM APMSettings `json:"advancedPasswordManagement"`
	}{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/config/cloud_directory/advanced_password_management", url.QueryEscape(tenantID)), &settings)
	return settings.APM, err
}

// UpdateTokenConfig ...
func (c *config) UpdateTokenConfig(tenantID string, config TokenConfig) error {
	_, err := c.client.Put(fmt.Sprintf("/management/v4/%s/config/tokens", url.QueryEscape(tenantID)), config, nil)
	return err
}
