package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
)

type ServerInit struct {
	srv *Server
}

type ServerPluginInit func(srv *ServerInit) error

var plugins = make(map[string]ServerPluginInit)

func RegisterPlugin(name string, srvPlugin ServerPluginInit) {
	if srvPlugin == nil {
		return
	}
	if _, ok := plugins[name]; ok {
		panic(fmt.Errorf("plugin already registered: %s", name))
	}

	plugins[name] = srvPlugin
}

// initPlugins initializes all registered plugins in sorted name order
func (srv *Server) initPlugins() (retErr error) {
	pluginNames := make([]string, 0, len(plugins))
	for name, _ := range plugins {
		pluginNames = append(pluginNames, name)
	}
	sort.Strings(pluginNames)

	srvInit := &ServerInit{srv: srv}

	for _, name := range pluginNames {
		pi := plugins[name]
		if err := pi(srvInit); err != nil {
			retErr = errors.Join(retErr, fmt.Errorf("plugin %s init failed: %w", name, err))
		} else {
			slog.Debug("plugin initialized", "name", name)
		}
	}
	return
}

func (srvInit *ServerInit) AddRootHandler(pattern string, handler http.Handler) {
	srvInit.srv.rootMux.Handle(pattern, handler)
	slog.Debug("Root handler added", "pattern", pattern)
}

func (srvInit *ServerInit) AddApiHandler(pattern string, handler http.Handler) {
	srvInit.srv.apiMux.Handle(pattern, handler)
	slog.Debug("Api handler added", "pattern", pattern)
}

type MiddlewareFactory func(http.Handler) http.Handler

func (srvInit *ServerInit) AddMiddleware(mwFactory MiddlewareFactory) {
	mwHandler := mwFactory(srvInit.srv.httpSrv.Handler)
	srvInit.srv.httpSrv.Handler = mwHandler
	slog.Debug("Middleware init", "middleware", fmt.Sprintf("%T", mwHandler))
}
