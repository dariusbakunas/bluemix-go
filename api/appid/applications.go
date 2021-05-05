package appid

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Applications interface {
	// Get returns a specific application registered with the App ID Instance.
	Get(tenantID string, clientID string) (Application, error)
	// List returns all applications registered with the App ID Instance
	List(tenantID string) ([]Application, error)
	// LitRoles returns defined roles for an application that is registered with an App ID instance
	ListRoles(tenantID string, clientID string) ([]Role, error)
	// ListScopes returns defined scopes for an application that is registered with an App ID instance
	ListScopes(tenantID string, clientID string) ([]string, error)
}

type applications struct {
	client *client.Client
}

func newApplicationsAPI(c *client.Client) Applications {
	return &applications{
		client: c,
	}
}

type Application struct {
	ClientID          string `json:"clientId"`
	TenantID          string `json:"tenantId"`
	Secret            string `json:"secret,omitempty"`
	Name              string `json:"name"`
	OAuthServerURL    string `json:"oAuthServerUrl"`
	ProfilesURL       string `json:"profilesURL"`
	DiscoveryEndpoint string `json:"discoveryEndpoint"`
	Type              string `json:"type"`
}

// List ...
func (a *applications) List(tenantID string) ([]Application, error) {
	response := struct {
		Applications []Application `json:"applications"`
	}{}

	_, err := a.client.Get(fmt.Sprintf("/management/v4/%s/applications", url.QueryEscape(tenantID)), &response)
	return response.Applications, err
}

// Get ...
func (a *applications) Get(tenantID string, clientID string) (Application, error) {
	app := Application{}
	_, err := a.client.Get(fmt.Sprintf("/management/v4/%s/applications/%s", url.QueryEscape(tenantID), url.QueryEscape(clientID)), &app)
	return app, err
}

// ListScopes ...
func (a *applications) ListScopes(tenantID string, clientID string) ([]string, error) {
	response := struct {
		Scopes []string `json:"scopes"`
	}{}

	_, err := a.client.Get(fmt.Sprintf("/management/v4/%s/applications/%s/scopes", url.QueryEscape(tenantID), url.QueryEscape(clientID)), &response)
	return response.Scopes, err
}

// ListRoles ...
func (a *applications) ListRoles(tenantID string, clientID string) ([]Role, error) {
	response := struct {
		Roles []Role `json:"roles"`
	}{}

	_, err := a.client.Get(fmt.Sprintf("/management/v4/%s/applications/%s/roles", url.QueryEscape(tenantID), url.QueryEscape(clientID)), &response)
	return response.Roles, err
}
