package server

import (
	"fmt"
	"net/http"

	demo_app "github.com/bukodi/demo-app"
)

func init() {
	RegisterPlugin("version", func(srv *ServerInit) error {
		srv.AddApiHandler("GET /version", http.HandlerFunc(versionHandler))
		return nil
	})
}

// VersionHandler is a simple HTTP handler that returns the version of the application.
func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%s (%s)", demo_app.Version, demo_app.GitCommit)))
}
