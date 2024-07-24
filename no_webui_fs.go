//go:build no_webui

package demo_app

import (
	"fmt"
	"io/fs"
)

func WebuiDist() (fs.FS, error) {
	return nil, fmt.Errorf("static web content isn't embedded into this binary")
}
