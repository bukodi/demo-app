package toauthsrv

import (
	"context"
	"github.com/go-oauth2/oauth2/v4/generates"
	"net"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	// Dump requests and responses
	dumpvar bool = true
)

type TestIDPServer struct {
	test           *testing.T
	httpSrv        *http.Server
	rootMux        *http.ServeMux
	listener       net.Listener
	manager        *manage.Manager
	clientCfgStore *store.ClientStore
	oauthSrv       *server.Server
}

func StartTestIDPServer(t *testing.T, addr string, srvCfg *server.Config) *TestIDPServer {
	idpSrv := &TestIDPServer{
		test:    t,
		rootMux: http.NewServeMux(),
	}

	idpSrv.manager = manage.NewDefaultManager()
	idpSrv.manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	idpSrv.manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	idpSrv.manager.MapAccessGenerate(generates.NewAccessGenerate())

	idpSrv.clientCfgStore = store.NewClientStore()
	idpSrv.manager.MapClientStorage(idpSrv.clientCfgStore)

	idpSrv.oauthSrv = server.NewServer(srvCfg, idpSrv.manager)

	idpSrv.oauthSrv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	idpSrv.oauthSrv.SetUserAuthorizationHandler(idpSrv.userAuthorizeHandler)

	idpSrv.oauthSrv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		idpSrv.test.Errorf("Internal Error: %+v", err)
		return
	})

	idpSrv.oauthSrv.SetResponseErrorHandler(func(re *errors.Response) {
		idpSrv.test.Errorf("Response Error: %+v", re.Error)
	})

	idpSrv.rootMux.HandleFunc("/login", idpSrv.loginHandler)
	idpSrv.rootMux.HandleFunc("/auth", idpSrv.authnHandler)

	idpSrv.rootMux.HandleFunc("/oauth/authorize", idpSrv.authzHandler)

	idpSrv.rootMux.HandleFunc("/oauth/token", idpSrv.tokenHandler)

	idpSrv.rootMux.HandleFunc("/test", idpSrv.testHandler)

	if addr == "" {
		idpSrv.test.Fatalf("http address not set")
		return nil
	}
	if l, err := net.Listen("tcp", addr); err != nil {
		idpSrv.test.Fatal(err)
		return nil
	} else {
		idpSrv.listener = l
	}
	idpSrv.httpSrv = &http.Server{
		Handler: idpSrv.rootMux,
	}
	go idpSrv.httpSrv.Serve(idpSrv.listener)

	idpSrv.Logf(`Server started on http://%s
  Point your OAuth client Auth endpoint to https://%s/oauth/authorize
  Point your OAuth client Token endpoint to https://%s/oauth/token`, idpSrv.TCPAddr(), idpSrv.TCPAddr(), idpSrv.TCPAddr())

	return idpSrv
}

func (idpSrv *TestIDPServer) Logf(format string, args ...any) {
	idpSrv.test.Logf("TestIDPSrv : "+format, args...)
}

func (idpSrv *TestIDPServer) TCPAddr() string {
	if idpSrv.listener == nil {
		idpSrv.test.Fatal("listener is nil")
		return ""
	}
	return idpSrv.listener.Addr().String()
}

func (idpSrv *TestIDPServer) Stop() {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := idpSrv.httpSrv.Shutdown(ctx); err != nil {
		idpSrv.test.Fatalf("shutdown failed: %+v", err)
	} else {
		idpSrv.Logf("shutdown")
	}
}

func (idpSrv *TestIDPServer) SetClientConfig(cliCfg *models.Client) error {
	return idpSrv.clientCfgStore.Set(cliCfg.ID, cliCfg)
}

func (idpSrv *TestIDPServer) dumpRequest(header string, r *http.Request) {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		idpSrv.test.Errorf("IDP Server Request - %s : dump failed %+v", header, err)
	} else {
		idpSrv.Logf("IDP Server Request - %s :\n%s", header, data)
	}
}
