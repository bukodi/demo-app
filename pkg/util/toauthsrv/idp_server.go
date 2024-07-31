package toauthsrv

import (
	"context"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"net/http"
	"net/http/httputil"
)

type IDPServerMock struct {
	GenericMock
	manager        *manage.Manager
	clientCfgStore *store.ClientStore
	oauthSrv       *server.Server
}

func (idpSrv *IDPServerMock) Start(srvCfg *server.Config) {
	idpSrv.initGeneric()

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

	idpSrv.rootMux.HandleFunc("/login", idpSrv.loginHandler)
	idpSrv.rootMux.HandleFunc("/auth", idpSrv.authnHandler)

	idpSrv.rootMux.HandleFunc("/oauth/authorize", idpSrv.authzHandler)

	idpSrv.rootMux.HandleFunc("/oauth/token", idpSrv.tokenHandler)

	idpSrv.rootMux.HandleFunc("/test", idpSrv.testHandler)

	idpSrv.startGeneric()
	idpSrv.Logf(`
  Point your OAuth client Auth endpoint to https://%s/oauth/authorize
  Point your OAuth client Token endpoint to https://%s/oauth/token`, idpSrv.TCPAddr(), idpSrv.TCPAddr())
}

func (idpSrv *IDPServerMock) SetClientConfig(cliCfg *models.Client) error {
	return idpSrv.clientCfgStore.Set(cliCfg.ID, cliCfg)
}

func (idpSrv *IDPServerMock) dumpRequest(header string, r *http.Request) {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		idpSrv.TestingT.Errorf("IDP Server Request - %s : dump failed %+v", header, err)
	} else {
		idpSrv.Logf("IDP Server Request - %s :\n%s", header, data)
	}
}
