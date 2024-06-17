package toauthsrv

import (
	"context"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/generates"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
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
	// The client id being passed in
	idvar string = "222222"
	// The client secret being passed in
	secretvar string = "22222222"
	// The domain of the redirect url
	domainvar string = "http://localhost:9094"
	// the base port for the server
	portvar int = 9094
)

type OAuthTestServer struct {
	listener net.Listener
	httpSrv  *http.Server
	rootMux  *http.ServeMux
	manager  *manage.Manager
	oauthSrv *server.Server
}

func StartOauthTestServer(srvAddr string, srvCfg *server.Config, clientCfg *models.Client) (*OAuthTestServer, error) {
	tsrv := &OAuthTestServer{
		rootMux: http.NewServeMux(),
		httpSrv: &http.Server{
			Addr: srvAddr,
		},
	}
	tsrv.httpSrv.Handler = tsrv.rootMux

	tsrv.manager = manage.NewDefaultManager()
	tsrv.manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	tsrv.manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	tsrv.manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	clientStore.Set(clientCfg.ID, clientCfg)
	tsrv.manager.MapClientStorage(clientStore)

	tsrv.oauthSrv = server.NewServer(srvCfg, tsrv.manager)

	tsrv.oauthSrv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	tsrv.oauthSrv.SetUserAuthorizationHandler(userAuthorizeHandler)

	tsrv.oauthSrv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		slog.Error(fmt.Sprintf("Internal Error: %s", err.Error()), "err", err)
		return
	})

	tsrv.oauthSrv.SetResponseErrorHandler(func(re *errors.Response) {
		slog.Error(fmt.Sprintf("Response Error: %s", re.Error.Error()), "err", re.Error)
	})

	tsrv.rootMux.HandleFunc("/login", loginHandler)
	tsrv.rootMux.HandleFunc("/auth", authnHandler)

	tsrv.rootMux.HandleFunc("/oauth/authorize", tsrv.authzHandler)

	tsrv.rootMux.HandleFunc("/oauth/token", tsrv.tokenHandler)

	tsrv.rootMux.HandleFunc("/test", tsrv.testHandler)

	l, err := net.Listen("tcp", tsrv.httpSrv.Addr)
	if err != nil {
		return nil, err
	}
	tsrv.listener = l
	go tsrv.httpSrv.Serve(l)
	slog.Info(fmt.Sprintf("Server started on http://%s", tsrv.listener.Addr()))
	slog.Info(fmt.Sprintf("Point your OAuth client Auth endpoint to https://%s/oauth/authorize", l.Addr().String()))
	slog.Info(fmt.Sprintf("Point your OAuth client Token endpoint to https://%s/oauth/token", l.Addr().String()))

	return tsrv, nil
}

func (tsrv *OAuthTestServer) Stop() error {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := tsrv.httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Server shutdown")
	return nil
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}
