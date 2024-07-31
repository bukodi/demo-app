package webui

import (
	"github.com/bukodi/demo-app/pkg/server"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestIndexHtml(t *testing.T) {
	srv := server.NewServer(":0")
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	resp, err := http.Get("http://" + srv.TCPAddr() + "/index.html")
	if err != nil {
		t.Errorf("The HTTP request failed with error %+v", err)
	} else if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	} else {
		data, _ := io.ReadAll(resp.Body)
		got := string(data)
		if !strings.Contains(got, "<title>") {
			t.Errorf("Wong html content: %s", got)
		}
	}
}
