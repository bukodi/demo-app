package toauthsrv

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"
)

type AppServerMock struct {
	Addr           string
	TestingT       *testing.T
	server         *httptest.Server
	mux            *http.ServeMux
	oauthClientCfg *models.Client
	oauthCfg       *oauth2.Config
	authServerURL  string
	globalToken    *oauth2.Token // Non-concurrent security
}

func (appSrv *AppServerMock) Start(oauthCfg *oauth2.Config, authServerURL string) {
	if appSrv.TestingT == nil {
		panic("TestingT not set")
	}

	appSrv.mux = http.NewServeMux()
	appSrv.oauthCfg = oauthCfg
	appSrv.authServerURL = authServerURL

	appSrv.mux.HandleFunc("/", appSrv.rootHandler)
	appSrv.mux.HandleFunc("/oauth2", appSrv.oauth2Handler)
	appSrv.mux.HandleFunc("/refresh", appSrv.refreshHandler)
	appSrv.mux.HandleFunc("/try", appSrv.tryHandler)
	appSrv.mux.HandleFunc("/pwd", appSrv.pwdHandler)
	appSrv.mux.HandleFunc("/client", appSrv.clientHandler)

	appSrv.server = httptest.NewUnstartedServer(appSrv.mux)

	// Set the server to listen on the specified address
	if listener, err := net.Listen("tcp", appSrv.Addr); err != nil {
		appSrv.TestingT.Fatalf("[AppSrv] : Failed to create listener: %v", err)
		return
	} else {
		appSrv.server.Listener = listener
	}

	appSrv.server.Start()
	appSrv.TestingT.Logf("[AppSrv] : started on http://%s", appSrv.TCPAddr())
}

func (appSrv *AppServerMock) rootHandler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("root", r)

	/*u := appSrv.oauthCfg.AuthCodeURL("xyz",
	oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
	oauth2.SetAuthURLParam("code_challenge_method", "S256"))*/
	u := appSrv.oauthCfg.AuthCodeURL("xyz")
	http.Redirect(w, r, u, http.StatusFound)
}

func (appSrv *AppServerMock) oauth2Handler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("oauth2", r)

	r.ParseForm()
	state := r.Form.Get("state")
	if state != "xyz" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := appSrv.oauthCfg.Exchange(context.Background(), code /*, oauth2.SetAuthURLParam("code_verifier", "s256example")*/)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appSrv.globalToken = token

	// Create a new HTTP client using the access token
	client := appSrv.oauthCfg.Client(context.Background(), token)

	// Make a request to the Google People API to get the user's email
	resp, err := client.Get(appSrv.authServerURL + "/test")
	if err != nil {
		fmt.Printf("Unable to get user info: %v\n", err)
		return
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("Invalid status code: %d\n", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read response body: %v\n", err)
		return
	}

	// Print the user's data
	fmt.Printf("User's data: %s\n", string(body))

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func (appSrv *AppServerMock) refreshHandler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("refresh", r)

	if appSrv.globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	appSrv.globalToken.Expiry = time.Now()
	token, err := appSrv.oauthCfg.TokenSource(context.Background(), appSrv.globalToken).Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appSrv.globalToken = token
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func (appSrv *AppServerMock) tryHandler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("try", r)

	if appSrv.globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", appSrv.authServerURL, appSrv.globalToken.AccessToken))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func (appSrv *AppServerMock) pwdHandler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("pwd", r)

	token, err := appSrv.oauthCfg.PasswordCredentialsToken(context.Background(), "test", "test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appSrv.globalToken = token
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func (appSrv *AppServerMock) clientHandler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("client", r)

	cfg := clientcredentials.Config{
		ClientID:     appSrv.oauthCfg.ClientID,
		ClientSecret: appSrv.oauthCfg.ClientSecret,
		TokenURL:     appSrv.oauthCfg.Endpoint.TokenURL,
	}

	token, err := cfg.Token(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func (appSrv *AppServerMock) TCPAddr() string {
	if appSrv.server == nil || appSrv.server.Listener == nil {
		return ""
	}
	return appSrv.server.Listener.Addr().String()
}

func (appSrv *AppServerMock) Stop() {
	if appSrv.server != nil {
		appSrv.server.Close()
		appSrv.TestingT.Logf("[AppSrv] : shutdown")
	}
}

func (appSrv *AppServerMock) dumpRequest(header string, r *http.Request) {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		appSrv.TestingT.Errorf("[AppSrv] %s : dump failed %+v", header, err)
	} else {
		appSrv.TestingT.Logf("[AppSrv] : ---- %s request ----\n%s", header, data)
	}
}
