//go:build !no_webui

package demo_app

import (
	"embed"
	"io/fs"
)

//go:embed _webui/dist
var webuiDist embed.FS

func WebuiDist() (fs.FS, error) {
	return fs.Sub(webuiDist, "_webui/dist")
}
