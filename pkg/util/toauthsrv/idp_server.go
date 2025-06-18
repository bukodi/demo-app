package toauthsrv

import (
	"context"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

type IDPServerMock struct {
	Addr           string
	TestingT       *testing.T
	server         *httptest.Server
	mux            *http.ServeMux
	manager        *manage.Manager
	clientCfgStore *store.ClientStore
	oauthSrv       *server.Server
}

func (idpSrv *IDPServerMock) Start(srvCfg *server.Config) {
	if idpSrv.TestingT == nil {
		panic("TestingT not set")
	}

	idpSrv.mux = http.NewServeMux()

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
		idpSrv.TestingT.Errorf("Internal Error: %+v", err)
		return
	})

	idpSrv.oauthSrv.SetResponseErrorHandler(func(re *errors.Response) {
		idpSrv.TestingT.Errorf("Response Error: %+v", re.Error)
	})

	idpSrv.mux.HandleFunc("/login", idpSrv.loginHandler)
	idpSrv.mux.HandleFunc("/auth", idpSrv.authnHandler)
	idpSrv.mux.HandleFunc("/oauth/authorize", idpSrv.authzHandler)
	idpSrv.mux.HandleFunc("/oauth/token", idpSrv.tokenHandler)
	idpSrv.mux.HandleFunc("/test", idpSrv.testHandler)

	idpSrv.server = httptest.NewUnstartedServer(idpSrv.mux)

	// Set the server to listen on the specified address
	if listener, err := net.Listen("tcp", idpSrv.Addr); err != nil {
		idpSrv.TestingT.Fatalf("[IDPSrv] : Failed to create listener: %v", err)
		return
	} else {
		idpSrv.server.Listener = listener
	}

	idpSrv.server.Start()
	idpSrv.TestingT.Logf("[IDPSrv] : started on http://%s", idpSrv.TCPAddr())
	idpSrv.TestingT.Logf("[IDPSrv] :"+`
  Point your OAuth client Auth endpoint to http://%s/oauth/authorize
  Point your OAuth client Token endpoint to http://%s/oauth/token`, idpSrv.TCPAddr(), idpSrv.TCPAddr())
}

func (idpSrv *IDPServerMock) SetClientConfig(cliCfg *models.Client) error {
	return idpSrv.clientCfgStore.Set(cliCfg.ID, cliCfg)
}

func (idpSrv *IDPServerMock) TCPAddr() string {
	if idpSrv.server == nil || idpSrv.server.Listener == nil {
		return ""
	}
	return idpSrv.server.Listener.Addr().String()
}

func (idpSrv *IDPServerMock) Stop() {
	if idpSrv.server != nil {
		idpSrv.server.Close()
		idpSrv.TestingT.Logf("[IDPSrv] : shutdown")
	}
}

func (idpSrv *IDPServerMock) dumpRequest(header string, r *http.Request) {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		idpSrv.TestingT.Errorf("[IDPSrv] %s : dump failed %+v", header, err)
	} else {
		idpSrv.TestingT.Logf("[IDPSrv] : ---- %s request ----\n%s", header, data)
	}
}
