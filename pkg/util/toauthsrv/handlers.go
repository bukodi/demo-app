package toauthsrv

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/go-session/session"
	"net/http"
	"net/url"
	"os"
	"time"
)

const sessKeyLoggedInUserID = "LoggedInUserID"

func (tsrv *OAuthTestServer) tokenHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "token", r) // Ignore the error
	}

	err := tsrv.oauthSrv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tsrv *OAuthTestServer) authzHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		dumpRequest(os.Stdout, "authorize", r)
	}

	sessData, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := sessData.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	sessData.Delete("ReturnUri")
	sessData.Save()

	err = tsrv.oauthSrv.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (tsrv *OAuthTestServer) testHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "test", r) // Ignore the error
	}
	token, err := tsrv.oauthSrv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func authnHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", r) // Ignore the error
	}
	sessData, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := sessData.Get(sessKeyLoggedInUserID); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	http.ServeContent(w, r, "authn.html", time.Now(), bytes.NewReader(authHtml))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "login", r) // Ignore the error
	}
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set(sessKeyLoggedInUserID, r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	http.ServeContent(w, r, "login.html", time.Now(), bytes.NewReader(loginHtml))
}

//go:embed auth.html
var authHtml []byte

//go:embed login.html
var loginHtml []byte

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
	}
	sessData, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := sessData.Get(sessKeyLoggedInUserID)
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		sessData.Set("ReturnUri", r.Form)
		sessData.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	sessData.Delete(sessKeyLoggedInUserID)
	sessData.Save()
	return
}
