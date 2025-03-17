package toauthsrv

import (
	"bytes"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var ()

func TestGoogleOAuth(t *testing.T) {
	//t.Skip("Skip TestGoogleOAuth")
	//  Open the https://console.cloud.google.com/auth/clients page select theproject and the demo-app client
	//  and set these enviroment variables copy the client id and secret
	DEMO_APP_GOOGLE_CLIENT_ID := os.Getenv("DEMO_APP_GOOGLE_CLIENT_ID")
	DEMO_APP_GOOGLE_CLIENT_SECRET := os.Getenv("DEMO_APP_GOOGLE_CLIENT_SECRET")

	// Create oauthClient srv
	authServerURL := "https://accounts.google.com/o/oauth2/v2/auth"
	config := oauth2.Config{
		ClientID:     DEMO_APP_GOOGLE_CLIENT_ID,
		ClientSecret: DEMO_APP_GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
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

	time.Sleep(300 * time.Second)

	// Create an http.Client with the cookie jar
	client := NewClientMock(t)

	// Login
	var loginSubmitUrl *url.URL
	resp, err := client.Post("http://"+appSrv.TCPAddr()+"/",
		"application/json",
		bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("%+v", err)
		return
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		loginSubmitPath := "/login"
		if !strings.Contains(string(data), fmt.Sprintf(`<form action="%s" method="POST">`, loginSubmitPath)) {
			t.Fatalf("Missing %s", loginSubmitPath)
			return
		}
		loginSubmitUrl, err = resp.Request.URL.Parse(loginSubmitPath)
		if err != nil {
			t.Fatalf("%+v", err)
			return
		}
	}

	var authorizeSubmitUrl *url.URL
	v := url.Values{}
	v.Add("username", "test")
	v.Add("password", "test")
	resp, err = client.PostForm(loginSubmitUrl.String(), v)
	if err != nil {
		t.Fatalf("%+v", err)
		return
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		authorizeSubmitPath := "/oauth/authorize"
		if !strings.Contains(string(data), fmt.Sprintf(`<form action="%s" method="POST">`, authorizeSubmitPath)) {
			t.Fatalf("Missing %s", authorizeSubmitPath)
			return
		}
		authorizeSubmitUrl, err = resp.Request.URL.Parse(authorizeSubmitPath)
		if err != nil {
			t.Fatalf("%+v", err)
			return
		}
	}

	resp, err = client.PostForm(authorizeSubmitUrl.String(), nil)
	if err != nil {
		t.Fatalf("%+v", err)
		return
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if !strings.Contains(string(data), `"access_token"`) {
			t.Fatalf("Missing %s", `"access_token"`)
			return
		}
		client.TestingT.Logf("Response: %s", data)
	}

}
