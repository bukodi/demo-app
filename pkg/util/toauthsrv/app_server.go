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
	"testing"
	"time"
)

type TestAppServer struct {
	test           *testing.T
	httpSrv        *http.Server
	rootMux        *http.ServeMux
	listener       net.Listener
	oauthClientCfg *models.Client
	globalToken    *oauth2.Token // Non-concurrent security
}

func StartTestAppServer(t *testing.T, addr string, oauthCfg *oauth2.Config, authServerURL string) *TestAppServer {
	appSrv := &TestAppServer{
		test:    t,
		rootMux: http.NewServeMux(),
	}

	appSrv.rootMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := oauthCfg.AuthCodeURL("xyz",
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		http.Redirect(w, r, u, http.StatusFound)
	})

	appSrv.rootMux.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
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
		token, err := oauthCfg.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appSrv.globalToken = token

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	appSrv.rootMux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if appSrv.globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		appSrv.globalToken.Expiry = time.Now()
		token, err := oauthCfg.TokenSource(context.Background(), appSrv.globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appSrv.globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	appSrv.rootMux.HandleFunc("/try", func(w http.ResponseWriter, r *http.Request) {
		if appSrv.globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, appSrv.globalToken.AccessToken))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	})

	appSrv.rootMux.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
		token, err := oauthCfg.PasswordCredentialsToken(context.Background(), "test", "test")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appSrv.globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	appSrv.rootMux.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		cfg := clientcredentials.Config{
			ClientID:     oauthCfg.ClientID,
			ClientSecret: oauthCfg.ClientSecret,
			TokenURL:     oauthCfg.Endpoint.TokenURL,
		}

		token, err := cfg.Token(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	if addr == "" {
		appSrv.test.Fatalf("http address not set")
		return nil
	}
	appSrv.httpSrv = &http.Server{
		Addr:    addr,
		Handler: appSrv.rootMux,
	}

	if l, err := net.Listen("tcp", appSrv.httpSrv.Addr); err != nil {
		appSrv.test.Fatal(err)
		return nil
	} else {
		appSrv.listener = l
	}
	go appSrv.httpSrv.Serve(appSrv.listener)
	appSrv.Logf("Server started on http://%s", appSrv.TCPAddr())

	return appSrv
}

func (appSrv *TestAppServer) Logf(format string, args ...any) {
	appSrv.test.Logf("TestAppSrv : "+format, args...)
}

func (appSrv *TestAppServer) TCPAddr() string {
	if appSrv.listener == nil {
		return ""
	}
	return appSrv.listener.Addr().String()
}

func (appSrv *TestAppServer) Stop() {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := appSrv.httpSrv.Shutdown(ctx); err != nil {
		appSrv.test.Fatalf("shutdown failed: %+v", err)
	} else {
		appSrv.Logf("shutdown")
	}
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
