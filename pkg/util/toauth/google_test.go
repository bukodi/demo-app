package toauth

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http/cookiejar"

	"github.com/bukodi/demo-app/pkg/util"

	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	//"github.com/bukodi/demo-app/pkg/util/toauthsrv"
	"golang.org/x/oauth2"
)

var (
	//  Open the https://console.cloud.google.com/auth/clients page select theproject and the demo-app client
	//  and set these enviroment variables copy the client id and secret
	DEMO_APP_GOOGLE_CLIENT_ID     = os.Getenv("DEMO_APP_GOOGLE_CLIENT_ID")
	DEMO_APP_GOOGLE_CLIENT_SECRET = os.Getenv("DEMO_APP_GOOGLE_CLIENT_SECRET")
)

func TestGoogleOAuth(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping this test in CI environment")
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	appSrv := httptest.NewUnstartedServer(http.NewServeMux())
	appSrv.Listener = util.Must(net.Listen("tcp", "localhost:9094"))
	appSrv.Start()
	t.Logf("App Server started on %s", appSrv.URL+"/cica")

	// Create oauthClient srv
	oauthCfg := oauth2.Config{
		ClientID:     DEMO_APP_GOOGLE_CLIENT_ID,
		ClientSecret: DEMO_APP_GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL:  appSrv.URL + "/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	appSrv.Config.Handler.(*http.ServeMux).HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		/*u := appSrv.oauthCfg.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))*/
		u := oauthCfg.AuthCodeURL("xyz")
		http.Redirect(w, r, u, http.StatusFound)
	})

	appSrv.Config.Handler.(*http.ServeMux).HandleFunc("/cica", func(w http.ResponseWriter, r *http.Request) {
		/*u := appSrv.oauthCfg.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))*/
		u := oauthCfg.AuthCodeURL("xyz",
			oauth2.SetAuthURLParam("redirectAfterOauth", r.URL.Path),
		)
		http.Redirect(w, r, u, http.StatusFound)
	})

	appSrv.Config.Handler.(*http.ServeMux).Handle("/oauth2", oauth2Handler(&oauthCfg))

	time.Sleep(300 * time.Second)

	// Create an http.Client with the cookie jar
	//client := toauthsrv.NewClientMock(t)
	client := &http.Client{
		Jar: util.Must(cookiejar.New(nil)),
	}

	// Login
	var loginSubmitUrl *url.URL
	resp, err := client.Post(appSrv.URL,
		"application/json",
		bytes.NewBuffer([]byte{}))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("%d %+v", resp.StatusCode, err)
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
		t.Logf("Response: %s", data)
	}

}
