package webui

import (
	demoapp "github.com/bukodi/demo-app"
	"github.com/bukodi/demo-app/pkg/server"
	"net/http"
)

func init() {
	server.RegisterPlugin("webui", func(srv *server.ServerInit) error {
		srv.AddRootHandler("/", StaticHttpHandler("/"))
		return nil
	})
}

func StaticHttpHandler(contextRoot string) http.Handler {
	fs, err := demoapp.WebuiDist()
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, err.Error(), http.StatusNotFound)
		})
	}

	httpFS := http.FS(fs)
	wrappedFS := notFoundRewriteToRootFS{wrappedFS: httpFS}
	httpFsHandler := http.FileServer(wrappedFS)

	return http.StripPrefix(contextRoot, httpFsHandler)
}

type notFoundRewriteToRootFS struct {
	wrappedFS http.FileSystem
}

func (w notFoundRewriteToRootFS) Open(name string) (http.File, error) {
	f, err := w.wrappedFS.Open(name)
	if err != nil {
		fRoot, err := w.wrappedFS.Open("/index.html")
		return fRoot, err
	}
	return f, err
}
