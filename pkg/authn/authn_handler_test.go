package authn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bukodi/demo-app/pkg/server"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

func TestSetUser(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	t.Log("TestSetUser")

	tIDP := &testIDP{
		users: []testUser{
			{
				userid:   "alice",
				password: "pswAlice",
			},
			{
				userid:   "bob",
				password: "pswBob",
			},
			{
				userid:   "admin",
				password: "adminPsw",
				roles:    []string{"admin"},
			},
		},
	}
	RegisterIdentityProvider("test_idp", tIDP)

	server.RegisterPlugin("authn_test", func(srv *server.ServerInit) error {
		srv.AddApiHandler("POST /login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var params map[string]string
			err := json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			tU, err := tIDP.FindAndAuthorize(r.Context(), params["userid"], params["password"])
			if err != nil || tU == nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
				return
			}

			SetUser(r.Context(), tU)
			w.WriteHeader(http.StatusOK)
		}))
		srv.AddApiHandler("GET /userid", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := GetUser(r.Context())
			if u == nil {
				http.Error(w, "no user", http.StatusNotFound)
				return
			}
			w.Write([]byte(u.Id()))
		}))
		srv.AddApiHandler("POST /logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SetUser(r.Context(), nil)
			w.WriteHeader(http.StatusOK)
		}))
		return nil
	})

	srv := server.NewServer(":0")
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	// Create a cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create an http.Client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	jsonData := []byte(`{"userid":"alice","password":"pswAlice"}`)
	resp, err := client.Post("http://"+srv.Addr()+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("%+v", err)
		return
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Logf("POST /login success: %s", data)
		baseUrl, _ := url.Parse("http://" + srv.Addr() + "/")
		for _, cookie := range jar.Cookies(baseUrl) {
			fmt.Printf("  %s: %s\n", cookieName, cookie.Value)
		}

	}

	resp, err = client.Get("http://" + srv.Addr() + "/api/v1/userid")
	if err != nil {
		t.Errorf("The HTTP request failed with error %+v", err)
	} else {
		data, _ := io.ReadAll(resp.Body)
		got := string(data)
		want := "alice"
		if got != want {
			t.Errorf("Expected: %s, but actual: %s", want, got)
		}
	}

}
