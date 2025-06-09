package toauthsrv

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpUtil(t *testing.T) {
	t.Skipf("Fix this test to finish")
	done := make(chan bool)
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		done <- true
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
	<-done
}
