package toauthsrv

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

type GenericMock struct {
	Name     string
	Addr     string
	TestingT *testing.T
	httpSrv  *http.Server
	rootMux  *http.ServeMux
	listener net.Listener
}

func (mockSrv *GenericMock) initGeneric() {
	if mockSrv.TestingT == nil {
		panic("TestingT not set")
	}
	if mockSrv.Name == "" {
		mockSrv.TestingT.Fatalf("Name not set")
		return
	}
	if mockSrv.Addr == "" {
		mockSrv.TestingT.Fatalf("[" + mockSrv.Name + "] : http address not set")
		return
	}
	mockSrv.rootMux = http.NewServeMux()
	mockSrv.httpSrv = &http.Server{
		Handler: mockSrv.rootMux,
	}
}

func (mockSrv *GenericMock) startGeneric() {
	if l, err := net.Listen("tcp", mockSrv.Addr); err != nil {
		mockSrv.TestingT.Fatalf("["+mockSrv.Name+"] : %+v", err)
		return
	} else {
		mockSrv.listener = l
	}
	go mockSrv.httpSrv.Serve(mockSrv.listener)
	mockSrv.Logf("started on http://%s", mockSrv.TCPAddr())
}

func (mockSrv *GenericMock) Logf(format string, args ...any) {
	mockSrv.TestingT.Logf("["+mockSrv.Name+"] : "+format, args...)
}

func (mockSrv *GenericMock) TCPAddr() string {
	if mockSrv.listener == nil {
		return ""
	}
	return mockSrv.listener.Addr().String()
}

func (mockSrv *GenericMock) Stop() {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := mockSrv.httpSrv.Shutdown(ctx); err != nil {
		mockSrv.TestingT.Fatalf("["+mockSrv.Name+"] : shutdown failed: %+v", err)
	} else {
		mockSrv.Logf("shutdown")
	}
}
