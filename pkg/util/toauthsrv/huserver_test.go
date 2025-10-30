package toauthsrv

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bukodi/demo-app/pkg/util"
)

func TestHttpUtil(t *testing.T) {
	//t.Skipf("Fix this test to finish")
	//done := make(chan bool)

	srv := httptest.NewUnstartedServer(http.NewServeMux())
	srv.Listener, _ = net.Listen("tcp", "localhost:9060")
	srv.Start()
	defer srv.Close()
	t.Logf("Server started on %s", srv.URL)

	srv.Config.Handler.(*http.ServeMux).HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
		t.Logf("Response sent")
	})

	resp, err := http.Get(util.Must(url.JoinPath(srv.URL, "hello")))
	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to get: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "Hello World" {
		t.Fatalf("Expected body Hello World, got %s", string(body))
	}
	t.Logf("Response readed")
}

func TestHttpUtilTLS(t *testing.T) {
	//t.Skipf("Fix this test to finish")
	//done := make(chan bool)

	srv := httptest.NewUnstartedServer(http.NewServeMux())
	srv.Listener, _ = net.Listen("tcp", "localhost:9443")
	srv.StartTLS()
	defer srv.Close()
	t.Logf("Server started on %s", srv.URL)

	srv.Config.Handler.(*http.ServeMux).HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
		t.Logf("Response sent")
	})

	resp, err := srv.Client().Get(util.Must(url.JoinPath(srv.URL, "hello")))
	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to get: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "Hello World" {
		t.Fatalf("Expected body Hello World, got %s", string(body))
	}
	t.Logf("Response readed")
}
