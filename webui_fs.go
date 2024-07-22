package demo_app

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed _webui/dist
var webuiDist embed.FS

func WebuiDist() fs.FS {
	fs, err := fs.Sub(webuiDist, "_webui/dist")
	if err != nil {
		log.Fatal(err)
	}

	return fs
}

func UIStaticStaticHandler(contextRoot string) (http.Handler, error) {
	fs, err := fs.Sub(webuiDist, "_webui/dist")
	if err != nil {
		log.Fatal(err)
	}

	return StaticHttpHandler("/app/", fs)
}

func StaticHttpHandler(contextRoot string, fs fs.FS) (http.Handler, error) {
	httpFS := http.FS(fs)
	wrappedFS := notFoundRewriteToRootFS{wrappedFS: httpFS}
	httpFsHandler := http.FileServer(wrappedFS)

	return http.StripPrefix(contextRoot, httpFsHandler), nil
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
