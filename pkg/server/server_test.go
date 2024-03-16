package server

import (
	"fmt"
	demo_app "github.com/bukodi/demo-app"
	"io"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	srv := NewServer(":0")
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	resp, err := http.Get("http://" + srv.Addr() + "/api/version")
	if err != nil {
		t.Errorf("The HTTP request failed with error %+v", err)
	} else {
		data, _ := io.ReadAll(resp.Body)
		got := string(data)
		want := fmt.Sprintf("%s (%s)", demo_app.Version, demo_app.GitCommit)
		if got != want {
			t.Errorf("Expected: %s, but actual: %s", want, got)
		}
	}
}
