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
	"net/http"
	"time"
)

type AppServerMock struct {
	GenericMock
	oauthClientCfg *models.Client
	oauthCfg       *oauth2.Config
	authServerURL  string
	globalToken    *oauth2.Token // Non-concurrent security
}

func (appSrv *AppServerMock) Start(oauthCfg *oauth2.Config, authServerURL string) {
	appSrv.initGeneric()

	appSrv.oauthCfg = oauthCfg
	appSrv.authServerURL = authServerURL

	appSrv.rootMux.HandleFunc("/", appSrv.rootHandler)
	appSrv.rootMux.HandleFunc("/oauth2", appSrv.oauth2Handler)
	appSrv.rootMux.HandleFunc("/refresh", appSrv.refreshHandler)
	appSrv.rootMux.HandleFunc("/try", appSrv.tryHandler)
	appSrv.rootMux.HandleFunc("/pwd", appSrv.pwdHandler)
	appSrv.rootMux.HandleFunc("/client", appSrv.clientHandler)

	appSrv.startGeneric()
}

func (appSrv *AppServerMock) rootHandler(w http.ResponseWriter, r *http.Request) {
	appSrv.dumpRequest("root", r)

	u := appSrv.oauthCfg.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
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
	token, err := appSrv.oauthCfg.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appSrv.globalToken = token

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
