package demo_app

import "embed"

//go:embed _webui/dist
var webuiDist embed.FS

func WebuiDist() embed.FS {
	return webuiDist
}
