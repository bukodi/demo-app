package main

import "pkg/server"

func main() {
	srv := server.NewServer()
	srv.Start()

}
