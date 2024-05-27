package authn

import (
	"context"
	"fmt"
)

const contextKeyAuthData = "CONTEXT_KEY_AUTH_DATA"

type authnData struct {
	user    User
	changed bool
}

func getAuthData(ctx context.Context) *authnData {
	authCtxData, ok := ctx.Value(contextKeyAuthData).(*authnData)
	if !ok || authCtxData == nil {
		return nil
	}
	return authCtxData
}

func GetUser(ctx context.Context) User {
	authCtxData := getAuthData(ctx)
	if authCtxData == nil {
		return nil
	}
	return authCtxData.user
}

func SetUser(ctx context.Context, user User) error {
	authCtxData := getAuthData(ctx)
	if authCtxData == nil {
		return fmt.Errorf("context does not have authentication data (ctx: %v)", ctx)
	}
	authCtxData.user = user
	authCtxData.changed = true
	return nil
}

func WithUser(ctx context.Context, user User) context.Context {
	retCtx := ctx
	authCtxData, ok := retCtx.Value(contextKeyAuthData).(*authnData)
	if !ok || authCtxData == nil {
		authCtxData = &authnData{}
		retCtx = context.WithValue(retCtx, contextKeyAuthData, authCtxData)
	}
	authCtxData.user = user
	authCtxData.changed = true
	return retCtx
}
