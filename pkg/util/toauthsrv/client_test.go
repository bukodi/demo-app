package toauthsrv

import (
	"bytes"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"golang.org/x/oauth2"
	"io"
	"net/url"
	"strings"
	"testing"
)

var ()

func TestOAuthFlow(t *testing.T) {
	idpSrv := &IDPServerMock{
		GenericMock: GenericMock{
			Name:     "IDPSrv",
			Addr:     "localhost:9096",
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
			Addr:     "localhost:9094",
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

	//time.Sleep(300 * time.Second)

	// Create an http.Client with the cookie jar
	client := NewClientMock(t)

	// Login
	var loginFormPath = "/login"
	resp, err := client.Post("http://"+appSrv.TCPAddr()+"/",
		"application/json",
		bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("%+v", err)
		return
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if !strings.Contains(string(data), loginFormPath) {
			t.Fatalf("Missing %s", loginFormPath)
			return
		}
	}
	_, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	loginFormPostUrl, err := resp.Request.URL.Parse(loginFormPath)
	if err != nil {
		t.Fatalf("%+v", err)
		return
	}
	v := url.Values{}
	v.Add("username", "test")
	v.Add("password", "test")
	resp, err = client.PostForm(loginFormPostUrl.String(), v)

	var authorizePath = "/oauth/authorize"
	if err != nil {
		t.Fatalf("%+v", err)
		return
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if !strings.Contains(string(data), authorizePath) {
			t.Fatalf("Missing %s", authorizePath)
			return
		}
	}

}
