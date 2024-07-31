package toauthsrv

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/go-session/session"
	"net/http"
	"net/url"
	"time"
)

const sessKeyLoggedInUserID = "LoggedInUserID"

func (idpSrv *TestIDPServer) tokenHandler(w http.ResponseWriter, r *http.Request) {
	idpSrv.dumpRequest("token", r)

	err := idpSrv.oauthSrv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (idpSrv *TestIDPServer) authzHandler(w http.ResponseWriter, r *http.Request) {
	idpSrv.dumpRequest("authorize", r)

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

	err = idpSrv.oauthSrv.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (idpSrv *TestIDPServer) testHandler(w http.ResponseWriter, r *http.Request) {
	idpSrv.dumpRequest("test", r)
	token, err := idpSrv.oauthSrv.ValidationBearerToken(r)
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

func (idpSrv *TestIDPServer) authnHandler(w http.ResponseWriter, r *http.Request) {
	idpSrv.dumpRequest("auth", r)
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

func (idpSrv *TestIDPServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	idpSrv.dumpRequest("login", r)

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

func (idpSrv *TestIDPServer) userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	idpSrv.dumpRequest("userAuthorizeHandler", r)

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
