package toauthsrv

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpUtil(t *testing.T) {
	//t.Skipf("Fix this test to finish")
	//done := make(chan bool)
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
		t.Logf("Response sent")
	})

	srv := httptest.NewUnstartedServer(helloHandler)
	if listener, err := net.Listen("tcp", "localhost:9060"); err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	} else {
		srv.Listener = listener
	}
	srv.Start()
	defer srv.Close()
	t.Logf("Server started on %s", srv.URL)

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "Hello World" {
		t.Fatalf("Expected body Hello World, got %s", string(body))
	}
	t.Logf("Response readed")
}
