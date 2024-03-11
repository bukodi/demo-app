package server

import "testing"

func TestServer(t *testing.T) {
	srv := NewServer(":8080")
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
}
