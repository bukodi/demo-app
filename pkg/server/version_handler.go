package server

import (
	"fmt"
	demo_app "github.com/bukodi/demo-app"
	"net/http"
)

func init() {
	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s (%s)", demo_app.Version, demo_app.GitCommit)))
	})
}
