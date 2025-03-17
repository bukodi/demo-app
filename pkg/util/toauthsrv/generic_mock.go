package toauthsrv

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
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
		mockSrv.TestingT.Fatalf("[%s] : http address not set", mockSrv.Name)
		return
	}
	mockSrv.rootMux = http.NewServeMux()
	mockSrv.httpSrv = &http.Server{
		Handler: mockSrv.rootMux,
	}
}

func (mockSrv *GenericMock) startGeneric() {
	if l, err := net.Listen("tcp", mockSrv.Addr); err != nil {
		mockSrv.TestingT.Fatalf("[%s] : %+v", mockSrv.Name, err)
		return
	} else {
		mockSrv.listener = l
	}
	go mockSrv.httpSrv.Serve(mockSrv.listener)
	mockSrv.TestingT.Logf("["+mockSrv.Name+"] : started on http://%s", mockSrv.TCPAddr())
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
		mockSrv.TestingT.Fatalf("[%s] : shutdown failed: %+v", mockSrv.Name, err)
	} else {
		mockSrv.TestingT.Logf("[%s] : shutdown", mockSrv.Name)
	}
}

func (mockSrv *GenericMock) dumpRequest(header string, r *http.Request) {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		mockSrv.TestingT.Errorf("[%s] %s : dump failed %+v", mockSrv.Name, header, err)
	} else {
		mockSrv.TestingT.Logf("[%s] : ---- %s request ----\n%s", mockSrv.Name, header, data)
	}
}
