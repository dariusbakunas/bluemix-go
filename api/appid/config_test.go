package appid

import (
	"log"
	"net/http"

	"github.com/onsi/gomega/ghttp"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	tokenConfig = `{
		"idTokenClaims": [
			{
				"source": "attributes",
				"sourceClaim": "theme"
			}
		],
		"accessTokenClaims": [
			{
				"source": "saml",
				"sourceClaim": "user_type",
				"destinationClaim": "type"
			}
		],
		"access": {
			"expires_in": 3600
		},
		"refresh": {
			"expires_in": 2592000,
			"enabled": true
		},
		"anonymousAccess": {
			"expires_in": 2592000,
			"enabled": true
		}
	}`
)

var _ = Describe("AppID", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	Describe("GetTokenConfig", func() {
		Context("When get token config is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/management/v4/074c2fb3-62b5-451d-a7c3-cf4efdc3266f/config/tokens"),
						ghttp.RespondWith(http.StatusOK, tokenConfig),
					),
				)
			})

			It("should get TokenConfig", func() {
				resp, err := newConfig(server.URL()).GetTokenConfig("074c2fb3-62b5-451d-a7c3-cf4efdc3266f")
				Expect(resp).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Access).ShouldNot(BeNil())
				Expect(resp.Access.ExpiresIn).Should(Equal(3600))
				Expect(resp.Refresh).ShouldNot(BeNil())
				Expect(resp.Refresh.ExpiresIn).Should(Equal(2592000))
				Expect(*resp.Refresh.Enabled).Should(BeTrue())
				Expect(resp.AnonymousAccess).ShouldNot(BeNil())
				Expect(resp.AnonymousAccess.ExpiresIn).Should(Equal(2592000))
				Expect(*resp.AnonymousAccess.Enabled).Should(BeTrue())
				Expect(resp.AccessTokenClaims).To(HaveLen(1))
				Expect(resp.AccessTokenClaims[0].DestinationClaim).Should(Equal("type"))
				Expect(resp.AccessTokenClaims[0].SourceClaim).Should(Equal("user_type"))
				Expect(resp.AccessTokenClaims[0].Source).Should(Equal("saml"))
				Expect(resp.IDTokenClaims).To(HaveLen(1))
				Expect(resp.IDTokenClaims[0].DestinationClaim).Should(Equal(""))
				Expect(resp.IDTokenClaims[0].SourceClaim).Should(Equal("theme"))
				Expect(resp.IDTokenClaims[0].Source).Should(Equal("attributes"))
			})
		})
	})

	Context("When unsuccessful", func() {
		BeforeEach(func() {
			server = ghttp.NewServer()
			server.SetAllowUnhandledRequests(true)
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/management/v4/074c2fb3-62b5-451d-a7c3-cf4efdc3266f/config/tokens"),
					ghttp.RespondWith(http.StatusInternalServerError, `Internal server error`),
				),
			)
		})
		It("should return error ", func() {
			_, err := newConfig(server.URL()).GetTokenConfig("074c2fb3-62b5-451d-a7c3-cf4efdc3266f")
			Expect(err).To(HaveOccurred())
		})
	})
})

func newConfig(url string) Config {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url
	conf.MaxRetries = helpers.Int(0) // otherwise it slows down the tests

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.AppIDService,
	}

	return newConfigAPI(&client)
}
