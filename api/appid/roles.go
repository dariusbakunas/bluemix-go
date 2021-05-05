package appid

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Roles interface {
	// ListRoles returns a list of the roles that are associated with your registered application.
	List(tenantID string) ([]Role, error)
	// Get by using the role ID, obtain the information for a specific role that is associated with a registered application.
	Get(tenantID string, roleID string) (Role, error)
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

func (r *roles) Get(tenantID string, roleID string) (Role, error) {
	role := Role{}

	_, err := r.client.Get(fmt.Sprintf("/management/v4/%s/roles/%s", url.QueryEscape(tenantID), url.QueryEscape(roleID)), &role)
	return role, err
}
