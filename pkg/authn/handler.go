package authn

import (
	"context"
	"encoding/json"
	"github.com/bukodi/demo-app/pkg/server"
	"log/slog"
	"net/http"
	"slices"
)

const cookieName = "X-User-Token"

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := convertCookieToUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authnCtx := WithUser(r.Context(), user)
		r2 := r.WithContext(authnCtx)

		w2 := server.NewResponseWrapper(r2.Context(), w, convertCtxUserToCookie)
		next.ServeHTTP(w2, r2)
	})
}

type userFromJWT struct {
	id      string
	idpName string
	roles   []string
}

func (u userFromJWT) HasRole(role string) bool {
	return slices.Contains(u.roles, role)
}

func (u userFromJWT) Id() string {
	return u.id
}

func (u userFromJWT) IDPName() string {
	return u.idpName
}

var _ User = (*userFromJWT)(nil)

func convertCookieToUser(r *http.Request) (User, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return nil, nil
	}

	var user userFromJWT = userFromJWT{
		id: cookie.Value,
	}
	return &user, nil

}

func convertCtxUserToCookie(requestCtx context.Context, w http.ResponseWriter) {
	authnData := getAuthData(requestCtx)
	if authnData == nil || authnData.user == nil {
		//w.Header().CHeader().Set(cookieName, "")
	} else if authnData.changed {
		http.SetCookie(w, &http.Cookie{
			Name:   cookieName,
			Value:  authnData.user.Id(),
			Secure: false,
		})
	}
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var idpName string
	var user User
	for name, idp := range plugins {
		u, err := idp.FindAndAuthorize(r.Context(), params["userid"], params["password"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if u == nil {
			continue
		}
		if user != nil {
			http.Error(w, "multiple users found", http.StatusInternalServerError)
			return
		}
		user = u
		idpName = name
		break
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	} else {
		SetUser(r.Context(), user)
		slog.Info("User authenticated", "idp", idpName, "userid", user.Id())
	}

	w.WriteHeader(http.StatusCreated)
}
