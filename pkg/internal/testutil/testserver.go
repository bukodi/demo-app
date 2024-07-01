package testutil

import (
	"github.com/bukodi/demo-app/pkg/server"
	"testing"
)

type TestServer struct {
	Server *server.Server
	T      *testing.T
}

func StartTestServer(t *testing.T) *TestServer {
	srv := server.NewServer(":0")
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
	return &TestServer{
		Server: srv,
		T:      t,
	}

}

func (ts *TestServer) Stop() {
	err := ts.Server.Stop()
	if err != nil {
		ts.T.Fatal(err)
	}
}
