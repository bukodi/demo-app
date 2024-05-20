package authn

import (
	"context"
)

type User interface {
	HasRole(role string) bool
}

const contextKeyAuthData = "CONTEXT_KEY_AUTH_DATA"

type authData struct {
	user    User
	changed bool
}

func getAuthData(ctx context.Context) *authData {
	authCtxData, ok := ctx.Value(contextKeyAuthData).(*authData)
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

func WithUser(ctx context.Context, user User) context.Context {
	retCtx := ctx
	authCtxData, ok := retCtx.Value(contextKeyAuthData).(*authData)
	if !ok || authCtxData == nil {
		authCtxData = &authData{}
		retCtx = context.WithValue(retCtx, contextKeyAuthData, authCtxData)
	}
	authCtxData.user = user
	authCtxData.changed = true
	return retCtx
}
