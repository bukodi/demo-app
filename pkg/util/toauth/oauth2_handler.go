package toauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bukodi/demo-app/pkg/util/errlog"
	"golang.org/x/oauth2"
)

type OAuth2IDP struct {
	Config *oauth2.Config
}

func oauth2Handler(oauthCfg *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errlog.SlogHttpRequest(r, "oauth2 request")
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

		token, err := oauthCfg.Exchange(r.Context(), code /*, oauth2.SetAuthURLParam("code_verifier", "s256example")*/)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new HTTP client using the access token
		client := oauthCfg.Client(r.Context(), token)

		// Make a request to the Google People API to get the user's email
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
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

		// Parse the response body to get the user's email
		var userInfo struct {
			ID            string `json:"id,omitempty"`
			Email         string `json:"email"`
			VerifiedEmail bool   `json:"verified_email"`
			PictureURL    string `json:"picture,omitempty"`
		}
		if err := json.Unmarshal(body, &userInfo); err != nil {
			fmt.Printf("Unable to parse user info: %v\n", err)
			return
		}

		// Print the user's email
		fmt.Printf("User's ID: %s\n", userInfo.ID)
		fmt.Printf("User's email: %s\n", userInfo.Email)
		fmt.Printf("Verified: %t\n", userInfo.VerifiedEmail)

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})
}
