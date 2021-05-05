package appid

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Roles interface {
	// Create create a new role for a registered application
	Create(tenantID string, input RoleInput) (Role, error)
	// Delete delete an existing role
	Delete(tenantID string, roleID string) error
	// Get by using the role ID, obtain the information for a specific role that is associated with a registered application.
	Get(tenantID string, roleID string) (Role, error)
	// ListRoles returns a list of the roles that are associated with your registered application.
	List(tenantID string) ([]Role, error)
	// Update update an existing role
	Update(tenantID string, roleID string, input RoleInput) (Role, error)
}

type roles struct {
	client *client.Client
}

func newRolesAPI(c *client.Client) Roles {
	return &roles{
		client: c,
	}
}

type RoleAccess struct {
	ApplicationID string   `json:"application_id"`
	Scopes        []string `json:"scopes"`
}

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Access      []RoleAccess `json:"access,omitempty"`
}

type RoleInput struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Access      []RoleAccess `json:"access"`
}

// ListRoles ...
func (r *roles) List(tenantID string) ([]Role, error) {
	response := struct {
		Roles []Role `json:"roles"`
	}{}

	_, err := r.client.Get(fmt.Sprintf("/management/v4/%s/roles", url.QueryEscape(tenantID)), &response)

	if err != nil {
		return nil, err
	}

	return response.Roles, nil
}

// Get ...
func (r *roles) Get(tenantID string, roleID string) (Role, error) {
	role := Role{}

	_, err := r.client.Get(fmt.Sprintf("/management/v4/%s/roles/%s", url.QueryEscape(tenantID), url.QueryEscape(roleID)), &role)
	return role, err
}

// Create
func (r *roles) Create(tenantID string, input RoleInput) (Role, error) {
	response := Role{}

	_, err := r.client.Post(fmt.Sprintf("/management/v4/%s/roles", url.QueryEscape(tenantID)), input, &response)
	return response, err
}

// Delete ...
func (r *roles) Delete(tenantID string, roleID string) error {
	_, err := r.client.Delete(fmt.Sprintf("/management/v4/%s/roles/%s", url.QueryEscape(tenantID), url.QueryEscape(roleID)))
	return err
}

// Update
func (r *roles) Update(tenantID string, roleID string, input RoleInput) (Role, error) {
	response := Role{}

	_, err := r.client.Put(fmt.Sprintf("/management/v4/%s/roles/%s", url.QueryEscape(tenantID), url.QueryEscape(roleID)), input, &response)
	return response, err
}
