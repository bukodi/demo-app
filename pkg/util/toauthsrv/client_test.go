package toauthsrv

import (
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"golang.org/x/oauth2"
	"testing"
)

var ()

func TestOAuthFlow(t *testing.T) {
	idpSrv := &IDPServerMock{
		GenericMock: GenericMock{
			Name:     "IDPSrv",
			Addr:     "localhost:0",
			TestingT: t,
		},
	}
	idpSrv.Start(server.NewConfig())
	defer idpSrv.Stop()

	// Create oauthClient srv
	authServerURL := "http://localhost:9096"
	config := oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}

	appSrv := &AppServerMock{
		GenericMock: GenericMock{
			Name:     "AppSrv",
			Addr:     "localhost:0",
			TestingT: t,
		},
	}
	appSrv.Start(&config, authServerURL)
	defer appSrv.Stop()

	idpSrv.SetClientConfig(&models.Client{
		ID:     config.ClientID,
		Secret: config.ClientSecret,
		Domain: "http://localhost:9094",
	})

}
