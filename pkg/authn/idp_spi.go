package authn

import (
	"context"
	"fmt"
)

type IdentityProvider interface {
	Name() string
	FindAndAuthorize(ctx context.Context, userid string, password string) (User, error)
	FindById(ctx context.Context, userid string) User
}

var plugins = make(map[string]IdentityProvider)

func RegisterIdentityProvider(name string, idp IdentityProvider) {
	if idp == nil {
		return
	}
	if _, ok := plugins[name]; ok {
		panic(fmt.Errorf("plugin already registered: %s", name))
	}

	plugins[name] = idp
}
