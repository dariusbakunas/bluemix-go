package appid

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Applications interface {
	// Create register a new application with the App ID instance, supported types: regularwebapp, singlepageapp
	Create(tenantID string, name string, applicationType string) (Application, error)
	// Delete deletes an application registered with the App ID instance. Note: This action cannot be undone
	Delete(tenantID string, clientID string) error
	// Get returns a specific application registered with the App ID Instance.
	Get(tenantID string, clientID string) (Application, error)
	// List returns all applications registered with the App ID Instance
	List(tenantID string) ([]Application, error)
	// LitRoles returns defined roles for an application that is registered with an App ID instance
	ListRoles(tenantID string, clientID string) ([]Role, error)
	// ListScopes returns defined scopes for an application that is registered with an App ID instance
	ListScopes(tenantID string, clientID string) ([]string, error)
	// Update update an application registered with the App ID instance, currently only name can be changed
	Update(tenantID string, clientID string, name string) (Application, error)
	// UpdateScopes update the scopes for a registered application. Important: Removing a scope from an array deletes it from any roles that it is associated with and the action cannot be undone
	UpdateScopes(tenantID string, clientID string, scopes []string) ([]string, error)
	// UpdateRoles update the roles for a registered application
	UpdateRoles(tenantID string, clientID string, roleIDs []string) ([]Role, error)
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

// Create ...
func (a *applications) Create(tenantID string, name string, applicationType string) (Application, error) {
	input := struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		Name: name,
		Type: applicationType,
	}

	response := Application{}

	_, err := a.client.Post(fmt.Sprintf("/management/v4/%s/applications", url.QueryEscape(tenantID)), input, &response)
	return response, err
}

// Delete ...
func (a *applications) Delete(tenantID string, clientID string) error {
	_, err := a.client.Delete(fmt.Sprintf("/management/v4/%s/applications/%s", url.QueryEscape(tenantID), url.QueryEscape(clientID)))
	return err
}

// Update ...
func (a *applications) Update(tenantID string, clientID string, name string) (Application, error) {
	input := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	response := Application{}

	_, err := a.client.Put(fmt.Sprintf("/management/v4/%s/applications/%s", url.QueryEscape(tenantID), url.QueryEscape(clientID)), input, &response)
	return response, err
}

// UpdateScopes ...
func (a *applications) UpdateScopes(tenantID string, clientID string, scopes []string) ([]string, error) {
	input := struct {
		Scopes []string `json:"scopes"`
	}{
		Scopes: scopes,
	}

	response := struct {
		Scopes []string `json:"scopes"`
	}{}

	_, err := a.client.Put(fmt.Sprintf("/management/v4/%s/applications/%s/scopes", url.QueryEscape(tenantID), url.QueryEscape(clientID)), input, &response)
	return response.Scopes, err
}

// UpdateRoles ...
func (a *applications) UpdateRoles(tenantID string, clientID string, roleIDs []string) ([]Role, error) {
	input := struct {
		Roles struct {
			IDS []string `json:"ids"`
		} `json:"roles"`
	}{
		Roles: struct {
			IDS []string `json:"ids"`
		}{
			IDS: roleIDs,
		},
	}

	response := struct {
		Roles []Role `json:"roles"`
	}{}

	_, err := a.client.Put(fmt.Sprintf("/management/v4/%s/applications/%s/roles", url.QueryEscape(tenantID), url.QueryEscape(clientID)), input, &response)
	return response.Roles, err
}
