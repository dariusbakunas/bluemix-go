package appid

import (
	gohttp "net/http"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

// AppIDServiceAPI is the IBM AppID client
type AppIDServiceAPI interface {
}

type scService struct {
	*client.Client
}

// New construct new AppID client
func New(sess *session.Session) (AppIDServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.AppIDService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	if config.IAMAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefresher, config)
		if err != nil {
			return nil, err
		}
	}

	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.AppIDEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &scService{
		Client: client.New(config, bluemix.AppIDService, tokenRefresher),
	}, nil
}
