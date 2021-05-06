package appid

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type CloudDirectoryUsers interface {
	List(tenantID string, opt *PageOptions) (PagedUsersResponse, error)
}

type PageOptions struct {
	StartIndex int
	Count      int
	Query      string
}

type UserEmail struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

type UserPhone struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type UserName struct {
	Given     string `json:"givenName"`
	Family    string `json:"familyName"`
	Formatted string `json:"formatted"`
}

type UserMeta struct {
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	LastLogin    time.Time `json:"lastLogin"`
	Location     string    `json:"location"`
	ResourceType string    `json:"resourceType"`
}

type User struct {
	ID           string      `json:"id"`
	DisplayName  string      `json:"displayName"`
	Active       bool        `json:"active"`
	Emails       []UserEmail `json:"emails"`
	PhoneNumbers []UserPhone `json:"phoneNumbers,omitempty"`
	Name         UserName    `json:"name"`
	Schemas      []string    `json:"schemas"`
	Meta         UserMeta    `json:"meta"`
	Status       string      `json:"status"`
}

type PagedUsersResponse struct {
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	Schemas      []string `json:"schemas"`
	Resources    []User   `json:"Resources"`
}

type cloudDirectoryUsers struct {
	client *client.Client
}

func newCloudDirectoryUsersAPI(c *client.Client) CloudDirectoryUsers {
	return &cloudDirectoryUsers{
		client: c,
	}
}

// List ...
func (c *cloudDirectoryUsers) List(tenantID string, opt *PageOptions) (PagedUsersResponse, error) {
	query := url.Values{}

	if opt != nil {
		query.Add("startIndex", strconv.Itoa(opt.StartIndex))

		if opt.Count != 0 {
			query.Add("count", strconv.Itoa(opt.Count))
		}

		if opt.Query != "" {
			query.Add("query", opt.Query)
		}
	}

	response := PagedUsersResponse{}
	_, err := c.client.Get(fmt.Sprintf("/management/v4/%s/cloud_directory/Users?%s", url.QueryEscape(tenantID), query.Encode()), &response)
	return response, err
}
